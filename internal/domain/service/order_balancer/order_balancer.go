package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"
)

type OrdersRepo interface {
	Save(ctx context.Context, order *aggregate.OrderWithParameter) error
	UpdateExecutor(ctx context.Context, id int, executorId int) error
	GetById(ctx context.Context, id int) (*entity.Order, error)
	GetParameters(ctx context.Context, id int) ([]entity.Parameter, error)
	SetEnabled(ctx context.Context, id int, enabled bool) error
	Delete(ctx context.Context, id int) error
}

type ExecutorsRepo interface {
	GetParameters(ctx context.Context, id int) ([]entity.Parameter, error)
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

func (s *OrderBalancerService) NewOrder(ctx context.Context, order *aggregate.OrderWithParameter) error {
	const op = "OrderBalancerService.HandleOrder"

	executorId, err := s.FindExecutor(ctx, order)
	if err != nil {
		return err
	}

	order.ExecutorId = executorId

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

		executorParamsMap := mapParams(lastExecutorParameters)

		for _, param := range params {
			if val, ok := executorParamsMap[param.Name]; !ok || val != param.Value {
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

				break
			}
		}

		if err := s.ordersRepo.SetEnabled(ctx, orderId, true); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

	case value.OrderStatusAwait:
		if err := s.ordersRepo.SetEnabled(ctx, orderId, false); err != nil {
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

	var executorsWithMinOrders []entity.Executor
	minValue := executors[0].OrderCount

	for _, executor := range executors {
		if executor.OrderCount < minValue {
			minValue = executor.OrderCount
			executorsWithMinOrders = []entity.Executor{}
		}
		if executor.OrderCount == minValue {
			executorsWithMinOrders = append(executorsWithMinOrders, executor.Executor)
		}
	}

	executor := executorsWithMinOrders[rand.IntN(len(executorsWithMinOrders))]

	return executor.Id, nil
}

func filterByParameters(executors []aggregate.ExecutorWithParams, parameters []entity.Parameter) []aggregate.ExecutorWithParams {
	paramMap := mapParams(parameters)

	var resExecutors []aggregate.ExecutorWithParams

	for _, executor := range executors {
		var equalParams int

		for _, param := range executor.Parameters {
			if v, ok := paramMap[param.Name]; ok && v == param.Value {
				equalParams++
				if equalParams == len(parameters) {
					resExecutors = append(resExecutors, executor)
				}
			}
		}
	}
	return resExecutors
}

func mapParams(params []entity.Parameter) map[string]string {
	m := make(map[string]string)
	for _, param := range params {
		m[param.Name] = param.Value
	}
	return m
}
