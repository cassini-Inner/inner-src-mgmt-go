package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *userResolver) Skills(ctx context.Context, obj *model.User) ([]*model.Skill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) CreatedJobs(ctx context.Context, obj *model.User) ([]*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) AppliedJobs(ctx context.Context, obj *model.User) ([]*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) JobStats(ctx context.Context, obj *model.User) (*model.UserStats, error) {
	panic(fmt.Errorf("not implemented"))
}
