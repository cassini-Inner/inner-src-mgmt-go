package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *userResolver) Skills(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Skill, error) {
	skills, err := r.SkillsRepo.GetByUserId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Skill
	for _, skill := range skills {
		var gqlskill gqlmodel.Skill
		gqlskill.MapDbToGql(*skill)
		result = append(result, &gqlskill)
	}
	return result, nil
}

func (r *userResolver) CreatedJobs(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.Job, error) {
	applications, err := r.JobsRepo.GetByUserId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Job

	for _, application := range applications {
		var gqlJob gqlmodel.Job
		gqlJob.MapDbToGql(*application)
		result = append(result, &gqlJob)
	}

	return result, nil
}

func (r *userResolver) JobStats(ctx context.Context, obj *gqlmodel.User) (*gqlmodel.UserStats, error) {
	return r.JobsRepo.GetStatsByUserId(obj.ID)
}

func (r *userResolver) AppliedJobs(ctx context.Context, obj *gqlmodel.User) ([]*gqlmodel.UserJobApplication, error) {
	applications, err := r.ApplicationsRepo.GetUserJobApplications(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.UserJobApplication

	for _, job := range applications {
		var gqlJob gqlmodel.Job
		gqlJob.MapDbToGql(*job)
		jobApplicationStatus, err := r.ApplicationsService.GetApplicationStatusForUserAndJob(ctx, obj.ID, job.Id)
		if err != nil {
			return nil, err
		}

		result = append(result, &gqlmodel.UserJobApplication{Job: &gqlJob, ApplicationStatus: gqlmodel.ApplicationStatus(jobApplicationStatus), UserJobStatus: gqlmodel.JobStatus(job.Status)})
	}

	return result, nil
}
