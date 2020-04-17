package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

type ApplicationsRepo struct {
	db *sqlx.DB
}

func NewApplicationsRepo(db *sqlx.DB) *ApplicationsRepo {
	return &ApplicationsRepo{db: db}
}

func (a *ApplicationsRepo) CreateApplication(jobId string, userId string) (*model.Application, error) {
	panic("Not implemented")
}
func (a *ApplicationsRepo) UpdateApplication(applicantId, jobId string, newStatus model.ApplicationStatus) (*model.Application, error) {
	panic("Not implemented")
}

func (a *ApplicationsRepo) DeleteApplication(jobId string, userId string) (*model.Application, error) {
	panic("Not implemented")
}

func (a *ApplicationsRepo) GetByJobId(jobId string) ([]*dbmodel.Application, error) {
	var result []*dbmodel.Application
	var application dbmodel.Application

	rows, err := a.db.Queryx(getApplicationsForJobID, jobId)
	if err != nil {
		return nil, err
	}
	for rows != nil && rows.Next() {
		rows.StructScan(&application)
		result = append(result, &application)
	}
	return result, nil
}

func (a *ApplicationsRepo) GetUserJobApplications(userId string) ([]*dbmodel.Job, error) {
	var result[]*dbmodel.Job
	var job dbmodel.Job

	rows, err := a.db.Queryx(getAppliedJobsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}

	for rows != nil && rows.Next() {
		rows.StructScan(&job)
		result = append(result, &job)
	}
	return result, nil
}

const (
	getApplicationsForJobID = `select applications.id,
		milestones.id as "milestone_id",
		applications.applicant_id as "applicant_id",
		applications.status as "status",
		applications.note,
		applications.time_created,
		applications.time_updated
		from milestones
		join applications on applications.milestone_id = milestones.id
		join users on applications.applicant_id = users.id
		where milestones.job_id = $1`

	getAppliedJobsByUserIdQuery = `select distinct(jobs.id),
		jobs.title
		from applications
		join milestones on milestones.id = applications.milestone_id
		join jobs on milestones.job_id = jobs.id
		where applicant_id = $1 and applications.status <> 'withdrawn'`
)
