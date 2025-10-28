package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/receiver", s.balancer.ApiHandleNewOrder).Methods(http.MethodPost)
	rtr.HandleFunc("/update_status", s.balancer.ApiHandleUpdateStatus).Methods(http.MethodPost)
	rtr.HandleFunc("/executor/set_status", s.balancer.ApiHandleSetExecutorStatus).Methods(http.MethodPost)

	rtr.HandleFunc("/metric/order_count", s.metric.ApiHandleGetOrderCountByLimit).Methods(http.MethodGet)
	rtr.HandleFunc("/metric/executors_order_count", s.metric.ApiHandleGetExecutorsOrdersCountList).Methods(http.MethodGet)

	rtr.HandleFunc("/parameters/create", s.parameters.ApiHandleCreateParameter).Methods(http.MethodPost)
	rtr.HandleFunc("/parameters", s.parameters.ApiHandleGetParameters).Methods(http.MethodGet)
	rtr.HandleFunc("/parameters/delete", s.parameters.ApiHandleDeleteParameter).Methods(http.MethodPost)

}
