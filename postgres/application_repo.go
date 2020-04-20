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
	rows, err := a.db.Queryx(selectApplicationsForJobIDQuery, jobId)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Application
	for rows != nil && rows.Next() {
		var application dbmodel.Application
		rows.StructScan(&application)
		result = append(result, &application)
	}
	return result, nil
}

func (a *ApplicationsRepo) GetUserJobApplications(userId string) ([]*dbmodel.Job, error) {

	rows, err := a.db.Queryx(selectAppliedJobsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Job
	for rows != nil && rows.Next() {
		var job dbmodel.Job
		rows.StructScan(&job)
		result = append(result, &job)
	}
	return result, nil
}

const (
	selectApplicationsForJobIDQuery = `select applications.id, 
		applications.milestone_id, 
		applications.applicant_id, 
		applications.status, 
		applications.note,
		applications.time_created, 
		applications.time_updated
		from applications
		join milestones on milestones.id = applications.milestone_id and milestones.is_deleted = false
		where milestones.job_id = $1 and applications.status in ('pending', 'accepted' )`

	selectAppliedJobsByUserIdQuery = `select distinct jobs.id,
		jobs.created_by,
		jobs.title,
		jobs.description,
		jobs.difficulty,
		jobs.status,
		jobs.time_created,
		jobs.time_updated,
		jobs.is_deleted
		from applications
		join milestones on milestones.id = applications.milestone_id and milestones.is_deleted = false
		join jobs on milestones.job_id = jobs.id and jobs.is_deleted = false
		where applicant_id = $1 and applications.status in ('pending', 'accepted', 'rejected')`
)
