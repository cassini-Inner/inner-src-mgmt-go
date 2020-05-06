package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewMilestoneByJobIdLoader(db *sqlx.DB) *generated.MilestoneByJobIdLoader {
	return generated.NewMilestoneByJobIdLoader(generated.MilestoneByJobIdLoaderConfig{
		Fetch: func(keys []string) ([]*gqlmodel.Milestones, []error) {
			resultMap := make(map[string][]*gqlmodel.Milestone)

			var milestones []*gqlmodel.Milestones

			query, args, err := sqlx.In(`SELECT * FROM milestones WHERE job_id in (?) and is_deleted = false`, keys)
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
				var tempMilestone dbmodel.Milestone
				err := res.rows.StructScan(&tempMilestone)
				if err != nil {
					return nil, []error{err}
				}
				var gqlMilestone gqlmodel.Milestone
				gqlMilestone.MapDbToGql(tempMilestone)

				_, ok := resultMap[tempMilestone.JobId]
				if !ok {
					resultMap[tempMilestone.JobId] = make([]*gqlmodel.Milestone, 0)
				}
				resultMap[tempMilestone.JobId] = append(resultMap[tempMilestone.JobId], &gqlMilestone)
			}

			for _, key := range keys {
				tempJob := resultMap[key]
				length := len(tempJob)
				milestones = append(milestones, &gqlmodel.Milestones{
					TotalCount: &length ,
					Milestones: tempJob,
				})
			}
			return milestones, nil
		},
		Wait:     2 * time.Millisecond,
		MaxBatch: 500,
	})
}
