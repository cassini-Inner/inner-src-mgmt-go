package postgres

import "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"

type JobsRepo struct{}

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
