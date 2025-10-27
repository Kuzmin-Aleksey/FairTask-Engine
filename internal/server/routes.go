package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	get  = http.MethodGet
	post = http.MethodPost
)

func (s *Server) InitRoutes(rtr *mux.Router) {
	rtr.HandleFunc("/receiver", s.balancer.ApiHandleNewOrder).Methods(http.MethodPost)
	rtr.HandleFunc("/update_status", s.balancer.ApiHandleUpdateStatus).Methods(http.MethodPost)
}
