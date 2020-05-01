package dataloader

import (
	"database/sql"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
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

			acceptedApplicationsQuery, acceptedApplicationsArgs, err := sqlx.In(acceptedApplicationsQuery, keys)
			if err != nil {
				return nil, []error{err}
			}

			pendingApplicationsQuery, pendingApplicationsArgs, err := sqlx.In(pendingApplicationsQuery, keys)
			if err != nil {
				return nil, []error{err}
			}

			rejectedApplicationsQuery, rejectedApplicationsArgs, err := sqlx.In(rejectedApplicationsQuery, keys)
			if err != nil {
				return nil, []error{err}
			}

			applicationsQuery = db.Rebind(applicationsQuery)
			rejectedApplicationsQuery = db.Rebind(rejectedApplicationsQuery)
			acceptedApplicationsQuery = db.Rebind(acceptedApplicationsQuery)
			pendingApplicationsQuery = db.Rebind(pendingApplicationsQuery)

			applicationRows, err := db.Queryx(applicationsQuery, applicationsArgs...)
			if err != nil {
				return nil, []error{err}
			}
			errors = mapApplicationRowsToGqlModel(applicationRows, applicationsMap)
			if errors != nil {
				return nil, errors
			}

			acceptedCountRows, err := db.Queryx(acceptedApplicationsQuery, acceptedApplicationsArgs...)
			if err != nil {
				return nil, []error{err}
			}
			mapCountRowToValue(acceptedMap, acceptedCountRows, keys)

			rejectedCountRows, err := db.Queryx(rejectedApplicationsQuery, rejectedApplicationsArgs...)
			if err != nil {
				return nil, []error{err}
			}

			mapCountRowToValue(rejectedMap, rejectedCountRows, keys)

			pendingCountRows, err := db.Queryx(pendingApplicationsQuery, pendingApplicationsArgs...)
			if err != nil {
				return nil, []error{err}
			}
			mapCountRowToValue(pendingMap, pendingCountRows, keys)

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
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})
}

func mapCountRowToValue(storageMap map[string]int, rows *sqlx.Rows, keys []string) error {
	for _, key := range keys {
		_, ok := storageMap[key]
		if !ok {
			storageMap[key] = 0
		}
	}
	for rows.Next() {
		var jobId string
		var count int
		err := rows.Scan(&jobId, &count)
		if err != nil {
			return err
		}
		storageMap[jobId] = count
	}

	return nil
}

func mapApplicationRowsToGqlModel(applicationRows *sqlx.Rows, applicationsMap map[string][]*gqlmodel.Application) []error {
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

	acceptedApplicationsQuery = `select job_id, count(distinct  a.applicant_id) from milestones
			join applications a on milestones.id = a.milestone_id where a.status = 'accepted'
				and milestones.job_id in (?)
				group by job_id`

	rejectedApplicationsQuery = `select job_id, count(distinct  a.applicant_id) from milestones
			join applications a on milestones.id = a.milestone_id where a.status = 'rejected'
				and milestones.job_id in (?)
				group by job_id`

	pendingApplicationsQuery = `select job_id, count(distinct  a.applicant_id) from milestones
			join applications a on milestones.id = a.milestone_id where a.status = 'pending'
				and milestones.job_id in (?)
				group by job_id`
)
