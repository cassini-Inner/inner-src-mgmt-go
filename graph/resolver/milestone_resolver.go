package resolver

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	var j model.Job 
	dbjob, err := r.JobsRepo.GetById(obj.JobID)
	j.MapDbToGql(*dbjob)
	return &j, err
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	return r.UsersRepo.GetById(obj.AssignedTo)
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return r.SkillsRepo.GetByMilestoneId(obj.ID)
}

func (m milestonesResolver) TotalCount(ctx context.Context, obj *model.Milestones) (*int, error) {
	totalCount := len(obj.Milestones)
	return &totalCount, nil
 }