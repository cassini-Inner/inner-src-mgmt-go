package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *gqlmodel.Milestone) (*gqlmodel.Job, error) {
	var job gqlmodel.Job
	dbjob, err := r.JobsRepo.GetById(obj.JobID)
	if err != nil {
		return nil, err
	}
	job.MapDbToGql(*dbjob)
	return &job, nil
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *gqlmodel.Milestone) (*gqlmodel.User, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.AssignedTo)
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *gqlmodel.Milestone) ([]*gqlmodel.Skill, error) {
	return dataloader.GetSkillByMilestoneIdLoader(ctx).Load(obj.ID)
}
