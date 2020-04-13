package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-source-mgmt-srv/graph/model"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Job(ctx context.Context, id int) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, id int, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
