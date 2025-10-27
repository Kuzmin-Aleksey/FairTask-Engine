package server

type Server struct {
	products *BalancerServer
}

func NewServer(
	products *BalancerServer,
) *Server {
	return &Server{
		products: products,
	}
}
