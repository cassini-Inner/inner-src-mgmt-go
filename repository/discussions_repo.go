package repository

import (
	"context"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type DiscussionsRepo interface {
	Repository

	CreateComment(jobId, comment, userId string, tx *sqlx.Tx, ctx context.Context) (*dbmodel.Discussion, error)
	UpdateComment(discussionId, content string, tx *sqlx.Tx, ctx context.Context) (*dbmodel.Discussion, error)
	DeleteComment(discussionId string, tx *sqlx.Tx, ctx context.Context) (*dbmodel.Discussion, error)
	GetByJobId(jobId string) ([]*dbmodel.Discussion, error)
	GetById(discussionId string, tx *sqlx.Tx) (*dbmodel.Discussion, error)
	DeleteAllCommentsForJob(tx *sqlx.Tx, jobID string) error
}
