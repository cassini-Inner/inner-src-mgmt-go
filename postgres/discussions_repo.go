package postgres

import "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"

type DiscussionsRepo struct{}

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
