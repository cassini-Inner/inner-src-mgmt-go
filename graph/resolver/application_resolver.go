package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.ApplicantID)
}

func (r *applicationResolver) Milestone(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.Milestone, error) {
	dbMilestone, err := r.JobsRepo.GetMilestoneById(obj.MilestoneID)
	if err != nil {
		return nil, err
	}

	var result gqlmodel.Milestone
	result.MapDbToGql(*dbMilestone)
	return &result, nil
}
