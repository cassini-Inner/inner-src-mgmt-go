package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
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
			rows, err := db.Queryx(query, args...)
			if err != nil {
				return nil, []error{err}
			}
			for rows.Next() {
				var tempMilestone dbmodel.Milestone
				err := rows.StructScan(&tempMilestone)
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
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})
}
