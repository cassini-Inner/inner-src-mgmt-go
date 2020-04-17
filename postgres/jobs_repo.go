package postgres

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/models"
	"github.com/jinzhu/gorm"
)

type JobsRepo struct {
	db *gorm.DB
}

func NewJobsRepo(db *gorm.DB) *JobsRepo {
	return &JobsRepo{db: db}
}

func (j *JobsRepo) CreateJob(input *gqlmodel.CreateJobInput) (*dbmodel.job, error) {
	// query := "INSERT INTO jobs (CreatedBy, Title, Description, Difficulty, Status, TimeCreated, TimeUpdated, IsDeleted) VALUES ("
	// err := j.db.QueryRowx()
}

func (j *JobsRepo) UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) DeleteJob(jobId string) (*dbmodel.job, error) {
	panic("Not implemented")
}

// Get the complete job details based on the job id
func (j *JobsRepo) GetById(jobId string) (*dbmodel.job, error) {
	var job dbmodel.job
	query := "SELECT * FROM jobs WHERE id = $1"
	err := j.db.QueryRowx(query, JobId).StructScan(&job)
	return job, err
}

// GetByUserId returns all jobs created by that user
func (j *JobsRepo) GetByUserId(userId string) ([]*dbmodel.job, error) {
	var job dbmodel.job
	query := "SELECT * FROM jobs WHERE created_by = $1" 
	
}

func (j *JobsRepo) GetStatsByUserId(userId string) (*dbmodel.UserStats, error) {
	panic("not implemented")
}

func (j *JobsRepo) GetAll(filters *gqlmodel.JobsFilterInput) ([]*dbmodel.job, error) {
	panic("not implemented")
}
