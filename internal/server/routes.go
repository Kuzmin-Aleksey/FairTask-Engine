package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/receiver", s.balancer.ApiHandleNewOrder).Methods(http.MethodPost)
	rtr.HandleFunc("/update_status", s.balancer.ApiHandleUpdateStatus).Methods(http.MethodPost)
	rtr.HandleFunc("/executors", s.balancer.ApiHandleGetExecutors).Methods(http.MethodGet)
	rtr.HandleFunc("/executors/create", s.balancer.ApiHandleNewExecutor).Methods(http.MethodPost)
	rtr.HandleFunc("/executors/set_status", s.balancer.ApiHandleSetExecutorStatus).Methods(http.MethodPost)
	rtr.HandleFunc("/executors/delete", s.balancer.ApiHandleDeleteExecutor).Methods(http.MethodPost)
	rtr.HandleFunc("/executors/parameters/add", s.balancer.ApiHandleAddExecutorParameter).Methods(http.MethodPost)
	rtr.HandleFunc("/executors/parameters/del", s.balancer.ApiHandleDeleteExecutorParameter).Methods(http.MethodPost)

	rtr.HandleFunc("/metric/order_count", s.metric.ApiHandleGetOllOrderCount).Methods(http.MethodGet)
	rtr.HandleFunc("/metric/order_count_limit", s.metric.ApiHandleGetOrderCountByLimit).Methods(http.MethodGet)
	rtr.HandleFunc("/metric/complete_order_count", s.metric.ApiHandleGetCompleteOrderCount).Methods(http.MethodGet)

	rtr.HandleFunc("/parameters/create", s.parameters.ApiHandleCreateParameter).Methods(http.MethodPost)
	rtr.HandleFunc("/parameters", s.parameters.ApiHandleGetParameters).Methods(http.MethodGet)
	rtr.HandleFunc("/parameters/delete", s.parameters.ApiHandleDeleteParameter).Methods(http.MethodPost)

}
