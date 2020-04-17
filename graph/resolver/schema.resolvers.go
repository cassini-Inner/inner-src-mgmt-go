package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
)

// Application returns generated.ApplicationResolver implementation.
func (r *Resolver) Application() generated.ApplicationResolver { return &applicationResolver{r} }

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Milestone returns generated.MilestoneResolver implementation.
func (r *Resolver) Milestone() generated.MilestoneResolver { return &milestoneResolver{r} }

// Milestones returns generated.MilestonesResolver implementation.
func (r *Resolver) Milestones() generated.MilestonesResolver { return &milestonesResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Skill returns generated.SkillResolver implementation.
func (r *Resolver) Skill() generated.SkillResolver { return &skillResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type applicationResolver struct{ *Resolver }
type commentResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type milestoneResolver struct{ *Resolver }
type milestonesResolver struct{ *Resolver }


type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type skillResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
