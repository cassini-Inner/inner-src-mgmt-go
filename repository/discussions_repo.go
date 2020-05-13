package repository

import (
	"context"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type DiscussionsRepo interface {
	Repository
	CreateComment(ctx context.Context, tx *sqlx.Tx, jobId, comment, userId string) (*dbmodel.Discussion, error)
	UpdateComment(ctx context.Context, tx *sqlx.Tx, discussionId, content string) (*dbmodel.Discussion, error)
	DeleteComment(ctx context.Context, tx *sqlx.Tx, discussionId string) (*dbmodel.Discussion, error)
	GetByJobId(jobId string) ([]*dbmodel.Discussion, error)
	GetById(tx *sqlx.Tx, discussionId string) (*dbmodel.Discussion, error)
	DeleteAllCommentsForJob(tx *sqlx.Tx, jobID string) error
}
