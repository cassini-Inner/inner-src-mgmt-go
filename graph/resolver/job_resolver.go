package resolver

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
)

func (r *jobResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}

func (r *jobResolver) Discussion(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Discussions, error) {
	discussionsList, err := r.DiscussionsRepo.GetByJobId(obj.ID)

	var commentsList []*gqlmodel.Comment
	if err != nil {
		return nil, err
	}
	for _, discussion := range discussionsList {
		var comment gqlmodel.Comment
		comment.MapDbToGql(*discussion)
		commentsList = append(commentsList, &comment)
	}
	commentsLength := len(commentsList)
	return &gqlmodel.Discussions{Discussions: commentsList, TotalCount: &commentsLength}, nil
}

//Get the list of milestones in dbmodel type, converts it to gqlmodel type and returns list of milestones
func (r *jobResolver) Milestones(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Milestones, error) {
	return dataloader.GetMilestonesByJobIdLoader(ctx).Load(obj.ID)

}

func (r *jobResolver) Skills(ctx context.Context, obj *gqlmodel.Job) ([]*gqlmodel.Skill, error) {
	return dataloader.GetSkillByJobIdLoader(ctx).Load(obj.ID)
}

func (r *jobResolver) Applications(ctx context.Context, obj *gqlmodel.Job) (*gqlmodel.Applications, error) {
	return dataloader.GetApplicationsByJobIdLoader(ctx).Load(obj.ID)
}


func (r *jobResolver) ViewerHasApplied(ctx context.Context, obj *model.Job) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}
