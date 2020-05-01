package repository

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type JobsRepo interface {
	CreateJob(ctx context.Context, tx *sqlx.Tx, input *gqlmodel.CreateJobInput, user *dbmodel.User) (*dbmodel.Job, error)

	UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.Job, error)

	DeleteJob(tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)

	GetById(db sqlx.Ext, jobId string) (*dbmodel.Job, error)

	GetByUserId(userId string) ([]*dbmodel.Job, error)

	GetStatsByUserId(userId string) (*gqlmodel.UserStats, error)

	GetAll(tx sqlx.Ext, skillNames []string, status []string) ([]dbmodel.Job, error)

	GetMilestonesByJobId(tx sqlx.Ext, jobId string) ([]*dbmodel.Milestone, error)

	GetMilestoneIdsByJobId(tx sqlx.Ext, jobId string) (result []string, err error)

	GetMilestoneById(milestoneId string) (*dbmodel.Milestone, error)

	GetAuthorFromMilestoneId(milestoneId string) (*dbmodel.User, error)

	MarkJobCompleted(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)

	ForceAutoUpdateJobStatus(ctx context.Context, tx *sqlx.Tx, jobId string) (*dbmodel.Job, error)

	ForceAutoUpdateMilestoneStatusByJobID(ctx context.Context, tx *sqlx.Tx, jobId string) error

	ForceAutoUpdateMilestoneStatusByMilestoneId(ctx context.Context, tx *sqlx.Tx, milestoneID string) error

	MarkMilestonesCompleted(tx *sqlx.Tx, ctx context.Context, milestoneIds ...string) error

	CreateMilestones(ctx context.Context, tx *sqlx.Tx, jobId string, milestones []*gqlmodel.MilestoneInput) (createdMilestones []*dbmodel.Milestone, err error)

	DeleteMilestonesByJobId(tx *sqlx.Tx, jobID string) error
}
