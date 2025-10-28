package server

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/service/order_balancer"
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/failure"
	"encoding/json"
	"net/http"
)

type BalancerServer struct {
	balancer *order_balancer.OrderBalancerService
}

func NewBalancerServer(balancer *order_balancer.OrderBalancerService) *BalancerServer {
	return &BalancerServer{balancer: balancer}
}

type ApiHandleNewOrderRequest struct {
	Order      entity.Order            `json:"order"`
	Parameters []entity.OrderParameter `json:"order_params"`
}

func (s *BalancerServer) ApiHandleNewOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(ApiHandleNewOrderRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	orderWithParams := &aggregate.OrderWithParameter{
		Order:      req.Order,
		Parameters: req.Parameters,
	}

	if err := s.balancer.NewOrder(ctx, orderWithParams); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}

type ApiHandleUpdateStatusRequest struct {
	OrderId int               `json:"order_id"`
	Status  value.OrderStatus `json:"status"`
}

func (s *BalancerServer) ApiHandleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(ApiHandleUpdateStatusRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
	}

	if err := s.balancer.UpdateOrderStatus(ctx, req.OrderId, req.Status); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}

type ApiHandleSetExecutorStatusRequest struct {
	ExecutorId int    `json:"executor_id"`
	Status     string `json:"status"`
}

func (s *BalancerServer) ApiHandleSetExecutorStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(ApiHandleSetExecutorStatusRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.balancer.SetExecutorStatus(ctx, req.ExecutorId, req.Status); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}
