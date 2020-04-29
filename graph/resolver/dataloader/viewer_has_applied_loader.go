package dataloader

import (
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewViewerHasAppliedByUserIdLoader(db *sqlx.DB) *generated.ViewerHasAppliedLoader {
	return generated.NewViewerHasAppliedLoader(generated.ViewerHasAppliedLoaderConfig{
		Fetch: func(keys []string) (bools []bool, errors []error) {

			userIds := make([]string, len(keys))
			jobIds := make([]string, len(keys))

			var resultSet = make(map[string]bool)

			for i, key := range keys {
				var jobId, userId string
				_, err := fmt.Sscan(key, &jobId, &userId)
				if err != nil {
					return nil, []error{err}
				}
				userIds[i] = userId
				jobIds[i] = jobId
			}

			stmt, args, err := sqlx.In(`select job_id, applicant_id, count(distinct job_id) from milestones join applications on milestones.id = applications.milestone_id and applications.applicant_id in (?) and applications.status in ('pending', 'accepted') where job_id in (?) group by job_id, applicant_id`,userIds, jobIds )

			if err != nil {
				return nil, []error{err}
			}

			stmt = db.Rebind(stmt)

			rows, err := db.Queryx(stmt, args...)
			if err != nil {
				return nil, []error{err}
			}

			for rows.Next() {
				var jobId, userId, count int
				err := rows.Scan(&jobId, &userId, &count)
				if err != nil {
					return nil, []error{err}
				}
				result := false;
				if count > 0 {
					result = true;
				}
				key := fmt.Sprintf("%v %v", jobId, userId)
				resultSet[key] = result
			}

			for _, key := range keys {
				value, ok := resultSet[key]
				if !ok {
					value = false
				}
				bools = append(bools, value)
			}

			return bools, nil
		},
		Wait:     1 * time.Millisecond,
		MaxBatch: 200,
	})

}
