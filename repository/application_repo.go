package repository

import (
	"context"
	"errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoExistingApplications = errors.New("user does not have any existing pending or accepted applications")
)

type ApplicationsRepo interface {
	Repository
	GetExistingUserApplications(milestones []*dbmodel.Milestone, userId string, tx *sqlx.Tx, applicationStatus ...string) ([]*dbmodel.Application, error)

	CreateApplication(milestones []*dbmodel.Milestone, userId string, ctx context.Context, tx *sqlx.Tx) ([]*dbmodel.Application, error)

	SetApplicationStatusForUserMilestone(milestoneIds []string, userId string, applicationStatus string, note string, tx *sqlx.Tx) ([]*dbmodel.Application, error)

	GetByJobId(jobId string) ([]*dbmodel.Application, error)

	GetApplicationStatusForUserAndJob(userId, jobId string, tx *sqlx.Tx) (string, error)

	SetApplicationStatusForUserAndJob(userId, jobId string, milestones []*dbmodel.Milestone, applicationStatus string, note *string, tx *sqlx.Tx, ctx context.Context) ([]*dbmodel.Application, error)

	GetAcceptedApplicationsByJobId(jobId string) ([]*dbmodel.Application, error)

	GetUserJobApplications(userId string) ([]*dbmodel.Job, error)

	DeleteAllJobApplications(tx *sqlx.Tx, jobId string) error
}
