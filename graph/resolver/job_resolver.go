package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *jobResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.User, error) {
	return r.UsersRepo.GetById(obj.CreatedBy)
}

func (r *jobResolver) Discussion(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Discussions, error) {
	return r.DiscussionsRepo.GetByJobId(obj.ID)
}

//Get the list of milestones in dbmodel type, converts it to gqlmodel type and returns list of milestones
func (r *jobResolver) Milestones(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Milestones, error) {
	var milestone gqlmodel.Milestone 
	var milestones gqlmodel.Milestones
	dbmilestones, err := r.MilestonesRepo.GetByJobId(obj.ID)
	for _, m := range dbmilestones {
		milestone.MapDbToGql(*m)
		milestones.Milestones = append(milestones.Milestones, &milestone)
	}
	return &milestones, err
}

func (r *jobResolver) Skills(ctx context.Context, obj *gqlmodel.Job) ([]*gqlmodel.Skill, error) {
	return r.SkillsRepo.GetByJobId(obj.ID)
}

func (r *jobResolver) Applications(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Applications, error) {
	return r.ApplicationsRepo.GetByJobId(obj.ID)
}
