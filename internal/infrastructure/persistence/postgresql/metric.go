package postgresql

import (
	"FairTask_Engine/internal/domain/value"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type MetricRepo struct {
	db *sqlx.DB
}

func NewMetricRepo(db *sqlx.DB) *MetricRepo {
	return &MetricRepo{
		db: db,
	}
}

func (r *MetricRepo) GetExecutorsOrdersCountList(ctx context.Context) ([]int, error) {
	const op = "MetricRepo.GetExecutorsOrdersCountList"

	var counts []int

	if err := r.db.SelectContext(ctx, &counts, "SELECT (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id AND orders.status = 'processed') AS order_count FROM executors WHERE executors.status = 'active'"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return counts, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return counts, nil

}

func (r *MetricRepo) GetCountByPeriod(ctx context.Context, start time.Time, end time.Time) (int, error) {
	const op = "MetricRepo.GetCountByPeriod"
	var count int

	if err := r.db.GetContext(ctx, &count, "SELECT count(*) FROM orders WHERE  $1 <= ts  AND ts < $2", start, end); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return count, fmt.Errorf("%s: %w", op, err)
		}
	}

	return count, nil
}

func (r *MetricRepo) GetAllOrderCount(ctx context.Context) (int, error) {
	const op = "MetricRepo.GetAllOrderCount"
	var count int

	if err := r.db.GetContext(ctx, &count, "SELECT count(*) FROM orders"); err != nil {
		return count, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

func (r *MetricRepo) GetOrderCountByStatus(ctx context.Context, status value.OrderStatus) (int, error) {
	const op = "MetricRepo.GetAllOrderCount"
	var count int

	if err := r.db.GetContext(ctx, &count, "SELECT count(*) FROM orders WHERE status=$1", status); err != nil {
		return count, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
