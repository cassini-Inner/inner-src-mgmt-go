package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ApplicationsRepo struct {
	db *sqlx.DB
}

func NewApplicationsRepo(db *sqlx.DB) *ApplicationsRepo {
	return &ApplicationsRepo{db: db}
}

func (a *ApplicationsRepo) CreateApplication(milestones []*dbmodel.Milestone, userId string, ctx context.Context) ([]*dbmodel.Application, error) {

	// get the list of milestoneIds from the job milestones
	var milestoneIds []string
	for _, milestone := range milestones {
		milestoneIds = append(milestoneIds, milestone.Id)
	}

	// check if the user already has applied to a given job and the application is in pending state
	// if they have, then just return their current applications
	stmt, args, err := sqlx.In(`select * from applications where applicant_id = ? and milestone_id in (?) and (status='pending' or status='accepted')`, userId, milestoneIds)
	if err != nil {
		return nil, err
	}
	stmt = a.db.Rebind(stmt)
	var existingPendingApplications []*dbmodel.Application
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
		existingPendingApplications = append(existingPendingApplications, &application)
	}

	if len(existingPendingApplications) == len(milestones) {
		return existingPendingApplications, nil
	}

	// if the user previously applied to a job but got rejected or the application was withdrawn
	// set the status of those applications as 'pending' so they will act as new application
	stmt, args, err = sqlx.In(`select * from applications where applicant_id = ? and milestone_id in (?) and (status='rejected' or status = 'withdrawn')`, userId, milestoneIds)
	if err != nil {
		return nil, err
	}
	stmt = a.db.Rebind(stmt)
	var existingApplications []*dbmodel.Application
	existingApplicationRows, err := a.db.Queryx(stmt, args...)
	if err != nil {
		return nil, err
	}
	for existingApplicationRows.Next() {
		var application dbmodel.Application
		err := existingApplicationRows.StructScan(&application)
		if err != nil {
			return nil, err
		}
		existingApplications = append(existingApplications, &application)
	}

	// the number of existing applications of that user is same as number of milestones then
	// update the exsting applications of users and return that
	if len(existingApplications) == len(milestones) {
		updateExistingStatement, updateExistingArgs, err := sqlx.In(updateApplicationsForMilestonesUser, dbmodel.ApplicationStatusPending, "", time.Now(), milestoneIds, userId)
		if err != nil {
			return nil, err
		}
		updateExistingStatement = a.db.Rebind(updateExistingStatement)
		rows, err := a.db.Queryx(updateExistingStatement, updateExistingArgs...)
		if err != nil {
			return nil, err
		}
		result, err := scanApplicationRowsx(rows)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	// if the user has not applied to a job already then begin a new transaction to create new applications
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

	insertedApplications, err := scanApplicationRowsById(insertRows)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
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

	// once the applications have been modified, update the status of the job
	// and milestones automatically
	// if active applicants = 0 -> job is open
	// if active applicants > 0 -> job is ongoing
	// job can only by marked completed once all the milestones are resolved

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func scanApplicationRowsById(rows *sql.Rows) (result []string, err error) {
	for rows.Next() {
		id := ""
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
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

	return scanApplicationRowsx(rows)
}

func scanApplicationRowsx(rows *sqlx.Rows) ([]*dbmodel.Application, error) {
	var result []*dbmodel.Application
	for rows != nil && rows.Next() {
		var application dbmodel.Application
		err := rows.StructScan(&application)
		if err != nil {
			return nil, err
		}
		result = append(result, &application)
	}

	return result, nil
}

func scanApplicationRows(rows *sql.Rows) ([]*dbmodel.Application, error) {
	var result []*dbmodel.Application
	for rows != nil && rows.Next() {
		var application dbmodel.Application
		err := rows.Scan(&application.Id, &application.MilestoneId, &application.ApplicantId, &application.Status, &application.Note, &application.TimeCreated, &application.TimeUpdated)
		if err != nil {
			return nil, err
		}
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

func (a *ApplicationsRepo) GetApplicationStatusForUserAndJob(userId, jobId string) (string, error) {
	// TODO: Will need to refactor this when we allow users to apply to milestones
	result := ""
	err := a.db.QueryRowx(selectApplicationStatusByUserIdJobId, jobId, userId).Scan(&result)
	if err != nil {
		return "", err
	}

	return strings.ToLower(result), nil
}

func (a *ApplicationsRepo) SetApplicationStatusForUserAndJob(userId, jobId string, milestones []*dbmodel.Milestone, applicationStatus string, note *string) ([]*dbmodel.Application, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}

	var milestoneIds []string
	for _, milestone := range milestones {
		milestoneIds = append(milestoneIds, milestone.Id)
	}

	updateApplicationsQuery, updateApplicationArgs, err := sqlx.In(updateApplicationsForMilestonesUser, applicationStatus, note, time.Now(), milestoneIds, userId)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	updateApplicationsQuery = a.db.Rebind(updateApplicationsQuery)
	rows, err := tx.Query(updateApplicationsQuery, updateApplicationArgs...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	updatedApplication, err := scanApplicationRows(rows)
	if err != nil {
		return nil, err
	}

	// after updating the applications we need to update the job status too
	// if the working count of users is 0 then job is open, other wise ongoing
	workingUsersCount := 0
	err = tx.QueryRow(selectNumberOfAcceptedApplicantsForJob, jobId).Scan(&workingUsersCount)
	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(updateJobStatusById, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(updateMilestoneStatusByJobId, jobId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return updatedApplication, nil
}

func (a *ApplicationsRepo) GetAcceptedApplicationsByJobId(jobId string) ([]*dbmodel.Application, error) {
	rows, err := a.db.Queryx(selectAcceptedApplicationsForJobIDQuery, &jobId)
	if err != nil {
		return nil, err
	}
	return scanApplicationRowsx(rows)
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
		where milestones.job_id = $1 and applications.status in ('pending', 'accepted', 'rejected')`

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

	selectApplicationStatusByUserIdJobId = `select applications.status from milestones join applications on milestones.id = applications.milestone_id
where milestones.job_id = $1 and applications.applicant_id = $2 limit 1`

	selectAcceptedApplicationsForJobIDQuery = `select applications.id, 
		applications.milestone_id, 
		applications.applicant_id, 
		applications.status, 
		applications.note,
		applications.time_created, 
		applications.time_updated
		from applications
		join milestones on milestones.id = applications.milestone_id and milestones.is_deleted = false
		where milestones.job_id = $1 and applications.status in ('pending', 'accepted' )`

	updateApplicationsForMilestonesUser = `update applications set status = ?, note = ?, time_updated = ? where applications.milestone_id in (?) and applications.applicant_id = ? returning *`

	selectNumberOfAcceptedApplicantsForJob = `select distinct (applicant_id)
		from milestones
	 	join applications a on milestones.id = a.milestone_id
		where milestones.job_id = $1 and a.status = 'accepted'`

	updateJobStatusById = `with acceptedCount as (select distinct count(distinct applicant_id) as count
						   from milestones
									join applications a on milestones.id = a.milestone_id
						   where milestones.job_id = $1
							 and a.status = 'accepted'
						   group by applicant_id)
	update jobs
	set status = case
					 when ((select count from acceptedCount) is not null and jobs.status != 'completed')
						 then 'ongoing'
					 when ((select count from acceptedCount) is null and jobs.status != 'completed')
						 then 'open'
					 else status
		end
	where jobs.id = $1`

	updateMilestoneStatusByJobId = `with acceptedCount as (select distinct count(distinct applicant_id) as count
							   from milestones
										join applications a on milestones.id = a.milestone_id
							   where milestones.job_id = $1
								 and a.status = 'accepted'
							   group by applicant_id)
		update milestones
		set status = case
						 when ((select count from acceptedCount) is not null and milestones.status != 'completed')
							 then 'ongoing'
						 when ((select count from acceptedCount) is null and milestones.status != 'completed')
							 then 'open'
						 else status
			end
		where milestones.job_id = $1`
)
