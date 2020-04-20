package resolver

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
)

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, user *gqlmodel.UpdateUserInput) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserProfile(ctx context.Context, user *gqlmodel.CreateUserInput) (*gqlmodel.User, error) {
	var dbuser *dbmodel.User
	var gqluser gqlmodel.User
	dbuser, err := r.UsersRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	gqluser.MapDbToGql(*dbuser)
	return &gqluser, err
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *gqlmodel.CreateJobInput) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *gqlmodel.UpdateJobInput) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCommment(ctx context.Context, id string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *gqlmodel.ApplicationStatus) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}
