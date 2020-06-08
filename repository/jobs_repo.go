package repository

import (
	"context"

	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type JobsRepo interface {
	Repository
	CreateJob(ctx context.Context, tx *sqlx.Tx, input *gqlmodel.CreateJobInput, user *dbmodel.User) (*dbmodel.Job, error)

	GetAll(skillNames []string, status []string) ([]dbmodel.Job, error)
	GetAllPaginated(skillNames []string, status []string, limit int, cursor *string) ([]dbmodel.Job, error)
	GetById(jobId string) (*dbmodel.Job, error)
	GetByUserId(userId string) ([]*dbmodel.Job, error)
	GetByTitle(jobTitle string, limit *int) ([]dbmodel.Job, error)
	GetStatsByUserId(userId string) (*gqlmodel.UserStats, error)

	UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.Job, error)
	MarkJobCompleted(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)
	ForceAutoUpdateJobStatus(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)

	DeleteJob(tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)
}
