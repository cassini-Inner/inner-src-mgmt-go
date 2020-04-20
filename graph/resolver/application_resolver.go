package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.User, error) {
	user, err := r.UsersRepo.GetById(obj.ApplicantID)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *applicationResolver) Milestone(ctx context.Context, obj *gqlmodel.Application) (*gqlmodel.Milestone, error) {
	dbMilestone, err := r.MilestonesRepo.GetById(obj.MilestoneID)
	if err != nil {
		return nil, err
	}

	var result gqlmodel.Milestone
	result.MapDbToGql(*dbMilestone)
	return &result, nil
}
