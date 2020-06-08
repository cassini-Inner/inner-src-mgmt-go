package dataloader

import (
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewJobMilestoneReviewLoader(db *sqlx.DB) *generated.JobMilestoneReviewLoader {
	return generated.NewJobMilestoneReviewLoader(generated.JobMilestoneReviewLoaderConfig{
		Fetch: func(keys []string) (i []*gqlmodel.Review, errors []error) {
			var userIds []string
			var milestoneIds []string
			resultMap := make(map[string]*gqlmodel.Review)

			for _, key := range keys {
				var milestoneId, assignedTo string
				_, err := fmt.Sscan(key, &milestoneId, &assignedTo)
				if err != nil {
					return nil, []error{err}
				}
				userIds = append(userIds, assignedTo)
				milestoneIds = append(milestoneIds, milestoneId)
			}
			stmt, args, err := sqlx.In(selectMilestoneReviewsQuery, milestoneIds, userIds)
			stmt = db.Rebind(stmt)
			if err != nil {
				return nil, []error{err}
			}
			resultChan := make(chan *FetchStruct)
			go func(result chan *FetchStruct) {
				rows, err := db.Queryx(stmt, args...)
				result <- &FetchStruct{
					rows: rows,
					err:  err,
				}
			}(resultChan)
			res := <-resultChan
			if res.err != nil {
				return nil, []error{res.err}
			}
			for res.rows.Next() {
				review := &dbmodel.Review{}
				err = res.rows.StructScan(review)
				if err != nil {
					return nil, []error{err}
				}
				var gqlReview gqlmodel.Review
				gqlReview.MapDbToGql(*review)
				resultMap[fmt.Sprintf("%v %v", review.MilestoneId, review.UserId)] = &gqlReview
			}

			for _, key := range keys {
				i = append(i, resultMap[key])
			}
			return i, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})
}

const (
	selectMilestoneReviewsQuery = "select * from reviews where milestone_id in (?) and user_id in (?)"
)
