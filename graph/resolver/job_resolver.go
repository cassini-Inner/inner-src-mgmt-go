package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-source-mgmt-srv/graph/model"
)

func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *jobResolver) Discussion(ctx context.Context, obj *model.Job) (*model.Discussions, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *jobResolver) Milestones(ctx context.Context, obj *model.Job) (*model.Milestones, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	panic(fmt.Errorf("not implemented"))
}
