package metric_service

import (
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"context"
	"fmt"
	"slices"
	"time"
)

type Repo interface {
	GetExecutorsOrdersCountList(ctx context.Context) ([]int, error)
	GetCountByPeriod(ctx context.Context, start time.Time, period time.Duration) ([]int, error)
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

func (s *MetricService) GetOrderCountByLimit(ctx context.Context, limit string) (map[time.Time]int, error) {
	const op = "MetricService.GetOrderCountByLimit"

	var delay time.Duration
	var start time.Time

	now := time.Now()

	switch limit {
	case "hour":
		delay = time.Minute
		start = now.Add(-time.Hour)
	case "day":
		delay = time.Hour * 4
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		delay = time.Hour * 24
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(-time.Hour * 24 * 7)
	default: // hour
		delay = time.Minute
		start = now.Add(-time.Hour)
	}

	contextx.GetLoggerOrDefault(ctx).InfoContext(ctx, op, logx.Stringer("delay", delay), logx.Stringer("start", start))

	counts, err := s.repo.GetCountByPeriod(ctx, start, delay)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	timetable := make(map[time.Time]int)

	for i, count := range counts {
		timetable[start.Add(time.Duration(i)*delay)] = count + 1

	}

	return timetable, nil
}
