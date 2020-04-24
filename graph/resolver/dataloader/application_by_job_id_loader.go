package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewApplicationByJobIdLoader(db *sqlx.DB) *generated.ApplicationsByJobIdLoader {
	return generated.NewApplicationsByJobIdLoader(generated.ApplicationsByJobIdLoaderConfig{
		Fetch: func(keys []string) (i [][]*gqlmodel.Applications, errors []error) {
			return nil, nil
		},
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})
}
