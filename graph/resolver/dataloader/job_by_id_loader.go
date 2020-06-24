package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewJobByIdLoader(db *sqlx.DB) *generated.JobByIdLoader {
	return generated.NewJobByIdLoader(generated.JobByIdLoaderConfig{
		Fetch: func(keys []string) ([]*gqlmodel.Job, []error) {
			resultMap := make(map[string]*gqlmodel.Job)

			var jobs []*gqlmodel.Job

			query, args, err := sqlx.In(`SELECT * FROM jobs WHERE id in (?) and is_deleted = false`, keys)
			if err != nil {
				return nil, []error{err}
			}
			query = db.Rebind(query)

			resultChan := make(chan *FetchStruct)
			go func(result chan *FetchStruct) {
				rows, err := db.Queryx(query, args...)
				result <- &FetchStruct{
					rows: rows,
					err:  err,
				}
			}(resultChan)
			res := <-resultChan

			if res.err != nil {
				return nil, []error{err}
			}
			defer res.rows.Close()

			if res.err != nil {
				return nil, []error{err}
			}
			defer res.rows.Close()
			for res.rows.Next() {
				var tempJob dbmodel.Job
				err := res.rows.StructScan(&tempJob)
				if err != nil {
					return nil, []error{err}
				}
				var gqlJob gqlmodel.Job
				gqlJob.MapDbToGql(tempJob)

				resultMap[tempJob.Id] = &gqlJob
			}

			for _, key := range keys {
				tempJob := resultMap[key]
				jobs = append(jobs, tempJob)
			}
			return jobs, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})
}
