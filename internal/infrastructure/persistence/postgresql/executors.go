package postgresql

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ExecutorsRepo struct {
	db *sqlx.DB
}

func NewExecutorsRepo(db *sqlx.DB) *ExecutorsRepo {
	return &ExecutorsRepo{
		db: db,
	}
}

func (r *ExecutorsRepo) GetParameters(ctx context.Context, id int) ([]entity.ExecutorParameter, error) {
	const op = "ExecutorsRepo.GetParameters"

	var params []entity.ExecutorParameter

	if err := r.db.SelectContext(ctx, &params, "SELECT parameters.id, executor_parameters.mask, parameters.type FROM executor_parameters INNER JOIN parameters ON parameters.id = executor_parameters.parameter_id WHERE executor_parameters.executor_id = $1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return params, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return params, nil
}

func (r *ExecutorsRepo) GetActive(ctx context.Context) ([]aggregate.ExecutorWithParams, error) {
	const op = "ExecutorsRepo.GetActive"

	var executors []aggregate.ExecutorWithParams

	if err := r.db.SelectContext(ctx, &executors, "SELECT id, name, (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id AND orders.status = 'processed') AS order_count, max_daily_limit FROM executors WHERE executors.status = 'active'"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i, executor := range executors {
		params, err := r.GetParameters(ctx, executor.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		executors[i].Parameters = params
	}

	return executors, nil
}
