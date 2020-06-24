package resolver

import (
	"context"
	"encoding/base64"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *queryResolver) AllJobs(ctx context.Context, filter *gqlmodel.JobsFilterInput) ([]*gqlmodel.Job, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

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

func (r *queryResolver) Skills(ctx context.Context, query string, limit *int) (result []*gqlmodel.Skill, err error) {
	skills, err := r.SkillsService.GetMatchingSkills(query, limit)
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}

	return result, nil
}

func (r *queryResolver) Search(ctx context.Context, query string, limit *int) (*gqlmodel.SearchResult, error) {
	//For fetching jobs with title similar to query string
	jobsFromDb, err := r.JobsService.GetByTitle(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var jobs []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		jobs = append(jobs, &tempJob)
	}

	//For fetching users with name similar to query string
	usersFromDb, err := r.UserService.GetByName(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var users []*gqlmodel.User
	for _, dbUser := range usersFromDb {
		var tempUser gqlmodel.User
		tempUser.MapDbToGql(dbUser)
		users = append(users, &tempUser)
	}

	//Search result with jobs and users
	searchResult := gqlmodel.SearchResult{
		Jobs:  jobs,
		Users: users,
	}

	return &searchResult, nil
}

func (r *queryResolver) Jobs(ctx context.Context, filter *gqlmodel.JobsFilterInput, limit int, after *string) (connection *gqlmodel.JobsConnection, err error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

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

	jobs, err := r.JobsService.GetAllJobsPaginated(ctx, skills, statuses, limit, after)
	if err != nil {
		return connection, err
	}
	var edges []*gqlmodel.JobEdge

	for i, job := range jobs {
		if i < limit {
			edges = append(edges, &gqlmodel.JobEdge{
				Node:   job,
				Cursor: base64.StdEncoding.EncodeToString([]byte(job.ID)),
			})
		}
	}
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &gqlmodel.JobsConnection{
		//TODO: Implement
		TotalCount: 10,
		Edges:      edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(jobs) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}
