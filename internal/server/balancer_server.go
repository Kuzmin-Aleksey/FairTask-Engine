package server

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/service/order_balancer"
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/failure"
	"encoding/json"
	"net/http"
	"strconv"
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

func (s *BalancerServer) ApiHandleGetExecutors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	executors, err := s.balancer.GetAllExecutors(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, executors, http.StatusOK)
}

func (s *BalancerServer) ApiHandleNewExecutor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	executor := &aggregate.ExecutorWithParams{}

	if err := json.NewDecoder(r.Body).Decode(executor); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.balancer.CreateExecutor(ctx, executor); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, IdResponse{executor.Id}, http.StatusOK)

}
func (s *BalancerServer) ApiHandleDeleteExecutor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
	}

	if err := s.balancer.DeleteExecutor(ctx, id); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}

type createExecutorParameterRequest struct {
	ExecutorId  int    `json:"executor_id"`
	ParameterId int    `json:"parameter_id"`
	Mask        string `json:"mask"`
}

func (s *BalancerServer) ApiHandleAddExecutorParameter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	param := new(createExecutorParameterRequest)

	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.balancer.AddExecutorParameter(ctx, param.ExecutorId, &entity.ExecutorParameter{
		Id:   param.ParameterId,
		Mask: param.Mask,
	}); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}

type deleteExecutorParameterRequest struct {
	ExecutorId  int `json:"executor_id"`
	ParameterId int `json:"parameter_id"`
}

func (s *BalancerServer) ApiHandleDeleteExecutorParameter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param := new(deleteExecutorParameterRequest)

	if err := json.NewDecoder(r.Body).Decode(param); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.balancer.DeleteExecutorParameter(ctx, param.ExecutorId, param.ParameterId); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}
}
