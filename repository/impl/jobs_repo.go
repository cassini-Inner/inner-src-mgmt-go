package impl

import (
	"context"
	"fmt"
	"strings"

	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type JobsRepoImpl struct {
	db *sqlx.DB
}

func NewJobsRepo(db *sqlx.DB) *JobsRepoImpl {
	return &JobsRepoImpl{db: db}
}

func (j *JobsRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return j.db.BeginTxx(ctx, nil)
}

func (j *JobsRepoImpl) CommitTx(ctx context.Context, tx *sqlx.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}

func (j *JobsRepoImpl) CreateJob(ctx context.Context, tx *sqlx.Tx, input *gqlmodel.CreateJobInput, user *dbmodel.User) (*dbmodel.Job, error) {
	var insertedJob dbmodel.Job
	// insert the information into the job table
	err := tx.QueryRowxContext(ctx, createJobQuery, input.Title, input.Desc, input.Difficulty, user.Id).StructScan(&insertedJob)

	if err != nil {
		return nil, err
	}
	return j.getByIdTx(tx, insertedJob.Id)
}

func (j *JobsRepoImpl) UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepoImpl) DeleteJob(tx *sqlx.Tx, jobId string) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := tx.QueryRowx(deleteJobQuery, jobId).StructScan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Get the complete job details based on the job id
func (j *JobsRepoImpl) GetById(jobId string) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := j.db.QueryRowx(selectJobByIdQuery, jobId).StructScan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByUserId returns all jobs created by that user
func (j *JobsRepoImpl) GetByUserId(userId string) ([]*dbmodel.Job, error) {

	rows, err := j.db.Queryx(selectJobsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}

	var jobs []*dbmodel.Job
	for rows.Next() {
		var job dbmodel.Job
		rows.StructScan(&job)
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

//TODO: Refactor this.
func (j *JobsRepoImpl) GetStatsByUserId(userId string) (*gqlmodel.UserStats, error) {
	panic("not implemented")
}

//TODO: Add sorting order functionality
func (j *JobsRepoImpl) GetAll(skillNames []string, status []string) ([]dbmodel.Job, error) {
	query, args, err := sqlx.In(selectAllJobsWithFiltersQuery, skillNames, status)
	if err != nil {
		return nil, err
	}
	query = j.db.Rebind(query)

	rows, err := j.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (j *JobsRepoImpl) GetAllPaginated(skillNames []string, status []string, limit int, cursor *string) ([]dbmodel.Job, error) {
	if cursor != nil {
		query, args, err := sqlx.In(selectAllJobsLimitedWithID, cursor, skillNames, status, limit)
		if err != nil {
			return nil, err
		}
		query = j.db.Rebind(query)

		rows, err := j.db.Queryx(query, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanRows(rows)
	}
	query, args, err := sqlx.In(selectAllJobsLimited, skillNames, status, limit)
	if err != nil {
		return nil, err
	}
	query = j.db.Rebind(query)

	rows, err := j.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (j *JobsRepoImpl) MarkJobCompleted(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error) {
	_, err := tx.ExecContext(ctx, updateJobStatusCompleted, jobId)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	return j.getByIdTx(tx, jobId)
}

func (j *JobsRepoImpl) ForceAutoUpdateJobStatus(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error) {
	_, err := tx.ExecContext(ctx, updateJobStatusByIdForce, jobId)
	if err != nil {
		return nil, err
	}

	return j.getByIdTx(tx, jobId)
}

func (j *JobsRepoImpl) GetByTitle(jobTitle string, limit *int) ([]dbmodel.Job, error) {
	rows, err := j.db.Queryx(selectJobByTitleQuery, jobTitle, limit)
	if err != nil {
		return nil, err
	}

	var jobs []dbmodel.Job
	for rows != nil && rows.Next() {
		var tempJob dbmodel.Job
		rows.StructScan(&tempJob)
		jobs = append(jobs, tempJob)
	}

	return jobs, nil
}

func (j *JobsRepoImpl) getByIdTx(tx *sqlx.Tx, id string) (result *dbmodel.Job, err error) {
	err = tx.QueryRowx(selectJobByIdQuery, id).StructScan(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func scanRows(rows *sqlx.Rows) (result []dbmodel.Job, err error) {
	for rows != nil && rows.Next() {
		var tempJob dbmodel.Job
		err = rows.StructScan(&tempJob)
		if err != nil {
			return nil, err
		}
		result = append(result, tempJob)
	}
	return result, nil
}

const (
	createJobQuery = `insert into jobs(title, description, difficulty,created_by) values ($1, $2, $3, $4) returning jobs.id`

	selectJobByIdQuery            = `SELECT * FROM jobs WHERE id = $1 and is_deleted=false`
	selectJobByTitleQuery         = `SELECT * FROM jobs WHERE title ~* $1 and is_deleted=false ORDER BY title LIMIT $2`
	selectJobsByUserIdQuery       = `SELECT * FROM jobs WHERE created_by = $1 and is_deleted=false`
	selectAllJobsWithFiltersQuery = `select * from jobs where jobs.id in (select milestones.job_id from globalskills 
		join milestoneskills on globalskills.id = milestoneskills.skill_id
		join milestones on milestones.id = milestoneskills.milestone_id
		where globalskills.value in (?)
		)
		and jobs.status in (?) and jobs.is_deleted = false`
	selectMilestonesByJobId      = `SELECT * FROM milestones WHERE job_id = $1 and is_deleted = false`
	selectMilestoneByIdQuery     = `SELECT * FROM milestones WHERE id = $1 and is_deleted=false`
	selectUserByMilestoneIdQuery = `select * from users where id in (select created_by from jobs join milestones m on 
	jobs.id = m.job_id and m.id = $1 and m.is_deleted = false and jobs.is_deleted = false)`

	// TODO: These 2 queries can be optimised further. Need feedback on this
	selectAllJobsLimited = `select *
from jobs
where id in (select distinct job_id
             from milestones
                      join milestoneskills on milestones.id = milestoneskills.milestone_id and milestoneskills.is_deleted = false and milestones.is_deleted = false
                      join globalskills g on milestoneskills.skill_id = g.id
             where value in (?)
             )
and status in (?)
and jobs.is_deleted = false
order by jobs.time_created desc
fetch first ? rows only`

	selectAllJobsLimitedWithID = `select *
from jobs
where id in (select distinct job_id
             from milestones
                      join milestoneskills on milestones.id = milestoneskills.milestone_id and job_id < ? and milestoneskills.is_deleted = false and milestones.is_deleted = false
                      join globalskills g on milestoneskills.skill_id = g.id
             where value in (?)
             )
and status in (?)
and jobs.is_deleted = false
order by jobs.time_created desc
fetch first ? rows only`

	updateMilestoneStatusCompleted = `update milestones set status = 'completed' where id in (?) and is_deleted = false`
	updateJobStatusCompleted       = `update jobs set status = 'completed' where id = $1 and is_deleted = false`
	deleteMilestonesByJobId        = `update milestones set is_deleted = true where job_id = $1`
	deleteJobQuery                 = `update jobs set is_deleted = true where id = $1 returning *`
)

func getInsertMilestonesStatement(milestoneInputs []*gqlmodel.MilestoneInput, insertedJobId string) (string, []interface{}) {
	var valueStrings []string
	var valueArgs []interface{}
	for _, milestone := range milestoneInputs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, milestone.Title, milestone.Desc, milestone.Duration, milestone.Status.String(), milestone.Resolution, insertedJobId)
	}
	stmt := fmt.Sprintf("INSERT INTO milestones (title, description, duration, status, resolution, job_id) VALUES %s returning *",
		strings.Join(valueStrings, ", "))
	return stmt, valueArgs
}
