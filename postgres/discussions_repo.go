package postgres

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/go-pg/pg/v9"
)

type DiscussionsRepo struct {
	db *pg.DB
}

func NewDiscussionsRepo(db *pg.DB) *DiscussionsRepo {
return &DiscussionsRepo{    db:db}
}

//TODO: Implement
func (d *DiscussionsRepo) CreateComment(jobId string, comment string) (*model.Comment, error) {
	panic("not implemented")
}
func (d *DiscussionsRepo) UpdateComment(commentId string, comment string) (*model.Comment, error) {
	panic("not implemented")
}
func (d *DiscussionsRepo) DeleteComment(commentId string) (*model.Comment, error) {
	panic("not implemented")
}

func (d *DiscussionsRepo) GetByJobId(jobId string) (*model.Discussions, error) {
	panic("not implemented")
}
