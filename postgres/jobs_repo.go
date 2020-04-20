package postgres

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type JobsRepo struct {
	db *sqlx.DB
}

func NewJobsRepo(db *sqlx.DB) *JobsRepo {
	return &JobsRepo{db: db}
}

func (j *JobsRepo) CreateJob(input *gqlmodel.CreateJobInput) (*dbmodel.Job, error) {
	// query := "INSERT INTO jobs (CreatedBy, Title, Description, Difficulty, Status, TimeCreated, TimeUpdated, IsDeleted) VALUES ("
	// err := j.db.QueryRowx()
	panic("Not implemented")
}

func (j *JobsRepo) UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) DeleteJob(jobId string) (*dbmodel.Job, error) {
	panic("Not implemented")
}

// Get the complete job details based on the job id
func (j *JobsRepo) GetById(jobId string) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := j.db.QueryRowx(selectJobByIdQuery, jobId).StructScan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByUserId returns all jobs created by that user
func (j *JobsRepo) GetByUserId(userId string) ([]*dbmodel.Job, error) {

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
func (j *JobsRepo) GetStatsByUserId(userId string) (*gqlmodel.UserStats, error) {
	panic("not implemented")
}

//TODO: Add sorting order functionality
func (j *JobsRepo) GetAll(filters *gqlmodel.JobsFilterInput) ([]*dbmodel.Job, error) {
	var jobSkills []string
	for _, skill := range filters.Skills {
		jobSkills = append(jobSkills, strings.ToLower(*skill))
	}

	var jobStatuses []string
	for _, status := range filters.Status {
		jobStatuses = append(jobStatuses, strings.ToLower(status.String()))
	}

	query, args, err := sqlx.In(selectAllJobsWithFiltersQuery, jobSkills, jobStatuses)
	if err != nil {
		return nil, err
	}
	query = j.db.Rebind(query)

	rows, err := j.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Job
	for rows != nil && rows.Next() {
		var tempJob dbmodel.Job
		rows.StructScan(&tempJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

const (
	selectJobByIdQuery            = `SELECT * FROM jobs WHERE id = $1 and is_deleted=false`
	selectJobsByUserIdQuery       = `SELECT * FROM jobs WHERE created_by = $1 and is_deleted=false`
	selectAllJobsWithFiltersQuery = `select * from jobs where jobs.id in (
		select milestones.job_id from globalskills 
		join milestoneskills on globalskills.id = milestoneskills.skill_id
		join milestones on milestones.id = milestoneskills.milestone_id
		where globalskills.value in (?)
		)
		and jobs.status in (?)`
)
