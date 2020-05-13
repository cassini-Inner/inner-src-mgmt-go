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
	GetExistingUserApplications(tx *sqlx.Tx, milestones []*dbmodel.Milestone, userId string, applicationStatus ...string) ([]*dbmodel.Application, error)
	CreateApplication(ctx context.Context, tx *sqlx.Tx, milestones []*dbmodel.Milestone, userId string) ([]*dbmodel.Application, error)
	SetApplicationStatusForUserMilestone(tx *sqlx.Tx, milestoneIds []string, userId string, applicationStatus string, note string) ([]*dbmodel.Application, error)
	GetByJobId(jobId string) ([]*dbmodel.Application, error)
	GetApplicationStatusForUserAndJob(userId string, tx *sqlx.Tx, jobId string) (string, error)
	SetApplicationStatusForUserAndJob(ctx context.Context, tx *sqlx.Tx, milestones []*dbmodel.Milestone, applicationStatus string, note *string, jobId, userId string) ([]*dbmodel.Application, error)
	GetAcceptedApplicationsByJobId(jobId string) ([]*dbmodel.Application, error)
	GetUserJobApplications(userId string) ([]*dbmodel.Job, error)
	DeleteAllJobApplications(tx *sqlx.Tx, jobId string) error
}
