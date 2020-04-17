package resolver

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	var job model.Job
	dbjob, err := r.JobsRepo.GetById(obj.JobID)
	if err != nil {
		return nil, err
	}
	job.MapDbToGql(*dbjob)
	return &job, nil
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	return r.UsersRepo.GetById(obj.AssignedTo)
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return r.SkillsRepo.GetByMilestoneId(obj.ID)
}

//TODO: Should not exist, or find a better way to implement
func (m milestonesResolver) TotalCount(ctx context.Context, obj *model.Milestones) (*int, error) {
	totalCount := len(obj.Milestones)
	return &totalCount, nil
}
