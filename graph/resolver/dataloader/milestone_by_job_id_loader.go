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
		Fetch: func(keys []string) (i [][]*gqlmodel.Milestone, errors []error) {
			var gqlMilestones []*gqlmodel.Milestone
			resultMap := make(map[string][]*gqlmodel.Milestone)

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
				gqlMilestones = append(gqlMilestones, &gqlMilestone)
			}

			for _, milestones := range gqlMilestones {
				_, ok := resultMap[milestones.JobID]
				if !ok {
					resultMap[milestones.JobID] = make([]*gqlmodel.Milestone, 0)
				}
				resultMap[milestones.JobID] = append(resultMap[milestones.JobID], milestones)
			}


			for _, key := range keys {
				tempJob := resultMap[key]
				i = append(i,tempJob)
			}
			return i, nil
		},
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})
}
