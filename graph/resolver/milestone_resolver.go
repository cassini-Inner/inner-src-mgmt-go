package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
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
	user, err := r.UsersRepo.GetById(obj.AssignedTo)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *gqlmodel.Milestone) ([]*gqlmodel.Skill, error) {
	skills, err := r.SkillsRepo.GetByMilestoneId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Skill
	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}
<<<<<<< HEAD
	return result, err
=======
	return result, nil
>>>>>>> 752a2c7c9ec346312a51fe3a084dcdf2d8c98bb2
}
