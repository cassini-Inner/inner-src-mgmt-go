package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jinzhu/gorm"
	"fmt"
)

type JobsRepo struct {
	db *gorm.DB
}

func NewJobsRepo(db *gorm.DB) *JobsRepo {
	return &JobsRepo{db: db}
}

func (j *JobsRepo) CreateJob(input *model.CreateJobInput) (*model.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) UpdateJob(input *model.UpdateJobInput) (*model.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) DeleteJob(jobId string) (*model.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) GetById(jobId string) (*model.Job, error) {
	panic("Not implemented")
}

// GetByUserId returns all jobs created by that user
func (j *JobsRepo) GetByUserId(userId string) ([]*model.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) GetStatsByUserId(userId string) (*model.UserStats, error) {
	panic("not implemented")
}

func (j *JobsRepo) GetAll(filters *model.JobsFilterInput) ([]*model.Job, error) {
	panic("not implemented")
}
