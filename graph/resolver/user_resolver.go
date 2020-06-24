package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
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

func (r *userResolver) Reviews(ctx context.Context, obj *gqlmodel.User) (result []*gqlmodel.JobReview, err error) {
	// get applied jobs
	appliedJobs, err := r.ApplicationsService.GetAcceptedAppliedJobs(ctx, obj.ID)
	var mappedJobs []*gqlmodel.Job
	for _, job := range appliedJobs {
		tempJob := &gqlmodel.Job{}
		tempJob.MapDbToGql(*job)
		mappedJobs = append(mappedJobs, tempJob)
	}

	// map milestones to job
	jobMilestones := make(map[string][]*gqlmodel.Milestone)
	if err != nil {
		return nil, err
	}
	var jobIds []string
	for _, job := range appliedJobs {
		jobIds = append(jobIds, job.Id)
	}

	// if the user has not applied to any jobs
	if len(jobIds) == 0 {
		return result, nil
	}
	milestones, err := r.JobsService.GetByMilestonesForJobIds(jobIds...)
	if err != nil {
		return nil, err
	}

	for _, milestone := range milestones {
		if _, ok := jobMilestones[milestone.JobId]; !ok {
			jobMilestones[milestone.JobId] = make([]*gqlmodel.Milestone, 0)
		}
		mappedMilestone := gqlmodel.Milestone{}
		mappedMilestone.MapDbToGql(*milestone)
		jobMilestones[milestone.JobId] = append(jobMilestones[milestone.JobId], &mappedMilestone)
	}

	// a map of milestoneId -> review
	reviewMilestoneMap := make(map[string]*gqlmodel.Review)
	reviews, err := r.ReviewsService.GetForUserId(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	for _, review := range reviews {
		gqlReview := &gqlmodel.Review{}
		gqlReview.MapDbToGql(*review)
		reviewMilestoneMap[review.MilestoneId] = gqlReview
	}

	for _, job := range mappedJobs {
		var milestoneReviews []*gqlmodel.MilestoneReview
		for _, milestone := range jobMilestones[job.ID] {
			milestoneReviews = append(milestoneReviews, &gqlmodel.MilestoneReview{
				Review:    reviewMilestoneMap[milestone.ID],
				Milestone: milestone,
			})
		}
		result = append(result, &gqlmodel.JobReview{
			Job:             job,
			MilestoneReview: milestoneReviews,
		})
	}

	return result, nil
}

func (r *userResolver) OverallRating(ctx context.Context, obj *gqlmodel.User) (*int, error) {
	return dataloader.GetUserAverageRatingLoader(ctx).Load(obj.ID)
}
