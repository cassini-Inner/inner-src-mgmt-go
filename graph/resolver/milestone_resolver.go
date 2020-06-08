package resolver

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *gqlmodel.Milestone) (*gqlmodel.Job, error) {
	return r.JobsService.GetById(ctx, obj.JobID)
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

func (r *milestoneResolver) Review(ctx context.Context, obj *gqlmodel.Milestone) (*gqlmodel.Review, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetJobMilestoneReviewLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, obj.AssignedTo))
}
