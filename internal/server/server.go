package server

type Server struct {
	balancer *BalancerServer
}

func NewServer(
	products *BalancerServer,
) *Server {
	return &Server{
		balancer: products,
	}
}
