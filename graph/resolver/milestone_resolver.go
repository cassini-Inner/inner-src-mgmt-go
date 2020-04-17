package resolver

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	//return r.JobsRepo.GetById(obj.JobID)
	panic("not implemented")
}

//TODO: Should not exist, or find a better way to implement
func (m milestonesResolver) TotalCount(ctx context.Context, obj *model.Milestones) (*int, error) {
	panic("implement me")
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	//return r.UsersRepo.GetById(obj.AssignedTo)
	panic("not implemented")
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return r.SkillsRepo.GetByMilestoneId(obj.ID)
}
