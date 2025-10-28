package postgresql

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/pkg/failure"
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

	if err := r.db.SelectContext(ctx, &executors, "SELECT id, name, (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id AND orders.status = 'processed') AS order_count FROM executors WHERE executors.status = 'active'"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return executors, nil
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

func (r *ExecutorsRepo) GetById(ctx context.Context, orderId int) (*aggregate.ExecutorWithParams, error) {
	const op = "ExecutorsRepo.GetByOrderId"

	executor := &aggregate.ExecutorWithParams{}

	if err := r.db.GetContext(ctx, executor, "SELECT id, name, status, (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id AND orders.status = 'processed') AS order_count FROM executors WHERE id = $1", orderId); err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	params, err := r.GetParameters(ctx, executor.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	executor.Parameters = params

	return executor, nil
}

func (r *ExecutorsRepo) SetStatus(ctx context.Context, id int, status string) error {
	const op = "ExecutorsRepo.SetStatus"

	if _, err := r.db.ExecContext(ctx, "UPDATE executors SET status = $1 WHERE id = $2", status, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ExecutorsRepo) Create(ctx context.Context, executor *aggregate.ExecutorWithParams) error {
	const op = "ExecutorsRepo.Create"

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.GetContext(ctx, &executor.Id, "INSERT INTO executors (name, status) VALUES ($1, $2) RETURNING id", executor.Name, executor.Status); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, param := range executor.Parameters {
		if _, err := tx.ExecContext(ctx, "INSERT INTO executor_parameters (parameter_id, executor_id, mask) VALUES ($1, $2, $3)", param.Id, executor.Id, param.Mask); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ExecutorsRepo) GetAll(ctx context.Context) ([]aggregate.ExecutorWithParams, error) {
	const op = "ExecutorsRepo.GetAll"

	var executors []aggregate.ExecutorWithParams

	if err := r.db.SelectContext(ctx, &executors, "SELECT id, name, status, (SELECT count(*) FROM orders WHERE orders.executor_id = executors.id AND orders.status = 'processed') AS order_count FROM executors"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return executors, nil
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

func (r *ExecutorsRepo) Delete(ctx context.Context, id int) error {
	const op = "ExecutorsRepo.Delete"
	if _, err := r.db.ExecContext(ctx, "DELETE FROM executors WHERE id = $1", id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ExecutorsRepo) AddParameter(ctx context.Context, id int, param *entity.ExecutorParameter) error {
	const op = "ExecutorsRepo.AddParameter"
	if err := r.db.GetContext(ctx, &param.Type, "SELECT type FROM parameters WHERE id=$1", param.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return failure.NewNotFoundError(fmt.Sprintf("parameter %d not found", param.Id))
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := r.db.ExecContext(ctx, "INSERT INTO executor_parameters (parameter_id, executor_id, mask) VALUES ($1, $2, $3)", param.Id, id, param.Mask); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ExecutorsRepo) DeleteParameter(ctx context.Context, id int, paramId int) error {
	const op = "ExecutorsRepo.DeleteParameter"
	if _, err := r.db.ExecContext(ctx, "DELETE FROM executor_parameters WHERE parameter_id=$1 AND executor_id=$2", paramId, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
