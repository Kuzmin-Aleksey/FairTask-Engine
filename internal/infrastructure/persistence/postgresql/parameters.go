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

	res, err := r.db.NamedExecContext(ctx, "INSERT INTO parameters (name, type) VALUES (:name, :type)", param)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	id, _ := res.LastInsertId()
	param.Id = int(id)

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
