package resolver

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *gqlmodel.JobsFilterInput) ([]*gqlmodel.Job, error) {
	var skills []string
	var statuses []string

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobsFromDb, err := r.JobsService.GetAllJobs(ctx, skills, statuses)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (*gqlmodel.Job, error) {
	return r.JobsService.GetById(ctx, id)
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *gqlmodel.JobStatus) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(id)
}

func (r *queryResolver) Skills(ctx context.Context, query *string) ([]*gqlmodel.Skill, error) {
	panic(fmt.Errorf("not implemented"))
}
