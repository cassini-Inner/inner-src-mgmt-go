package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ApplicationsRepoImpl struct {
	db *sqlx.DB
}

var (
	ErrNoExistingApplications = errors.New("user does not have any existing pending or accepted applications")
)

func NewApplicationsRepo(db *sqlx.DB) *ApplicationsRepoImpl {
	return &ApplicationsRepoImpl{db: db}
}

func (a *ApplicationsRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return a.db.BeginTxx(ctx, nil)
}

func (a *ApplicationsRepoImpl) CommitTx(ctx context.Context, tx *sqlx.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}

// GetExistingUserApplications return existing user applications on the basis a applicationStatus filter
// if the number of applications is equal to the number of milestones then the user has properly applied to
// all the milestones. Returns ErrNoExistingApplications if this is not the case

// will get a null list in other case
func (a *ApplicationsRepoImpl) GetExistingUserApplications(tx *sqlx.Tx, milestones []*dbmodel.Milestone, userId string, applicationStatus ...string) ([]*dbmodel.Application, error) {

	// get the list of milestoneIds from the job milestones
	var milestoneIds []string
	for _, milestone := range milestones {
		milestoneIds = append(milestoneIds, milestone.Id)
	}

	// check if the user already has applied to a given job and the application is in pending state
	// if they have, then just return their current applications
	var statuses []string
	for _, status := range applicationStatus {
		statuses = append(statuses, fmt.Sprintf("status = '%v'", status))
	}
	stmt := fmt.Sprintf(`select * from applications where applicant_id = ? and milestone_id in (?) and (%v)`, strings.Join(statuses, " or "))
	stmt, args, err := sqlx.In(stmt, userId, milestoneIds)
	if err != nil {
		return nil, err
	}
	stmt = tx.Rebind(stmt)
	var existingApplications []*dbmodel.Application
	rows, err := tx.Queryx(stmt, args...)
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

	return nil, custom_errors.ErrNoExistingApplications
}

func (a *ApplicationsRepoImpl) CreateApplication(ctx context.Context, tx *sqlx.Tx, milestones []*dbmodel.Milestone, userId string) ([]*dbmodel.Application, error) {
	var applicationInsertArgs []interface{}
	var applicationInsertValues []string

	for _, milestone := range milestones {
		applicationInsertValues = append(applicationInsertValues, "(?, ?)")
		applicationInsertArgs = append(applicationInsertArgs, milestone.Id, userId)
	}

	stmt := fmt.Sprintf(`insert into applications(milestone_id, applicant_id) values %s returning *`, strings.Join(applicationInsertValues, ","))

	stmt = tx.Rebind(stmt)
	insertRows, err := tx.QueryxContext(ctx, stmt, applicationInsertArgs...)
	if err != nil {
		return nil, err
	}

	insertedApplications, err := scanApplicationRowsx(insertRows)
	if err != nil {
		return nil, err
	}
	insertRows.Close()
	return insertedApplications, nil
}

func (a *ApplicationsRepoImpl) SetApplicationStatusForUserMilestone(tx *sqlx.Tx, milestoneIds []string, userId string, applicationStatus string, note string) ([]*dbmodel.Application, error) {
	updateExistingStatement, updateExistingArgs, err := sqlx.In(updateApplicationsForMilestonesUser, applicationStatus, note, time.Now(), milestoneIds, userId)
	if err != nil {
		return nil, err
	}
	updateExistingStatement = tx.Rebind(updateExistingStatement)
	rows, err := tx.Queryx(updateExistingStatement, updateExistingArgs...)
	if err != nil {
		return nil, err
	}
	result, err := scanApplicationRowsx(rows)
	rows.Close()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ApplicationsRepoImpl) GetByJobId(jobId string) ([]*dbmodel.Application, error) {
	rows, err := a.db.Queryx(selectApplicationsForJobIDQuery, jobId)
	if err != nil {
		return nil, err
	}

	return scanApplicationRowsx(rows)
}

func (a *ApplicationsRepoImpl) GetApplicationStatusForUserAndJob(userId string, tx *sqlx.Tx, jobId string) (string, error) {
	// TODO: Will need to refactor this when we allow users to apply to milestones
	result := ""
	err := tx.QueryRowx(selectApplicationStatusByUserIdJobId, jobId, userId).Scan(&result)
	if err != nil {
		return "", err
	}
	return strings.ToLower(result), nil
}

func (a *ApplicationsRepoImpl) SetApplicationStatusForUserAndJob(ctx context.Context, tx *sqlx.Tx, milestones []*dbmodel.Milestone, applicationStatus string, note *string, jobId, userId string) ([]*dbmodel.Application, error) {
	var milestoneIds []string
	for _, milestone := range milestones {
			milestoneIds = append(milestoneIds, milestone.Id)
	}
	updateApplicationsQuery, updateApplicationArgs, err := sqlx.In(updateApplicationsForMilestonesUser, applicationStatus, note, time.Now(), milestoneIds, userId)
	if err != nil {
		return nil, err
	}
	updateApplicationsQuery = tx.Rebind(updateApplicationsQuery)
	rows, err := tx.Query(updateApplicationsQuery, updateApplicationArgs...)
	if err != nil {
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
		return nil, err
	}
	_, err = tx.Exec(updateJobStatusById, jobId)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(updateMilestoneStatusByJobId, jobId)
	if err != nil {
		return nil, err
	}
	return updatedApplication, nil
}

func (a *ApplicationsRepoImpl) GetAcceptedApplicationsByJobId(jobId string) ([]*dbmodel.Application, error) {
	rows, err := a.db.Queryx(selectAcceptedApplicationsForJobIDQuery, &jobId)
	if err != nil {
		return nil, err
	}
	return scanApplicationRowsx(rows)
}

func (a *ApplicationsRepoImpl) GetUserJobApplications(userId string) ([]*dbmodel.Job, error) {
	rows, err := a.db.Queryx(selectAppliedJobsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Job
	for rows != nil && rows.Next() {
		job := &dbmodel.Job{}
		rows.StructScan(job)
		result = append(result, job)
	}
	return result, nil
}

func (a *ApplicationsRepoImpl) DeleteAllJobApplications(tx *sqlx.Tx, jobId string) error {
	_, err := tx.Exec(deleteJobApplicationsByJobId, jobId)
	if err != nil {
		return err
	}
	return nil
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

func scanApplicationRowsx(rows *sqlx.Rows) ([]*dbmodel.Application, error) {
	var result []*dbmodel.Application
	for rows != nil && rows.Next() {
		application := &dbmodel.Application{}
		err := rows.StructScan(application)
		if err != nil {
			return nil, err
		}
		result = append(result, application)
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
where milestones.job_id = $1 and applications.applicant_id = $2 and milestones.is_deleted = false limit 1`

	selectAcceptedApplicationsForJobIDQuery = `select applications.id, 
		applications.milestone_id, 
		applications.applicant_id, 
		applications.status, 
		applications.note,
		applications.time_created, 
		applications.time_updated
		from applications
		join milestones on milestones.id = applications.milestone_id and milestones.is_deleted = false
		where milestones.job_id = $1 and milestones.is_delete = false and applications.status in ('pending', 'accepted' )`

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

	updateJobStatusByIdForce = `with acceptedCount as (select distinct count(distinct applicant_id) as count
	from milestones join applications a on milestones.id = a.milestone_id  where milestones.job_id = $1 and milestones.is_deleted = false and a.status = 'accepted' group by applicant_id)
	update jobs
	set status = case
					 when ((select count from acceptedCount) is not null)
						 then 'ongoing'
					 when ((select count from acceptedCount) is null )
						 then 'open'
					 else status
		end
	where jobs.id = $1 and jobs.is_deleted = false`

	updateMilestoneStatusByJobIdForce = `with acceptedCount as (select distinct count(distinct applicant_id) as count
							   from milestones
										join applications a on milestones.id = a.milestone_id and milestones.is_deleted = false
							   where milestones.job_id = $1
								 and a.status = 'accepted'
							   group by applicant_id)
		update milestones
		set status = case
						 when ((select count from acceptedCount) is not null)
							 then 'ongoing'
						 when ((select count from acceptedCount) is null)
							 then 'open'
						 else status
			end
		where milestones.job_id = $1 and milestones.is_deleted = false`

	updateMilestoneStatusByMilestoneIDForce = `with acceptedCount as (select distinct count(distinct applicant_id) as count
							   from milestones
										join applications a on milestones.id = a.milestone_id and milestones.is_deleted = false
							   where milestones.id = $1
								 and a.status = 'accepted'
							   group by applicant_id)
		update milestones
		set status = case
						 when ((select count from acceptedCount) is not null)
							 then 'ongoing'
						 when ((select count from acceptedCount) is null)
							 then 'open'
						 else status
			end
		where milestones.id = $1 and milestones.is_deleted = false`

	deleteJobApplicationsByJobId = `delete from applications
where id in (select applications.id
             from applications
                      join milestones m on applications.milestone_id = m.id
             where m.job_id = $1)`
)
