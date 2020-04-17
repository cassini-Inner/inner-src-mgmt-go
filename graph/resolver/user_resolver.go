package resolver

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *userResolver) Skills(ctx context.Context, obj *model.User) ([]*model.Skill, error) {
	return r.SkillsRepo.GetByUserId(obj.ID)
}

func (r *userResolver) CreatedJobs(ctx context.Context, obj *model.User) ([]*model.Job, error) {
	// return r.JobsRepo.GetByUserId(obj.ID)
	panic("Not implemented")
}

func (r *userResolver) AppliedJobs(ctx context.Context, obj *model.User) ([]*model.Job, error) {
	// return r.ApplicationsRepo.GetUserJobApplications(obj.ID)
	panic("Not implemented")
}

func (r *userResolver) JobStats(ctx context.Context, obj *model.User) (*model.UserStats, error) {
	return r.JobsRepo.GetStatsByUserId(obj.ID)
}
