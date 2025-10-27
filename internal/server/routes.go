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

}
