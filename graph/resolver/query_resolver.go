package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *gqlmodel.JobsFilterInput) ([]*gqlmodel.Job, error) {

	if filter == nil {
		filter =&gqlmodel.JobsFilterInput{
			Status:    []*gqlmodel.JobStatus{},
			Skills:    []*string{},
			SortOrder: nil,
		}
	}

	// if the list of skills is empty, return all jobs
	var skills []*string
	if len(filter.Skills) == 0 {
		dbSkills, err := r.SkillsRepo.GetAll()
		if err != nil {
			return nil, err
		}
		for _, skill := range dbSkills {
			skillValue := skill.Value
			skills = append(skills, &skillValue)
		}

		filter.Skills = skills
	}

	var status []*gqlmodel.JobStatus
	if len(filter.Status) == 0 {
		open := gqlmodel.JobStatus("open")
		ongoing := gqlmodel.JobStatus("ongoing")
		completed := gqlmodel.JobStatus("completed")
		status = append(status, &open, &ongoing, &completed)

		filter.Status = status
	}

	jobsFromDb, err := r.JobsRepo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(*dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (*gqlmodel.Job, error) {
	var j gqlmodel.Job
	job, err := r.JobsRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	j.MapDbToGql(*job)
	return &j, err
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *gqlmodel.JobStatus) (*gqlmodel.User, error) {
	user, err := r.UsersRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
