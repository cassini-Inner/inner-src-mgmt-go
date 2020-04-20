package resolver

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *gqlmodel.JobsFilterInput) ([]*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Job(ctx context.Context, id string) (*gqlmodel.Job, error) {
	var j gqlmodel.Job
	job, err := r.JobsRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	j.MapDbToGql(*job)
	return &j, err
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *gqlmodel.JobStatus) (*gqlmodel.User, error) {
	user, err := r.UsersRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
