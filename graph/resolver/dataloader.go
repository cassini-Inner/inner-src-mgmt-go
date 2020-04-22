package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

const userLoaderKey = "userloader"
const milestoneByJobIdLoaderKey = "milestoneByIdLoader"

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := createUserLoader(db)
		milestoneByJobIdLoader := createMilestoneByIdLoader(db)
		ctxWithuserloader := context.WithValue(r.Context(), userLoaderKey, userloader)
		ctxWithmilestonebyidloader := context.WithValue(ctxWithuserloader, milestoneByJobIdLoaderKey, milestoneByJobIdLoader)
		next.ServeHTTP(w, r.WithContext(ctxWithmilestonebyidloader))
	})
}

func createMilestoneByIdLoader(db *sqlx.DB) *MilestoneByJobIdLoader {
	return NewMilestoneByJobIdLoader(MilestoneByJobIdLoaderConfig{
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

func getUserLoader(ctx context.Context) *UserLoader {
	userLoader := ctx.Value(userLoaderKey).(*UserLoader)
	return userLoader
}
func getMilestoneByJobIdLoader(ctx context.Context) *MilestoneByJobIdLoader {
	milestoneByJobIdLoader := ctx.Value(milestoneByJobIdLoaderKey).(*MilestoneByJobIdLoader)
	return milestoneByJobIdLoader
}
func createUserLoader(db *sqlx.DB) *UserLoader {
	return NewUserLoader(UserLoaderConfig{
		Fetch: func(keys []string) ([]*gqlmodel.User, []error) {
			users := make(map[string]dbmodel.User)
			var result []*gqlmodel.User
			query, args, err := sqlx.In(`select * from users where users.id in (?) and users.is_deleted = false`, keys)
			if err != nil {
				return nil, []error{err}
			}

			query = db.Rebind(query)
			rows, err := db.Queryx(query, args...)
			if err != nil {
				return nil, []error{err}
			}

			for rows.Next() {
				var tempUser dbmodel.User
				err := rows.StructScan(&tempUser)
				if err != nil {
					return nil, []error{err}
				}
				users[tempUser.Id] = tempUser
			}

			for _, id := range keys {
				var tempGqlUser gqlmodel.User
				tempGqlUser.MapDbToGql(users[id])
				result = append(result, &tempGqlUser)
			}

			return result, nil
		},
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})

}
