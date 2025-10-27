package postgresql

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type ExecutorsRepo struct {
	db *sqlx.DB
}

func NewExecutors(db *sqlx.DB) *ExecutorsRepo {
	return &ExecutorsRepo{
		db: db,
	}
}

func (r *ExecutorsRepo) GetParameters(ctx context.Context, id int) ([]entity.Parameter, error) {
	const op = "ExecutorsRepo.GetParameters"

	var params []entity.Parameter

	if err := r.db.SelectContext(ctx, &params, "SELECT parameters.id, parameters.name, executor_params.value FROM executor_params INNER JOIN parameters ON parameters.id = executor_params.param_id WHERE executor_params.user_id = ?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return params, nil
}

func (r *ExecutorsRepo) GetActive(ctx context.Context) ([]aggregate.ExecutorWithParams, error) {
	var executors []aggregate.ExecutorWithParams

	if err := r.db.SelectContext(ctx, &executors, "SELECT executors.id, executor.name, (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id) AS order_count FROM executors WHERE executors.active"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	for i, executor := range executors {
		params, err := r.GetParameters(ctx, executor.Id)
		if err != nil {
			return nil, err
		}

		executors[i].Parameters = params
	}

	return executors, nil
}
