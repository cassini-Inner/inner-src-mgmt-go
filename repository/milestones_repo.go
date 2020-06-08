package repository

import (
	"context"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type MilestonesRepo interface {
	Repository

	CreateMilestones(ctx context.Context, tx *sqlx.Tx, jobId string, milestones []*dbmodel.Milestone) (createdMilestones []*dbmodel.Milestone, err error)

	GetByJobId(tx sqlx.Ext, jobId string) ([]*dbmodel.Milestone, error)
	GetByJobIds(jobIds ...string) ([]*dbmodel.Milestone, error)
	GetIdsByJobId(tx sqlx.Ext, jobId string) (result []string, err error)
	GetById(milestoneId string) (*dbmodel.Milestone, error)
	GetAuthor(milestoneId string) (*dbmodel.User, error)

	ForceAutoUpdateMilestoneStatusByJobID(ctx context.Context, tx *sqlx.Tx, jobId string) error
	ForceAutoUpdateMilestoneStatusByMilestoneId(ctx context.Context, tx *sqlx.Tx, milestoneID string) error
	MarkMilestonesCompleted(tx *sqlx.Tx, ctx context.Context, milestoneIds ...string) error
	SetMilestoneAssignedTo(tx *sqlx.Tx, milestoneId string, userId *string) (*dbmodel.Milestone, error)
	DeleteMilestonesByJobId(tx *sqlx.Tx, jobID string) error
}
