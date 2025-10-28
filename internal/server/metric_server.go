package server

import (
	"FairTask_Engine/internal/domain/service/metric_service"
	"FairTask_Engine/pkg/failure"
	"net/http"
)

type MetricServer struct {
	metric *metric_service.MetricService
}

func NewMetricServer(metric *metric_service.MetricService) *MetricServer {
	return &MetricServer{
		metric: metric,
	}
}

func (s *MetricServer) ApiHandleGetExecutorsOrdersCountList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.metric.GetExecutorsOrdersCountList(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}

	writeJson(ctx, w, list, http.StatusOK)
}

func (s *MetricServer) ApiHandleGetOrderCountByLimit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := r.FormValue("limit")

	timetable, err := s.metric.GetOrderCountByLimit(ctx, limit)
	if err != nil {
		writeAndLogErr(ctx, w, err)
		return
	}

	writeJson(ctx, w, timetable, http.StatusOK)
}

type responseCount struct {
	Count int `json:"count"`
}

func (s *MetricServer) ApiHandleGetOllOrderCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	count, err := s.metric.GetOllOrderCount(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}
	writeJson(ctx, w, responseCount{count}, http.StatusOK)
}

func (s *MetricServer) ApiHandleGetCompleteOrderCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	count, err := s.metric.GetCompleteOrderCount(ctx)
	if err != nil {
		writeAndLogErr(ctx, w, failure.NewInvalidRequestError(err.Error()))
		return
	}
	writeJson(ctx, w, responseCount{count}, http.StatusOK)
}
