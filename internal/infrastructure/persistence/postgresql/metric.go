package postgresql

import (
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

func (r *MetricRepo) GetCountByPeriod(ctx context.Context, start time.Time, period time.Duration) ([]int, error) {
	const op = "MetricRepo.GetCountByPeriod"
	var counts []int

	if err := r.db.SelectContext(ctx, &counts, `
WITH time_intervals AS (
    SELECT 
        generate_series(
            $2,
            DATE_TRUNC('hour', MAX(orders.ts)) + INTERVAL '1 hour',
            ($1 || ' seconds')::INTERVAL
        ) as interval_start
    FROM orders
)
SELECT 
    COUNT(orders.id) as orders_count
FROM time_intervals
LEFT JOIN orders ON 
    orders.ts >= interval_start 
    AND orders.ts < interval_start + ($1 || ' seconds')::INTERVAL
GROUP BY interval_start, (interval_start + ($1 || ' seconds')::INTERVAL)
ORDER BY (interval_start + ($1 || ' seconds')::INTERVAL);
`, period.Seconds(), start); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return counts, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return counts, nil
}
