package server

type Server struct {
	balancer   *BalancerServer
	metric     *MetricServer
	parameters *ParametersServer
}

func NewServer(
	products *BalancerServer,
	metric *MetricServer,
	parameters *ParametersServer,
) *Server {
	return &Server{
		balancer:   products,
		metric:     metric,
		parameters: parameters,
	}
}
