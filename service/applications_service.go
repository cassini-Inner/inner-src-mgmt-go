package service

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

type ApplicationsService interface {
	CreateUserJobApplication(ctx context.Context, jobId string) ([]*gqlmodel.Application, error)

	GetApplicationStatusForUserAndJob(ctx context.Context, userId string, joinId string) (string, error)

	UpdateJobApplicationStatus(ctx context.Context, applicantId string, jobId string, status *gqlmodel.ApplicationStatus, note *string) ([]*gqlmodel.Application, error)

	DeleteUserJobApplication(ctx context.Context, jobId string) ([]*gqlmodel.Application, error)
}
