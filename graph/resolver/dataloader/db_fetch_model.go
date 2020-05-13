package dataloader

import "github.com/jmoiron/sqlx"

type FetchStruct struct {
	rows *sqlx.Rows
	err  error
}
