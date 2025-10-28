package postgresql

import (
	"FairTask_Engine/internal/domain/aggregate"
	"FairTask_Engine/internal/domain/entity"
	"FairTask_Engine/internal/domain/value"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{
		db: db,
	}
}

func (r *OrdersRepo) Save(ctx context.Context, order *aggregate.OrderWithParameter) error {
	const op = "OrdersRepo.Save"

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := tx.NamedExecContext(ctx, "INSERT INTO orders (id, parent_id, text, status, executor_id) VALUES (:id, :parent_id, :text, :status, :executor_id)", order.Order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, param := range order.Parameters {
		if _, err := tx.ExecContext(ctx, "INSERT INTO order_parameters (order_id, parameter_id, value) VALUES ($1, $2, $3)", order.Id, param.Id, param.Value); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrdersRepo) UpdateExecutor(ctx context.Context, id int, executorId int) error {
	const op = "OrdersRepo.UpdateExecutor"

	if _, err := r.db.ExecContext(ctx, `UPDATE orders SET executor_id=$1 WHERE id=$2`, executorId, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *OrdersRepo) GetById(ctx context.Context, id int) (*entity.Order, error) {
	const op = "OrdersRepo.GetById"

	order := new(entity.Order)

	if err := r.db.GetContext(ctx, order, "SELECT id, parent_id, text, status, executor_id FROM orders WHERE id=$1", id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (r *OrdersRepo) GetParameters(ctx context.Context, id int) ([]entity.OrderParameter, error) {
	const op = "OrdersRepo.GetParameters"

	var params []entity.OrderParameter

	if err := r.db.SelectContext(ctx, &params, "SELECT parameter_id AS id, value FROM order_parameters WHERE order_parameters.order_id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return params, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return params, nil
}

func (r *OrdersRepo) GetWithoutExecutor(ctx context.Context) ([]aggregate.OrderWithParameter, error) {
	const op = "OrdersRepo.GetWithoutExecutor"
	var orders []aggregate.OrderWithParameter

	if err := r.db.SelectContext(ctx, &orders, "SELECT id, parent_id, text, status, executor_id FROM orders WHERE executor_id=0 AND status='processed'"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i, order := range orders {
		params, err := r.GetParameters(ctx, order.Id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		orders[i].Parameters = params
	}

	return orders, nil
}

func (r *OrdersRepo) SetStatus(ctx context.Context, id int, status value.OrderStatus) error {
	const op = "OrdersRepo.SetEnabled"

	if _, err := r.db.ExecContext(ctx, `UPDATE orders SET status=$1 WHERE id=$2`, status, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrdersRepo) Delete(ctx context.Context, id int) error {
	const op = "OrdersRepo.Delete"

	if _, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id=$1`, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
