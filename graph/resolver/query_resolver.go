package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	return r.jobsRepo.GetById(id)
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	return r.usersRepo.GetById(id)
}
