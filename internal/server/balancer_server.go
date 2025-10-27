package server

import (
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/service/order_balancer"
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
	Order      entity.Order       `json:"order"`
	Parameters []entity.Parameter `json:"parameters"`
}

func (s *BalancerServer) ApiHandleNewOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := new(ApiHandleNewOrderRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	if err := s.balancer.NewOrder(ctx, nil); err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

}

/*
{
  "order": {
    "id": 0,
    "parent_id": 0,
    "text": "string",
    "status": "processed"
  },
  "order_params": [
		{
			"id": 1,
			"params_id": 0,
			"value": "string"
		},
		{
			"id": 2,
			"params_id": 0,
			"value": "string"
		}
  ]
}

*/
