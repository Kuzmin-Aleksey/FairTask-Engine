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

	if _, err := tx.NamedExecContext(ctx, "INSERT INTO orders (id, parent_id, text, status) VALUES (:parent_id, :text, :status)", order.Order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, param := range order.Parameters {
		if _, err := tx.NamedExecContext(ctx, "INSERT INTO order_parametrs (order_id, parameter_id, value) VALUES (:order_id, :parametr_id, :value)", param); err != nil {
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

	if _, err := r.db.ExecContext(ctx, `UPDATE orders SET executor_id=? WHERE id=?`, executorId, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *OrdersRepo) GetById(ctx context.Context, id int) (*entity.Order, error) {
	const op = "OrdersRepo.GetById"

	order := new(entity.Order)

	if err := r.db.GetContext(ctx, order, "SELECT * FROM orders WHERE id=?", id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return order, nil
}

func (r *OrdersRepo) GetParameters(ctx context.Context, id int) ([]entity.Parameter, error) {
	const op = "OrdersRepo.GetParameters"

	var params []entity.Parameter

	if err := r.db.SelectContext(ctx, &params, "SELECT parameters.id, parameter.name, order_parameters.value FROM order_parameters INNER JOIN parameters ON parameters.id = order_parameters.parameter_id WHERE order_parameters.order_id=?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return params, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return params, nil
}

func (r *OrdersRepo) SetEnabled(ctx context.Context, id int, enabled bool) error {
	const op = "OrdersRepo.SetEnabled"

	if _, err := r.db.ExecContext(ctx, `UPDATE orders SET enabled=? WHERE id=?`, enabled, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OrdersRepo) Delete(ctx context.Context, id int) error {
	const op = "OrdersRepo.Delete"

	if _, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id=?`, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
