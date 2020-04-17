package postgres

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
)

type DiscussionsRepo struct {
	db *sqlx.DB
}

func NewDiscussionsRepo(db *sqlx.DB) *DiscussionsRepo {
	return &DiscussionsRepo{db: db}
}

//TODO: Implement
func (d *DiscussionsRepo) CreateComment(jobId string, comment string) (*gqlmodel.Comment, error) {
	panic("not implemented")
}
func (d *DiscussionsRepo) UpdateComment(commentId string, comment string) (*gqlmodel.Comment, error) {
	panic("not implemented")
}
func (d *DiscussionsRepo) DeleteComment(commentId string) (*gqlmodel.Comment, error) {
	panic("not implemented")
}

func (d *DiscussionsRepo) GetByJobId(jobId string) (*gqlmodel.Discussions, error) {
	panic("not implemented")
}
