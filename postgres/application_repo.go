package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ApplicationsRepo struct {
	db *sqlx.DB
}

func NewApplicationsRepo(db *sqlx.DB) *ApplicationsRepo {
	return &ApplicationsRepo{db: db}
}

func (a *ApplicationsRepo) CreateApplication(milestones []*dbmodel.Milestone, userId string, ctx context.Context) ([]*dbmodel.Application, error) {

	// check if the user already has applied to a given job
	// if they have, then just return their current applications
	var milestoneIds []string
	for _, milestone := range milestones {
		milestoneIds = append(milestoneIds, milestone.Id)
	}
	stmt, args, err := sqlx.In(`select * from applications where applicant_id = ? and milestone_id in (?) and status='pending'`, userId, milestoneIds)
	if err != nil {
		return nil, err
	}
	stmt = a.db.Rebind(stmt)
	var existingApplications []*dbmodel.Application
	rows, err := a.db.Queryx(stmt, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var application dbmodel.Application
		err := rows.StructScan(&application)
		if err != nil {
			return nil, err
		}
		existingApplications = append(existingApplications, &application)
	}


	if len(existingApplications) == len(milestones) {
		return existingApplications, nil
	}

	// if the user has not applied to a job already then
	// create new applications for the user
	var result []*dbmodel.Application

	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}

	var applicationInsertArgs []interface{}
	var applicationInsertValues []string

	for _, milestone := range milestones {
		applicationInsertValues = append(applicationInsertValues, "(?, ?)")
		applicationInsertArgs = append(applicationInsertArgs, milestone.Id, userId)
	}

	stmt = fmt.Sprintf(`insert into applications(milestone_id, applicant_id) values %s returning id`, strings.Join(applicationInsertValues, ","))

	stmt = a.db.Rebind(stmt)
	insertRows, err := tx.QueryContext(ctx, stmt, applicationInsertArgs...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var insertedApplications []string
	for insertRows.Next() {
		id := ""
		err = insertRows.Scan(&id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		insertedApplications = append(insertedApplications, id)
	}
	insertRows.Close()

	stmt, args, err = sqlx.In(`select id, milestone_id, applicant_id, status, note, time_created, time_updated from applications where id in (?)`, insertedApplications)
	stmt = a.db.Rebind(stmt)
	getApplicationsRows, err := tx.Query(stmt, args...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for getApplicationsRows.Next() {
		var id, milestoneId, applicantId, status, timeCreated, timeUpdated string
		var note sql.NullString
		if err = getApplicationsRows.Scan(&id, &milestoneId, &applicantId, &status, &note, &timeCreated, &timeUpdated); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		result = append(result, &dbmodel.Application{
			Id:          id,
			MilestoneId: milestoneId,
			ApplicantId: applicantId,
			Status:      status,
			Note:        note,
			TimeCreated: timeCreated,
			TimeUpdated: timeUpdated,
		})
	}
	getApplicationsRows.Close()
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return result, nil
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
