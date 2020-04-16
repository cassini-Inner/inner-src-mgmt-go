package resolver

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	return r.UsersRepo.GetById(obj.CreatedBy)
}

func (r *jobResolver) Discussion(ctx context.Context, obj *model.Job) (*model.Discussions, error) {
	return r.DiscussionsRepo.GetByJobId(obj.ID)
}

func (r *jobResolver) Milestones(ctx context.Context, obj *model.Job) (*model.Milestones, error) {
	return r.MilestonesRepo.GetByJobId(obj.ID)
}

func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	return r.SkillsRepo.GetByJobId(obj.ID)
}

func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	return r.ApplicationsRepo.GetByJobId(obj.ID)
}
