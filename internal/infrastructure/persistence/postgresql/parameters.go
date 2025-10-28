package postgresql

import (
	"FairTask_Engine/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ParametersRepo struct {
	db *sqlx.DB
}

func NewParametersRepo(db *sqlx.DB) *ParametersRepo {
	return &ParametersRepo{
		db: db,
	}
}

func (r *ParametersRepo) Save(ctx context.Context, param *entity.Parameter) error {
	const op = "ParametersRepo.Save"

	if err := r.db.GetContext(ctx, &param.Id, "INSERT INTO parameters (name, type) VALUES ($1, $2) RETURNING id", param.Name, param.Type); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ParametersRepo) GetAll(ctx context.Context) ([]entity.Parameter, error) {
	const op = "ParametersRepo.GetAll"

	var params []entity.Parameter

	if err := r.db.SelectContext(ctx, &params, "SELECT * FROM parameters ORDER BY id DESC"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return params, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return params, nil
}

func (r *ParametersRepo) Delete(ctx context.Context, id int) error {
	const op = "ParametersRepo.Delete"

	if _, err := r.db.ExecContext(ctx, "DELETE FROM parameters WHERE id=$1", id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
