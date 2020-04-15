package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-src-mgmt-go/graph/model"
)

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, user *model.UpdateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserProfile(ctx context.Context, user *model.CreateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *model.CreateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *model.UpdateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCommment(ctx context.Context, id string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *model.ApplicationStatus) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}
