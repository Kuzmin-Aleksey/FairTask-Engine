package order_balancer

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"context"
	"fmt"
	"log"
	"log/slog"
	"slices"
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
	//GetParameters(ctx context.Context, id int) ([]entity.ExecutorParameter, error)
	GetById(ctx context.Context, id int) (*aggregate.ExecutorWithParams, error)
	GetActive(ctx context.Context) ([]aggregate.ExecutorWithParams, error)
	SetStatus(ctx context.Context, id int, status string) error
}

type AIS interface {
	SendOrderExecutor(ctx context.Context, orderId, executorId int) error
}

type OrderBalancerService struct {
	ordersRepo    OrdersRepo
	executorsRepo ExecutorsRepo
	AIS           AIS

	freeExecutors *executorsList

	ordersPool []aggregate.OrderWithParameter
}

func New(ordersRepo OrdersRepo, executorsRepo ExecutorsRepo, aic AIS) *OrderBalancerService {
	s := &OrderBalancerService{
		ordersRepo:    ordersRepo,
		executorsRepo: executorsRepo,
		AIS:           aic,
		freeExecutors: &executorsList{},
	}

	if err := s.loadFreeExecutors(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println(s.freeExecutors.len())
	return s
}

func (s *OrderBalancerService) NewOrder(ctx context.Context, order *aggregate.OrderWithParameter) error {
	const op = "OrderBalancerService.HandleOrder"

	order.Status = value.OrderStatusProcessed

	executorId, err := s.FindExecutor(ctx, order)
	if err != nil {
		return err
	}

	if executorId == 0 {
		contextx.GetLoggerOrDefault(ctx).InfoContext(ctx, "executor not found", slog.Any("order", order))
	} else {
		order.ExecutorId = executorId

		if err := s.AIS.SendOrderExecutor(ctx, order.Id, executorId); err != nil {
			contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "send order error", logx.Error(err))
		}
	}

	contextx.GetLoggerOrDefault(ctx).InfoContext(ctx, "executor found", slog.Any("order", order))

	if err := s.ordersRepo.Save(ctx, order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *OrderBalancerService) SetExecutorStatus(ctx context.Context, id int, status string) error {
	const op = "OrderBalancerService.SetExecutorStatus"
	if err := s.executorsRepo.SetStatus(ctx, id, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if status == "inactive" {
		s.freeExecutors.del(id)
	} else {
		executor, err := s.executorsRepo.GetById(ctx, id)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		s.freeExecutors.insertByOrderCount(executor)
	}

	return nil
}

func (s *OrderBalancerService) loadFreeExecutors(ctx context.Context) error {
	const op = "OrderBalancerService.loadFreeExecutors"

	executors, err := s.executorsRepo.GetActive(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	slices.SortFunc(executors, func(e1 aggregate.ExecutorWithParams, e2 aggregate.ExecutorWithParams) int {
		return e1.OrderCount - e2.OrderCount
	})

	for _, executor := range executors {
		s.freeExecutors.addToEnd(&executor)
	}

	return nil
}

func (s *OrderBalancerService) UpdateOrderStatus(ctx context.Context, orderId int, status value.OrderStatus) error {
	const op = "OrderBalancerService.UpdateOrderStatus"

	order, err := s.ordersRepo.GetById(ctx, orderId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if status == value.OrderStatusProcessed {
		params, err := s.ordersRepo.GetParameters(ctx, orderId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		executor, err := s.executorsRepo.GetById(ctx, order.ExecutorId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if executor.Status == "inactive" || !checkExecutorParams(executor, params) {
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
	}

	if err := s.ordersRepo.SetStatus(ctx, orderId, status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *OrderBalancerService) FindExecutor(ctx context.Context, order *aggregate.OrderWithParameter) (int, error) {
	const op = "OrderBalancerService.FindExecutor"

	if order.ParentId != 0 {
		lastOrder, err := s.ordersRepo.GetById(ctx, order.ParentId)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		lastExecutor, err := s.executorsRepo.GetById(ctx, lastOrder.ExecutorId)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		if lastExecutor.Status == "active" && checkExecutorParams(lastExecutor, order.Parameters) {
			return lastExecutor.Id, nil
		}
	}

	if len(order.Parameters) == 0 {
		executor := s.freeExecutors.getFirsAndDel()
		if executor == nil {
			return 0, nil
		}
		return executor.Id, nil
	}

	executor := s.freeExecutors.findByParamsAndDelete(order.Parameters)
	if executor == nil {
		contextx.GetLoggerOrDefault(ctx).WarnContext(ctx, "executor not found", slog.Any("order", order))
		return 0, nil
	}

	return executor.Id, nil
}
