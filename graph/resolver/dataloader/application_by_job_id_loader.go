package dataloader

import (
	"database/sql"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

func NewApplicationByJobIdLoader(db *sqlx.DB) *generated.ApplicationsByJobIdLoader {
	return generated.NewApplicationsByJobIdLoader(generated.ApplicationsByJobIdLoaderConfig{
		Fetch: func(keys []string) (i []*gqlmodel.Applications, errors []error) {
			applicationsMap := make(map[string][]*gqlmodel.Application)
			acceptedMap := make(map[string]int)
			rejectedMap := make(map[string]int)
			pendingMap := make(map[string]int)

			applicationsQuery, applicationsArgs, err := sqlx.In(applicationsQuery, keys)
			if err != nil {
				return nil, []error{err}
			}

			applicationsQuery = db.Rebind(applicationsQuery)
			resultChan := make(chan *FetchStruct)

			go func(result chan *FetchStruct) {
				rows, err := db.Queryx(applicationsQuery, applicationsArgs...)
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

			errors = mapApplicationRowsToGqlModel(res.rows, applicationsMap, &acceptedMap, &rejectedMap, &pendingMap)
			if errors != nil {
				return nil, errors
			}

			for _, key := range keys {
				pendingCount := pendingMap[key]
				acceptedCount := acceptedMap[key]
				rejectedCount := rejectedMap[key]
				i = append(i, &gqlmodel.Applications{
					PendingCount:  &pendingCount,
					AcceptedCount: &acceptedCount,
					RejectedCount: &rejectedCount,
					Applications:  applicationsMap[key],
				})
			}
			return i, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})
}

func mapApplicationRowsToGqlModel(applicationRows *sqlx.Rows, applicationsMap map[string][]*gqlmodel.Application, acceptedCountMap, rejectedCountMap, pendingCountMap *map[string]int) []error {
	applicationCountsMap := make(map[string]map[string]bool)
	for applicationRows.Next() {
		var jobId, id, milestoneId, applicantId, status, timeCreated, timeUpdated string
		var note sql.NullString
		err := applicationRows.Scan(&jobId, &id, &milestoneId, &applicantId, &status, &note, &timeCreated, &timeUpdated)
		if err != nil {
			return []error{err}
		}
		_, ok := applicationsMap[jobId]
		if !ok {
			applicationsMap[jobId] = make([]*gqlmodel.Application, 0)
		}

		application := &gqlmodel.Application{
			ID:          id,
			ApplicantID: applicantId,
			MilestoneID: milestoneId,
			Status:      gqlmodel.ApplicationStatus(status),
			CreatedOn:   timeCreated,
		}

		_, ok = applicationCountsMap[jobId][applicantId]
		if !ok {
			applicationCountsMap[jobId] = make(map[string]bool)
			_, ok = applicationCountsMap[jobId][applicantId]
			if !ok {
				applicationCountsMap[jobId][applicantId] = true
				switch application.Status.String() {
				case strings.ToLower(gqlmodel.ApplicationStatusAccepted.String()):
					(*acceptedCountMap)[jobId]++
				case strings.ToLower(gqlmodel.ApplicationStatusPending.String()):
					(*pendingCountMap)[jobId]++
				case strings.ToLower(gqlmodel.ApplicationStatusRejected.String()):
					(*rejectedCountMap)[jobId]++
				}
			}
		}

		if note.Valid {
			application.Note = &note.String
		}
		applicationsMap[jobId] = append(applicationsMap[jobId], application)
	}
	return nil
}

const (
	applicationsQuery = `select m.job_id, applications.id, applications.milestone_id, applications.applicant_id, applications.status, applications.note, applications.time_created, applications.time_updated
			from applications
					 join milestones m on applications.milestone_id = m.id
			where m.job_id in (?) and applications.status <> 'withdrawn'`
)
