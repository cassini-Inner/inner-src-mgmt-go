package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	user, err := r.UsersRepo.GetById(obj.ApplicantID)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *applicationResolver) Milestone(ctx context.Context, obj *model.Application) (*model.Milestone, error) {
	dbMilestone, err := r.MilestonesRepo.GetById(obj.MilestoneID)
	if err != nil {
		return nil, err
	}

	var result gqlmodel.Milestone
	result.MapDbToGql(*dbMilestone)
	return &result, nil
}

func (r *commentResolver) CreatedBy(ctx context.Context, obj *model.Comment) (*model.User, error) {
	user, err := r.UsersRepo.GetById(obj.CreatedBy)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	user, err := r.UsersRepo.GetById(obj.CreatedBy)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	skills, err := r.SkillsRepo.GetByJobId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Skill
	for _, skill := range skills {
		var gqlskill gqlmodel.Skill
		gqlskill.MapDbToGql(*skill)
		result = append(result, &gqlskill)
	}
	return result, nil
}

func (r *jobResolver) Discussion(ctx context.Context, obj *model.Job) (*model.Discussions, error) {
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

func (r *jobResolver) Milestones(ctx context.Context, obj *model.Job) (*model.Milestones, error) {
	var milestones gqlmodel.Milestones
	dbmilestones, err := r.MilestonesRepo.GetByJobId(obj.ID)
	for _, m := range dbmilestones {
		var milestone gqlmodel.Milestone
		milestone.MapDbToGql(*m)
		milestones.Milestones = append(milestones.Milestones, &milestone)
	}
	totalLength := len(milestones.Milestones)
	milestones.TotalCount = &totalLength
	return &milestones, err
}

func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	applications, err := r.ApplicationsRepo.GetByJobId(obj.ID)
	if err != nil {
		return nil, err
	}

	var gqlApplicationsList []*gqlmodel.Application
	for _, application := range applications {
		var gqlApplication gqlmodel.Application
		gqlApplication.MapDbToGql(*application)
		gqlApplicationsList = append(gqlApplicationsList, &gqlApplication)
	}

	//TODO: Implement the counters
	return &gqlmodel.Applications{Applications: gqlApplicationsList}, nil
}

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	var job gqlmodel.Job
	dbjob, err := r.JobsRepo.GetById(obj.JobID)
	if err != nil {
		return nil, err
	}
	job.MapDbToGql(*dbjob)
	return &job, nil
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	user, err := r.UsersRepo.GetById(obj.AssignedTo)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	skills, err := r.SkillsRepo.GetByMilestoneId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Skill
	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}
	return result, nil
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, user *model.UpdateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserProfile(ctx context.Context, user *model.CreateUserInput) (*model.User, error) {
	var dbuser *dbmodel.User
	var gqluser gqlmodel.User
	dbuser, err := r.UsersRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	gqluser.MapDbToGql(*dbuser)
	return &gqluser, err
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *model.CreateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *model.UpdateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCommment(ctx context.Context, id string) (*model.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *model.ApplicationStatus) (*model.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*model.UserAuthenticationPayload, error) {
	user, err := r.UsersRepo.AuthenticateAndGetUser(githubCode)
	if err != nil {
		return nil, err
	}

	var resultUser gqlmodel.User
	resultUser.MapDbToGql(*user)

	resultPayload := &gqlmodel.UserAuthenticationPayload{
		Profile: &resultUser,
		Token:   "",
	}
	return resultPayload, nil
}

func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	jobsFromDb, err := r.JobsRepo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(*dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	var j gqlmodel.Job
	job, err := r.JobsRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	j.MapDbToGql(*job)
	return &j, err
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	user, err := r.UsersRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *skillResolver) CreatedBy(ctx context.Context, obj *model.Skill) (*model.User, error) {
	user, err := r.UsersRepo.GetById(obj.CreatedBy)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}

func (r *userResolver) Onboarded(ctx context.Context, obj *model.User) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Skills(ctx context.Context, obj *model.User) ([]*model.Skill, error) {
	skills, err := r.SkillsRepo.GetByUserId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Skill
	for _, skill := range skills {
		var gqlskill gqlmodel.Skill
		gqlskill.MapDbToGql(*skill)
		result = append(result, &gqlskill)
	}
	return result, nil
}

func (r *userResolver) CreatedJobs(ctx context.Context, obj *model.User) ([]*model.Job, error) {
	applications, err := r.JobsRepo.GetByUserId(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.Job

	for _, application := range applications {
		var gqlJob gqlmodel.Job
		gqlJob.MapDbToGql(*application)
		result = append(result, &gqlJob)
	}

	return result, nil
}

func (r *userResolver) AppliedJobs(ctx context.Context, obj *model.User) ([]*model.UserJobApplication, error) {
	applications, err := r.ApplicationsRepo.GetUserJobApplications(obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*gqlmodel.UserJobApplication

	for _, application := range applications {
		var gqlJob gqlmodel.Job
		gqlJob.MapDbToGql(*application)
		result = append(result, &gqlmodel.UserJobApplication{Job: &gqlJob})
	}

	return result, nil
}

func (r *userResolver) JobStats(ctx context.Context, obj *model.User) (*model.UserStats, error) {
	return r.JobsRepo.GetStatsByUserId(obj.ID)
}

// Application returns generated.ApplicationResolver implementation.
func (r *Resolver) Application() generated.ApplicationResolver { return &applicationResolver{r} }

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Milestone returns generated.MilestoneResolver implementation.
func (r *Resolver) Milestone() generated.MilestoneResolver { return &milestoneResolver{r} }

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
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type skillResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
