package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	CommitTx(ctx context.Context, tx *sqlx.Tx) error
}
