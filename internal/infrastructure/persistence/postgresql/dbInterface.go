package postgresql

import (
	"database/sql"
	"golang.org/x/net/context"
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
}

type Scanner interface {
	Scan(dest ...any) error
}
