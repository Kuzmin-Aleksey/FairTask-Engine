package metric_service

import (
	"FairTask_Engine/internal/domain/value"
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"context"
	"fmt"
	"slices"
	"time"
)

type Repo interface {
	GetExecutorsOrdersCountList(ctx context.Context) ([]int, error)
	GetCountByPeriod(ctx context.Context, start time.Time, end time.Time) (int, error)
	GetAllOrderCount(ctx context.Context) (int, error)
	GetOrderCountByStatus(ctx context.Context, status value.OrderStatus) (int, error)
}

type MetricService struct {
	repo Repo
}

func NewMetricService(repo Repo) *MetricService {
	return &MetricService{
		repo: repo,
	}
}

func (s *MetricService) GetExecutorsOrdersCountList(ctx context.Context) ([]int, error) {
	const op = "MetricService.GetExecutorsOrdersCountList"
	countList, err := s.repo.GetExecutorsOrdersCountList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slices.Sort(countList)

	return countList, nil
}

type TimeCount struct {
	Ts    time.Time `json:"ts"`
	Count int       `json:"count"`
}

func (s *MetricService) GetOrderCountByLimit(ctx context.Context, limit string) ([]TimeCount, error) {
	const op = "MetricService.GetOrderCountByLimit"

	var delay time.Duration
	var start time.Time

	now := time.Now()

	switch limit {
	case "hour":
		delay = time.Minute
		start = now.Add(-time.Hour)
	case "day":
		delay = time.Hour
		start = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Add(-time.Hour * 24)
	case "week":
		delay = time.Hour * 24
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(-time.Hour * 24 * 7)
	default: // hour
		delay = time.Minute
		start = now.Add(-time.Hour)
	}

	contextx.GetLoggerOrDefault(ctx).InfoContext(ctx, op, logx.Stringer("delay", delay), logx.Stringer("start", start))

	var timetable []TimeCount

	for ts := start; ts.Before(now); ts = ts.Add(delay) {
		endPeriod := ts.Add(delay)
		count, err := s.repo.GetCountByPeriod(ctx, ts, endPeriod)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		timetable = append(timetable, TimeCount{
			Ts:    ts,
			Count: count,
		})
	}

	return timetable, nil
}

func (s *MetricService) GetAllOrderCount(ctx context.Context) (int, error) {
	const op = "MetricService.GetAllOrderCount"
	count, err := s.repo.GetAllOrderCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return count, nil
}

func (s *MetricService) GetCompleteOrderCount(ctx context.Context) (int, error) {
	const op = "MetricService.GetCompleteOrderCount"
	countAccept, err := s.repo.GetOrderCountByStatus(ctx, value.OrderStatusAccept)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	countReject, err := s.repo.GetOrderCountByStatus(ctx, value.OrderStatusReject)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return countAccept + countReject, nil
}
