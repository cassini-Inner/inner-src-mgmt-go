package postgres

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v9"
)

// DBLogger is used for logging database calls
type DBLogger struct{}

func (d DBLogger) BeforeQuery(ctx context.Context, q *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}
func (d DBLogger) AfterQuery(ctx context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func New(options *pg.Options) *pg.DB {
	return pg.Connect(options)
}
