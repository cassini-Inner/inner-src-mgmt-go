package postgres

import (
	"github.com/jmoiron/sqlx"
)

type MilestonesRepo struct {
	db *sqlx.DB
}

func NewMilestonesRepo(db *sqlx.DB) *MilestonesRepo {
	return &MilestonesRepo{db: db}
}

