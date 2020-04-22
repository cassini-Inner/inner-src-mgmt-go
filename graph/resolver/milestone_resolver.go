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
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return getUserLoader(ctx).Load(obj.AssignedTo)
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
	return result, nil
}
