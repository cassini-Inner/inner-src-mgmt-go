package impl

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type JobsRepoImpl struct {
	db *sqlx.DB
}

func NewJobsRepo(db *sqlx.DB) *JobsRepoImpl {
	return &JobsRepoImpl{db: db}
}

func (j *JobsRepoImpl) CreateJob(ctx context.Context, tx *sqlx.Tx, input *gqlmodel.CreateJobInput, user *dbmodel.User) (*dbmodel.Job, error) {
	var insertedJob dbmodel.Job
	// insert the information into the job table
	err := tx.QueryRowxContext(ctx, createJobQuery, input.Title, input.Desc, input.Difficulty, user.Id).StructScan(&insertedJob)

	if err != nil {
		return nil, err
	}
	return j.GetByIdTx(insertedJob.Id, tx)
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
func (j *JobsRepoImpl) GetById(db sqlx.Ext, jobId string) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := db.QueryRowx(selectJobByIdQuery, jobId).StructScan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *JobsRepoImpl) GetByIdTx(jobId string, tx *sqlx.Tx) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := tx.QueryRowx(selectJobByIdQuery, jobId).StructScan(&job)
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
func (j *JobsRepoImpl) GetAll(tx sqlx.Ext, skillNames []string, status []string) ([]dbmodel.Job, error) {
	query, args, err := sqlx.In(selectAllJobsWithFiltersQuery, skillNames, status)
	if err != nil {
		return nil, err
	}
	query = j.db.Rebind(query)

	rows, err := j.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var result []dbmodel.Job
	for rows != nil && rows.Next() {
		var tempJob dbmodel.Job
		rows.StructScan(&tempJob)
		result = append(result, tempJob)
	}
	return result, nil
}

func (j *JobsRepoImpl) GetMilestonesByJobId(tx sqlx.Ext, jobId string) ([]*dbmodel.Milestone, error) {
	rows, err := tx.Queryx(selectMilestonesByJobId, jobId)
	if err != nil {
		return nil, err
	}

	var milestones []*dbmodel.Milestone
	for rows.Next() {
		var milestone dbmodel.Milestone
		rows.StructScan(&milestone)
		milestones = append(milestones, &milestone)
	}
	return milestones, nil
}

func (j *JobsRepoImpl) GetMilestoneIdsByJobId(tx sqlx.Ext, jobId string) (result []string, err error) {
	milestones, err := j.GetMilestonesByJobId(tx, jobId)
	if err != nil {
		return nil, err
	}

	for _, milestone := range milestones {
		result = append(result, milestone.Id)
	}

	return result, nil
}

func (j *JobsRepoImpl) GetMilestoneById(milestoneId string) (*dbmodel.Milestone, error) {
	var milestone dbmodel.Milestone
	err := j.db.QueryRowx(selectMilestoneByIdQuery, milestoneId).StructScan(&milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

func (j *JobsRepoImpl) GetAuthorFromMilestoneId(milestoneId string) (*dbmodel.User, error) {
	var user dbmodel.User
	err := j.db.QueryRowx(selectUserByMilestoneIdQuery, milestoneId).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (j *JobsRepoImpl) MarkJobCompleted(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error) {
	// mark the job status as completed
	_, err := tx.ExecContext(ctx, updateJobStatusCompleted, jobId)
	if err != nil {
		return nil, err
	}
	// commit the transaction
	return j.GetById(tx, jobId)
}

func (j *JobsRepoImpl) ForceAutoUpdateJobStatus(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error) {
	_, err := tx.ExecContext(ctx, updateJobStatusByIdForce, jobId)
	if err != nil {
		return nil, err
	}

	return j.GetByIdTx(jobId, tx)
}

func (j *JobsRepoImpl) ForceAutoUpdateMilestoneStatusByJobID(ctx context.Context, tx *sqlx.Tx, jobId string) error {
	_, err := tx.ExecContext(ctx, updateMilestoneStatusByJobIdForce, jobId)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobsRepoImpl) ForceAutoUpdateMilestoneStatusByMilestoneId(ctx context.Context, tx *sqlx.Tx, milestoneID string) error {
	_, err := tx.ExecContext(ctx, updateMilestoneStatusByMilestoneIDForce, milestoneID)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobsRepoImpl) MarkMilestonesCompleted(tx *sqlx.Tx, ctx context.Context, milestoneIds ...string) error {
	stmt, args, err := sqlx.In(updateMilestoneStatusCompleted, milestoneIds)
	if err != nil {
		return err
	}

	stmt = tx.Rebind(stmt)
	_, err = tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		return nil
	}
	return nil
}

func (j *JobsRepoImpl) CreateMilestones(ctx context.Context, tx *sqlx.Tx, jobId string, milestones []*gqlmodel.MilestoneInput) (createdMilestones []*dbmodel.Milestone, err error) {
	stmt, valueArgs := getInsertMilestonesStatement(milestones, jobId)
	stmt = tx.Rebind(stmt)
	// get the ids of newly inserted milestones
	milestonesInsertResult, err := tx.QueryxContext(ctx, stmt, valueArgs...)
	if err != nil {
		return nil, err
	}
	for milestonesInsertResult.Next() {
		var tempMilestone dbmodel.Milestone
		err := milestonesInsertResult.StructScan(&tempMilestone)
		if err != nil {
			return nil, err
		}
		createdMilestones = append(createdMilestones, &tempMilestone)
	}

	milestonesInsertResult.Close()
	return createdMilestones, nil
}

func (j *JobsRepoImpl) DeleteMilestonesByJobId(tx *sqlx.Tx, jobID string) error {
	_, err := tx.Exec(deleteMilestonesByJobId, jobID)
	if err != nil {
		return err
	}
	return nil
}

const (
	createJobQuery = `insert into jobs(title, description, difficulty,created_by) values ($1, $2, $3, $4) returning jobs.id`

	selectJobByIdQuery            = `SELECT * FROM jobs WHERE id = $1 and is_deleted=false`
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
