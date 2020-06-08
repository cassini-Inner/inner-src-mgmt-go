package service

import (
	"context"

	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type JobsService interface {
	CreateJob(ctx context.Context, job *gqlmodel.CreateJobInput) (result *gqlmodel.Job, err error)
	GetAllJobs(ctx context.Context, skills, status []string) ([]dbmodel.Job, error)
	GetAllJobsPaginated(ctx context.Context, skills, status []string, limit int, cursor *string) ([]*gqlmodel.Job, error)
	UpdateJobDiscussion(ctx context.Context, commentId, comment string) (*gqlmodel.Comment, error)
	DeleteJobDiscussion(ctx context.Context, commentId string) (*gqlmodel.Comment, error)
	GetById(ctx context.Context, jobId string) (*gqlmodel.Job, error)
	GetByTitle(ctx context.Context, jobTitle string, limit *int) ([]dbmodel.Job, error)
	ToggleJobCompleted(ctx context.Context, jobID string) (*gqlmodel.Job, error)
	DeleteJob(ctx context.Context, jobID string) (*gqlmodel.Job, error)
	ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*gqlmodel.Milestone, error)
	AddDiscussionToJob(ctx context.Context, comment, jobId string) (*gqlmodel.Comment, error)
	GetByMilestonesForJobIds(jobIds ...string)([]*dbmodel.Milestone, error)
}
