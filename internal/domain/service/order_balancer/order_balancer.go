package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

type OrdersRepo interface {
	Save(ctx context.Context, order *aggregate.OrderWithParameter) error
	UpdateExecutor(ctx context.Context, id int, executorId int) error
	GetById(ctx context.Context, id int) (*entity.Order, error)
	GetParameters(ctx context.Context, id int) ([]entity.OrderParameter, error)
	SetStatus(ctx context.Context, id int, status value.OrderStatus) error
	Delete(ctx context.Context, id int) error
}

type ExecutorsRepo interface {
	GetParameters(ctx context.Context, id int) ([]entity.ExecutorParameter, error)
	GetActive(ctx context.Context) ([]aggregate.ExecutorWithParams, error)
}

type AIC interface {
	SendOrderExecutor(ctx context.Context, orderId, executorId int) error
}

type OrderBalancerService struct {
	ordersRepo    OrdersRepo
	executorsRepo ExecutorsRepo
	aic           AIC
	mu            sync.Mutex
}

func New(ordersRepo OrdersRepo, executorsRepo ExecutorsRepo, aic AIC) *OrderBalancerService {
	return &OrderBalancerService{
		ordersRepo:    ordersRepo,
		executorsRepo: executorsRepo,
		aic:           aic,
	}
}

var errExecutorNotFound = errors.New("executor not found")

func (s *OrderBalancerService) NewOrder(ctx context.Context, order *aggregate.OrderWithParameter) error {
	const op = "OrderBalancerService.HandleOrder"

	executorId, err := s.FindExecutor(ctx, order)
	if err != nil {
		return err
	}

	if executorId == 0 {
		return fmt.Errorf("%s: %w", op, errExecutorNotFound)
	}

	order.ExecutorId = executorId
	order.Status = value.OrderStatusProcessed

	contextx.GetLoggerOrDefault(ctx).InfoContext(ctx, "new order", slog.Any("order", order))

	if err := s.ordersRepo.Save(ctx, order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := s.aic.SendOrderExecutor(ctx, order.Id, executorId); err != nil {
			contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "send order error", logx.Error(err))
		}
	}()

	return nil
}

func (s *OrderBalancerService) UpdateOrderStatus(ctx context.Context, orderId int, status value.OrderStatus) error {
	const op = "OrderBalancerService.UpdateOrderStatus"

	order, err := s.ordersRepo.GetById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch status {
	case value.OrderStatusReject, value.OrderStatusAccept:
		if err := s.ordersRepo.Delete(ctx, orderId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case value.OrderStatusProcessed:
		params, err := s.ordersRepo.GetParameters(ctx, orderId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		lastExecutorParameters, err := s.executorsRepo.GetParameters(ctx, order.ExecutorId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if len(filterByParameters([]aggregate.ExecutorWithParams{
			{
				Parameters: lastExecutorParameters,
			},
		}, params)) == 0 {
			// find new executor
			executorId, err := s.FindExecutor(ctx, &aggregate.OrderWithParameter{
				Order:      *order,
				Parameters: params,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}

			if err := s.ordersRepo.UpdateExecutor(ctx, orderId, executorId); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

		if err := s.ordersRepo.SetStatus(ctx, orderId, status); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case value.OrderStatusAwait:
		if err := s.ordersRepo.SetStatus(ctx, orderId, status); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *OrderBalancerService) FindExecutor(ctx context.Context, order *aggregate.OrderWithParameter) (int, error) {
	const op = "OrderBalancerService.FindExecutor"

	s.mu.Lock()
	defer s.mu.Unlock()

	executors, err := s.executorsRepo.GetActive(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	executors = filterByParameters(executors, order.Parameters)

	if len(executors) == 0 {
		contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "executor not found", slog.Any("order", order))
		return 0, nil
	}

	if order.ParentId != 0 {
		lastOrder, err := s.ordersRepo.GetById(ctx, order.ParentId)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		for _, executor := range executors {
			if lastOrder.ExecutorId == executor.Id {
				return executor.Id, nil
			}
		}
	}

	var executorsWithMinOrders []aggregate.ExecutorWithParams
	minValue := executors[0].OrderCount

	for _, executor := range executors {
		if executor.OrderCount < minValue {
			minValue = executor.OrderCount
			executorsWithMinOrders = []aggregate.ExecutorWithParams{}
		}
		if executor.OrderCount == minValue {
			executorsWithMinOrders = append(executorsWithMinOrders, executor)
		}
	}

	executors = executorsWithMinOrders

	// filter by MaxDailyLimit

	var freeExecutors []aggregate.ExecutorWithParams

	for _, executor := range executors {
		if executor.MaxDailyLimit-executor.OrderCount > 0 {
			freeExecutors = append(freeExecutors, executor)
		}
	}

	if len(freeExecutors) != 0 {
		executors = freeExecutors
	}

	executor := executors[rand.IntN(len(executors))]

	return executor.Id, nil
}

func filterByParameters(executors []aggregate.ExecutorWithParams, parameters []entity.OrderParameter) []aggregate.ExecutorWithParams {
	mappedParams := mapParams(parameters)
	paramsLen := len(parameters)

	var resExecutors []aggregate.ExecutorWithParams

	for _, executor := range executors {
		var count int
		for _, executorParam := range executor.Parameters {
			if param, ok := mappedParams[executorParam.Id]; ok && matchParam(executorParam, param) {
				log.Println("param", param, "ok")
				count++
				if count == paramsLen {
					resExecutors = append(resExecutors, executor)
				}
			}
		}

	}

	return resExecutors
}

func mapParams(params []entity.OrderParameter) map[int]string {
	m := make(map[int]string)
	for _, param := range params {
		m[param.Id] = param.Value
	}
	return m
}

func matchParam(param entity.ExecutorParameter, val string) bool {
	switch param.Type {
	case value.ParameterTypeText, value.ParameterTypeBool:
		return param.Mask == val
	case value.ParameterTypeInt, value.ParameterTypeFloat, value.ParameterTypeDatetime:
		parsedMask := parseMask(param.Mask, val)
		if len(parsedMask) == 1 {
			return param.Mask == val
		}

		switch param.Type {
		case value.ParameterTypeInt:
			return compareInt(parsedMask)
		case value.ParameterTypeFloat:
			return compareFloat(parsedMask)
		case value.ParameterTypeDatetime:
			return compareDatetime(parsedMask)
		}
	}

	return false
}

// TODO use interface
func compareInt(s []string) bool {
	lastN, _ := strconv.Atoi(s[0])
	for _, v := range s[1:] {
		n, _ := strconv.Atoi(v)
		if lastN >= n {
			return false
		}
	}
	return true
}
func compareFloat(s []string) bool {
	lastN, _ := strconv.ParseFloat(s[0], 64)
	for _, v := range s[1:] {
		n, _ := strconv.ParseFloat(v, 64)
		if lastN >= n {
			return false
		}
	}
	return true
}
func compareDatetime(s []string) bool {
	lastTime, _ := time.Parse(time.DateTime, s[0])
	for _, v := range s[1:] {
		t, _ := time.Parse(time.DateTime, v)
		if lastTime.After(t) {
			return false
		}
	}
	return true
}

func parseMask(mask string, val string) []string {
	xIdx := strings.Index(mask, "x")
	if xIdx == -1 {
		return []string{mask}
	}

	if xIdx == 0 {
		return []string{val, mask[1:]}
	}
	if xIdx == len(mask)-1 {
		return []string{mask[:len(val)-1], val}
	}

	return []string{mask[:xIdx], val, mask[xIdx+1:]}
}
