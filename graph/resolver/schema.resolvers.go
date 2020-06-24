package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.ApplicantID)
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
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}

func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}

func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	return dataloader.GetSkillByJobIdLoader(ctx).Load(obj.ID)
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
	return dataloader.GetMilestonesByJobIdLoader(ctx).Load(obj.ID)
}

func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	return dataloader.GetApplicationsByJobIdLoader(ctx).Load(obj.ID)
}

func (r *jobResolver) ViewerHasApplied(ctx context.Context, obj *model.Job) (bool, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return false, err
	}
	return dataloader.GetViewerHasAppliedLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, user.Id))
}

func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	return r.JobsService.GetById(ctx, obj.JobID)
}

func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.AssignedTo)
}

func (r *milestoneResolver) Review(ctx context.Context, obj *model.Milestone) (*model.Review, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetJobMilestoneReviewLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, obj.AssignedTo))
}

func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return dataloader.GetSkillByMilestoneIdLoader(ctx).Load(obj.ID)
}

func (r *mutationResolver) UpdateProfile(ctx context.Context, user *model.UpdateUserInput) (*model.User, error) {
	return r.UserService.UpdateProfile(ctx, updatedUserDetails)
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *model.CreateJobInput) (*model.Job, error) {
	createdJob, err := r.JobsService.CreateJobs(ctx, job)
	if err != nil {
		return nil, err
	}
	return createdJob[0], nil
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *model.UpdateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.DeleteJob(ctx, jobID)
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*model.Comment, error) {
	return r.JobsService.AddDiscussionToJob(ctx, comment, jobID)
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	return r.JobsService.UpdateJobDiscussion(ctx, id, comment)
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*model.Comment, error) {
	return r.JobsService.DeleteJobDiscussion(ctx, id)
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.CreateUserJobApplication(ctx, jobID)
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.DeleteUserJobApplication(ctx, jobID)
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *model.ApplicationStatus, note *string) ([]*model.Application, error) {
	return r.ApplicationsService.UpdateJobApplicationStatus(ctx, applicantID, jobID, status, note)
}

func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*model.UserAuthenticationPayload, error) {
	// authenticate the user with github and store them in db
	resultUser, err := r.AuthenticationService.AuthenticateAndGetUser(ctx, githubCode)
	if err != nil {
		return nil, err
	}
	//generate a token for the user and return
	authToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	refreshToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	resultPayload := &gqlmodel.UserAuthenticationPayload{
		Profile:      resultUser,
		Token:        *authToken,
		RefreshToken: *refreshToken,
	}
	return resultPayload, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*model.UserAuthenticationPayload, error) {
	// get the claims for the user
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("error while refreshing refreshToken %v", refreshToken)
		log.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid refreshToken signature")
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("refreshToken is not valid")
	}
	// only refresh the refreshToken if it's expiring in 2 minutes
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > (time.Minute * 2) {
		return nil, errors.New("refreshToken can only be refreshed 2 minutes from expiry time")
	}
	//generate a new refreshToken for the user
	gqlUser, err := r.UserService.GetById(ctx, claims.Id)
	if err != nil {
		log.Printf("error getting user from claims for user id %v", claims.Id)
		return nil, err
	}
	newToken, err := gqlUser.GenerateAccessToken()
	newRefreshToken, err := gqlUser.GenerateRefreshToken()
	if err != nil {
		log.Printf("error generating refreshToken for user %+v", gqlUser)
		return nil, err
	}
	return &gqlmodel.UserAuthenticationPayload{
		Profile:      gqlUser,
		Token:        *newToken,
		RefreshToken: *newRefreshToken,
	}, nil
}

func (r *mutationResolver) ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*model.Milestone, error) {
	return r.JobsService.ToggleMilestoneCompleted(ctx, milestoneID)
}

func (r *mutationResolver) ToggleJobCompleted(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.ToggleJobCompleted(ctx, jobID)
}

func (r *mutationResolver) CreateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, milestoneID string) (*model.Review, error) {
	createdReview, err := r.ReviewsService.ReviewAssignedUser(ctx, review.Rating, review.Remark, milestoneID)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*createdReview)
	return gqlReview, nil
}

func (r *mutationResolver) UpdateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, id string) (*model.Review, error) {
	updatedReview, err := r.ReviewsService.UpdateReview(ctx, review.Rating, review.Remark, id)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*updatedReview)
	return gqlReview, nil
}

func (r *mutationResolver) RestoreJobsBackup(ctx context.Context, jobs []*model.CreateJobInput) ([]*model.Job, error) {
	return r.JobsService.CreateJobs(ctx, jobs...)
}

func (r *mutationResolver) MarkAllViewerNotificationsRead(ctx context.Context) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) MarkViewerNotificationsRead(ctx context.Context, ids []string) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *notificationItemResolver) Recipient(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Recipient.ID)
}

func (r *notificationItemResolver) Sender(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Sender.ID)
}

func (r *notificationItemResolver) Job(ctx context.Context, obj *model.NotificationItem) (*model.Job, error) {
	return dataloader.GetJobByJobIdLoader(ctx).Load(obj.Job.ID)
}

func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobsFromDb, err := r.JobsService.GetAllJobs(ctx, skills, statuses)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	return r.JobsService.GetById(ctx, id)
}

func (r *queryResolver) Jobs(ctx context.Context, filter *model.JobsFilterInput, limit int, after *string) (*model.JobsConnection, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobs, err := r.JobsService.GetAllJobsPaginated(ctx, skills, statuses, limit, after)
	if err != nil {
		return connection, err
	}
	var edges []*gqlmodel.JobEdge

	for i, job := range jobs {
		if i < limit {
			edges = append(edges, &gqlmodel.JobEdge{
				Node:   job,
				Cursor: base64.StdEncoding.EncodeToString([]byte(job.ID)),
			})
		}
	}
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &gqlmodel.JobsConnection{
		//TODO: Implement
		TotalCount: 10,
		Edges:      edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(jobs) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}

func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(id)
}

func (r *queryResolver) Skills(ctx context.Context, query string, limit *int) ([]*model.Skill, error) {
	skills, err := r.SkillsService.GetMatchingSkills(query, limit)
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}

	return result, nil
}

func (r *queryResolver) Search(ctx context.Context, query string, limit *int) (*model.SearchResult, error) {
	//For fetching jobs with title similar to query string
	jobsFromDb, err := r.JobsService.GetByTitle(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var jobs []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		jobs = append(jobs, &tempJob)
	}

	//For fetching users with name similar to query string
	usersFromDb, err := r.UserService.GetByName(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var users []*gqlmodel.User
	for _, dbUser := range usersFromDb {
		var tempUser gqlmodel.User
		tempUser.MapDbToGql(dbUser)
		users = append(users, &tempUser)
	}

	//Search result with jobs and users
	searchResult := gqlmodel.SearchResult{
		Jobs:  jobs,
		Users: users,
	}

	return &searchResult, nil
}

func (r *queryResolver) ViewerNotifications(ctx context.Context, limit int, after *string) (*model.NotificationConnection, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	var afterString *string
	if after != nil {
		if *after == "" {
			return nil, custom_errors.ErrInvalidCursor
		}
		decodedAfter, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, custom_errors.ErrInvalidCursor
		}
		cursorStr := string(decodedAfter)
		afterString = &cursorStr
	}

	notifications, err := r.NotificationsService.GetAllPaginated(user.Id, afterString, limit+1)
	if err != nil {
		return nil, err
	}

	totalNotificationsForUser, err := r.NotificationsService.GetNotificationCountForReceiver(user.Id)

	if err != nil {
		return nil, err
	}

	notificationEdges := make([]*gqlmodel.NotificationEdge, 0)
	for i, notification := range notifications {
		if i < limit {
			notificationEdges = append(notificationEdges, &gqlmodel.NotificationEdge{
				Node:   notification,
				Cursor: base64.StdEncoding.EncodeToString([]byte(notification.ID)),
			})
		}
	}

	var endCursor *string
	if len(notificationEdges) > 0 {
		endCursor = &notificationEdges[len(notificationEdges)-1].Cursor
	}

	return &gqlmodel.NotificationConnection{
		TotalCount: totalNotificationsForUser,
		Edges:      notificationEdges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(notificationEdges) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}

func (r *reviewResolver) CreatedFor(ctx context.Context, obj *model.Review) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedFor)
}

func (r *skillResolver) CreatedBy(ctx context.Context, obj *model.Skill) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
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

	for _, job := range applications {
		var gqlJob gqlmodel.Job
		gqlJob.MapDbToGql(*job)
		jobApplicationStatus, err := r.ApplicationsService.GetApplicationStatusForUserAndJob(ctx, obj.ID, job.Id)
		if err != nil {
			return nil, err
		}

		result = append(result, &gqlmodel.UserJobApplication{Job: &gqlJob, ApplicationStatus: gqlmodel.ApplicationStatus(jobApplicationStatus), UserJobStatus: gqlmodel.JobStatus(job.Status)})
	}

	return result, nil
}

func (r *userResolver) JobStats(ctx context.Context, obj *model.User) (*model.UserStats, error) {
	return r.JobsRepo.GetStatsByUserId(obj.ID)
}

func (r *userResolver) Reviews(ctx context.Context, obj *model.User) ([]*model.JobReview, error) {
	// get applied jobs
	appliedJobs, err := r.ApplicationsService.GetAppliedJobs(ctx, obj.ID)
	var mappedJobs []*gqlmodel.Job
	for _, job := range appliedJobs {
		tempJob := &gqlmodel.Job{}
		tempJob.MapDbToGql(*job)
		mappedJobs = append(mappedJobs, tempJob)
	}

	// map milestones to job
	jobMilestones := make(map[string][]*gqlmodel.Milestone)
	if err != nil {
		return nil, err
	}
	var jobIds []string
	for _, job := range appliedJobs {
		jobIds = append(jobIds, job.Id)
	}

	// if the user has not applied to any jobs
	if len(jobIds) == 0 {
		return result, nil
	}
	milestones, err := r.JobsService.GetByMilestonesForJobIds(jobIds...)
	if err != nil {
		return nil, err
	}

	for _, milestone := range milestones {
		if _, ok := jobMilestones[milestone.JobId]; !ok {
			jobMilestones[milestone.JobId] = make([]*gqlmodel.Milestone, 0)
		}
		mappedMilestone := gqlmodel.Milestone{}
		mappedMilestone.MapDbToGql(*milestone)
		jobMilestones[milestone.JobId] = append(jobMilestones[milestone.JobId], &mappedMilestone)
	}

	// a map of milestoneId -> review
	reviewMilestoneMap := make(map[string]*gqlmodel.Review)
	reviews, err := r.ReviewsService.GetForUserId(ctx, obj.ID)
	if err != nil {
		return nil, err
	}
	for _, review := range reviews {
		gqlReview := &gqlmodel.Review{}
		gqlReview.MapDbToGql(*review)
		reviewMilestoneMap[review.MilestoneId] = gqlReview
	}

	for _, job := range mappedJobs {
		var milestoneReviews []*gqlmodel.MilestoneReview
		for _, milestone := range jobMilestones[job.ID] {
			milestoneReviews = append(milestoneReviews, &gqlmodel.MilestoneReview{
				Review:    reviewMilestoneMap[milestone.ID],
				Milestone: milestone,
			})
		}
		result = append(result, &gqlmodel.JobReview{
			Job:             job,
			MilestoneReview: milestoneReviews,
		})
	}

	return result, nil
}

func (r *userResolver) OverallRating(ctx context.Context, obj *model.User) (*int, error) {
	return dataloader.GetUserAverageRatingLoader(ctx).Load(obj.ID)
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

// NotificationItem returns generated.NotificationItemResolver implementation.
func (r *Resolver) NotificationItem() generated.NotificationItemResolver {
	return &notificationItemResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Review returns generated.ReviewResolver implementation.
func (r *Resolver) Review() generated.ReviewResolver { return &reviewResolver{r} }

// Skill returns generated.SkillResolver implementation.
func (r *Resolver) Skill() generated.SkillResolver { return &skillResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type applicationResolver struct{ *Resolver }
type commentResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type milestoneResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type notificationItemResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type reviewResolver struct{ *Resolver }
type skillResolver struct{ *Resolver }
type userResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.ApplicantID)
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
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}
func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}
func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	return dataloader.GetSkillByJobIdLoader(ctx).Load(obj.ID)
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
	return dataloader.GetMilestonesByJobIdLoader(ctx).Load(obj.ID)
}
func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	return dataloader.GetApplicationsByJobIdLoader(ctx).Load(obj.ID)
}
func (r *jobResolver) ViewerHasApplied(ctx context.Context, obj *model.Job) (bool, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return false, err
	}
	return dataloader.GetViewerHasAppliedLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, user.Id))
}
func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	return r.JobsService.GetById(ctx, obj.JobID)
}
func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.AssignedTo)
}
func (r *milestoneResolver) Review(ctx context.Context, obj *model.Milestone) (*model.Review, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetJobMilestoneReviewLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, obj.AssignedTo))
}
func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return dataloader.GetSkillByMilestoneIdLoader(ctx).Load(obj.ID)
}
func (r *mutationResolver) UpdateProfile(ctx context.Context, user *model.UpdateUserInput) (*model.User, error) {
	return r.UserService.UpdateProfile(ctx, updatedUserDetails)
}
func (r *mutationResolver) CreateJob(ctx context.Context, job *model.CreateJobInput) (*model.Job, error) {
	createdJob, err := r.JobsService.CreateJobs(ctx, job)
	if err != nil {
		return nil, err
	}
	return createdJob[0], nil
}
func (r *mutationResolver) UpdateJob(ctx context.Context, job *model.UpdateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.DeleteJob(ctx, jobID)
}
func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*model.Comment, error) {
	return r.JobsService.AddDiscussionToJob(ctx, comment, jobID)
}
func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	return r.JobsService.UpdateJobDiscussion(ctx, id, comment)
}
func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*model.Comment, error) {
	return r.JobsService.DeleteJobDiscussion(ctx, id)
}
func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.CreateUserJobApplication(ctx, jobID)
}
func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.DeleteUserJobApplication(ctx, jobID)
}
func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *model.ApplicationStatus, note *string) ([]*model.Application, error) {
	return r.ApplicationsService.UpdateJobApplicationStatus(ctx, applicantID, jobID, status, note)
}
func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*model.UserAuthenticationPayload, error) {
	// authenticate the user with github and store them in db
	resultUser, err := r.AuthenticationService.AuthenticateAndGetUser(ctx, githubCode)
	if err != nil {
		return nil, err
	}
	//generate a token for the user and return
	authToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	refreshToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	resultPayload := &gqlmodel.UserAuthenticationPayload{
		Profile:      resultUser,
		Token:        *authToken,
		RefreshToken: *refreshToken,
	}
	return resultPayload, nil
}
func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*model.UserAuthenticationPayload, error) {
	// get the claims for the user
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("error while refreshing refreshToken %v", refreshToken)
		log.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid refreshToken signature")
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("refreshToken is not valid")
	}
	// only refresh the refreshToken if it's expiring in 2 minutes
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > (time.Minute * 2) {
		return nil, errors.New("refreshToken can only be refreshed 2 minutes from expiry time")
	}
	//generate a new refreshToken for the user
	gqlUser, err := r.UserService.GetById(ctx, claims.Id)
	if err != nil {
		log.Printf("error getting user from claims for user id %v", claims.Id)
		return nil, err
	}
	newToken, err := gqlUser.GenerateAccessToken()
	newRefreshToken, err := gqlUser.GenerateRefreshToken()
	if err != nil {
		log.Printf("error generating refreshToken for user %+v", gqlUser)
		return nil, err
	}
	return &gqlmodel.UserAuthenticationPayload{
		Profile:      gqlUser,
		Token:        *newToken,
		RefreshToken: *newRefreshToken,
	}, nil
}
func (r *mutationResolver) ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*model.Milestone, error) {
	return r.JobsService.ToggleMilestoneCompleted(ctx, milestoneID)
}
func (r *mutationResolver) ToggleJobCompleted(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.ToggleJobCompleted(ctx, jobID)
}
func (r *mutationResolver) CreateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, milestoneID string) (*model.Review, error) {
	createdReview, err := r.ReviewsService.ReviewAssignedUser(ctx, review.Rating, review.Remark, milestoneID)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*createdReview)
	return gqlReview, nil
}
func (r *mutationResolver) UpdateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, id string) (*model.Review, error) {
	updatedReview, err := r.ReviewsService.UpdateReview(ctx, review.Rating, review.Remark, id)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*updatedReview)
	return gqlReview, nil
}
func (r *mutationResolver) RestoreJobsBackup(ctx context.Context, jobs []*model.CreateJobInput) ([]*model.Job, error) {
	return r.JobsService.CreateJobs(ctx, jobs...)
}
func (r *mutationResolver) MarkAllViewerNotificationsRead(ctx context.Context) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) MarkViewerNotificationsRead(ctx context.Context, ids []string) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *notificationItemResolver) Recipient(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Recipient.ID)
}
func (r *notificationItemResolver) Sender(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Sender.ID)
}
func (r *notificationItemResolver) Job(ctx context.Context, obj *model.NotificationItem) (*model.Job, error) {
	return dataloader.GetJobByJobIdLoader(ctx).Load(obj.Job.ID)
}
func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobsFromDb, err := r.JobsService.GetAllJobs(ctx, skills, statuses)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}
func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	return r.JobsService.GetById(ctx, id)
}
func (r *queryResolver) Jobs(ctx context.Context, filter *model.JobsFilterInput, limit int, after *string) (*model.JobsConnection, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobs, err := r.JobsService.GetAllJobsPaginated(ctx, skills, statuses, limit, after)
	if err != nil {
		return connection, err
	}
	var edges []*gqlmodel.JobEdge

	for i, job := range jobs {
		if i < limit {
			edges = append(edges, &gqlmodel.JobEdge{
				Node:   job,
				Cursor: base64.StdEncoding.EncodeToString([]byte(job.ID)),
			})
		}
	}
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &gqlmodel.JobsConnection{
		//TODO: Implement
		TotalCount: 10,
		Edges:      edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(jobs) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}
func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(id)
}
func (r *queryResolver) Skills(ctx context.Context, query string, limit *int) ([]*model.Skill, error) {
	skills, err := r.SkillsService.GetMatchingSkills(query, limit)
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}

	return result, nil
}
func (r *queryResolver) Search(ctx context.Context, query string, limit *int) (*model.SearchResult, error) {
	//For fetching jobs with title similar to query string
	jobsFromDb, err := r.JobsService.GetByTitle(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var jobs []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		jobs = append(jobs, &tempJob)
	}

	//For fetching users with name similar to query string
	usersFromDb, err := r.UserService.GetByName(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var users []*gqlmodel.User
	for _, dbUser := range usersFromDb {
		var tempUser gqlmodel.User
		tempUser.MapDbToGql(dbUser)
		users = append(users, &tempUser)
	}

	//Search result with jobs and users
	searchResult := gqlmodel.SearchResult{
		Jobs:  jobs,
		Users: users,
	}

	return &searchResult, nil
}
func (r *queryResolver) ViewerNotifications(ctx context.Context, limit int, after *string) (*model.NotificationConnection, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	var afterString *string
	if after != nil {
		if *after == "" {
			return nil, custom_errors.ErrInvalidCursor
		}
		decodedAfter, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, custom_errors.ErrInvalidCursor
		}
		cursorStr := string(decodedAfter)
		afterString = &cursorStr
	}

	notifications, err := r.NotificationsService.GetAllPaginated(user.Id, afterString, limit+1)
	if err != nil {
		return nil, err
	}

	totalNotificationsForUser, err := r.NotificationsService.GetNotificationCountForReceiver(user.Id)

	if err != nil {
		return nil, err
	}

	notificationEdges := make([]*gqlmodel.NotificationEdge, 0)
	for i, notification := range notifications {
		if i < limit {
			notificationEdges = append(notificationEdges, &gqlmodel.NotificationEdge{
				Node:   notification,
				Cursor: base64.StdEncoding.EncodeToString([]byte(notification.ID)),
			})
		}
	}

	var endCursor *string
	if len(notificationEdges) > 0 {
		endCursor = &notificationEdges[len(notificationEdges)-1].Cursor
	}

	return &gqlmodel.NotificationConnection{
		TotalCount: totalNotificationsForUser,
		Edges:      notificationEdges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(notificationEdges) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}
func (r *reviewResolver) CreatedFor(ctx context.Context, obj *model.Review) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedFor)
}
func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.ApplicantID)
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
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}
func (r *jobResolver) CreatedBy(ctx context.Context, obj *model.Job) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}
func (r *jobResolver) Skills(ctx context.Context, obj *model.Job) ([]*model.Skill, error) {
	return dataloader.GetSkillByJobIdLoader(ctx).Load(obj.ID)
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
	return dataloader.GetMilestonesByJobIdLoader(ctx).Load(obj.ID)
}
func (r *jobResolver) Applications(ctx context.Context, obj *model.Job) (*model.Applications, error) {
	return dataloader.GetApplicationsByJobIdLoader(ctx).Load(obj.ID)
}
func (r *jobResolver) ViewerHasApplied(ctx context.Context, obj *model.Job) (bool, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return false, err
	}
	return dataloader.GetViewerHasAppliedLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, user.Id))
}
func (r *milestoneResolver) Job(ctx context.Context, obj *model.Milestone) (*model.Job, error) {
	return r.JobsService.GetById(ctx, obj.JobID)
}
func (r *milestoneResolver) AssignedTo(ctx context.Context, obj *model.Milestone) (*model.User, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.AssignedTo)
}
func (r *milestoneResolver) Review(ctx context.Context, obj *model.Milestone) (*model.Review, error) {
	if obj.AssignedTo == "" {
		return nil, nil
	}
	return dataloader.GetJobMilestoneReviewLoader(ctx).Load(fmt.Sprintf("%v %v", obj.ID, obj.AssignedTo))
}
func (r *milestoneResolver) Skills(ctx context.Context, obj *model.Milestone) ([]*model.Skill, error) {
	return dataloader.GetSkillByMilestoneIdLoader(ctx).Load(obj.ID)
}
func (r *mutationResolver) UpdateProfile(ctx context.Context, user *model.UpdateUserInput) (*model.User, error) {
	return r.UserService.UpdateProfile(ctx, updatedUserDetails)
}
func (r *mutationResolver) CreateJob(ctx context.Context, job *model.CreateJobInput) (*model.Job, error) {
	createdJob, err := r.JobsService.CreateJobs(ctx, job)
	if err != nil {
		return nil, err
	}
	return createdJob[0], nil
}
func (r *mutationResolver) UpdateJob(ctx context.Context, job *model.UpdateJobInput) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.DeleteJob(ctx, jobID)
}
func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*model.Comment, error) {
	return r.JobsService.AddDiscussionToJob(ctx, comment, jobID)
}
func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*model.Comment, error) {
	return r.JobsService.UpdateJobDiscussion(ctx, id, comment)
}
func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*model.Comment, error) {
	return r.JobsService.DeleteJobDiscussion(ctx, id)
}
func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.CreateUserJobApplication(ctx, jobID)
}
func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) ([]*model.Application, error) {
	return r.ApplicationsService.DeleteUserJobApplication(ctx, jobID)
}
func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *model.ApplicationStatus, note *string) ([]*model.Application, error) {
	return r.ApplicationsService.UpdateJobApplicationStatus(ctx, applicantID, jobID, status, note)
}
func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*model.UserAuthenticationPayload, error) {
	// authenticate the user with github and store them in db
	resultUser, err := r.AuthenticationService.AuthenticateAndGetUser(ctx, githubCode)
	if err != nil {
		return nil, err
	}
	//generate a token for the user and return
	authToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	refreshToken, err := resultUser.GenerateAccessToken()

	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	resultPayload := &gqlmodel.UserAuthenticationPayload{
		Profile:      resultUser,
		Token:        *authToken,
		RefreshToken: *refreshToken,
	}
	return resultPayload, nil
}
func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*model.UserAuthenticationPayload, error) {
	// get the claims for the user
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("error while refreshing refreshToken %v", refreshToken)
		log.Println(err)
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid refreshToken signature")
		}
		return nil, err
	}

	if !tkn.Valid {
		return nil, errors.New("refreshToken is not valid")
	}
	// only refresh the refreshToken if it's expiring in 2 minutes
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > (time.Minute * 2) {
		return nil, errors.New("refreshToken can only be refreshed 2 minutes from expiry time")
	}
	//generate a new refreshToken for the user
	gqlUser, err := r.UserService.GetById(ctx, claims.Id)
	if err != nil {
		log.Printf("error getting user from claims for user id %v", claims.Id)
		return nil, err
	}
	newToken, err := gqlUser.GenerateAccessToken()
	newRefreshToken, err := gqlUser.GenerateRefreshToken()
	if err != nil {
		log.Printf("error generating refreshToken for user %+v", gqlUser)
		return nil, err
	}
	return &gqlmodel.UserAuthenticationPayload{
		Profile:      gqlUser,
		Token:        *newToken,
		RefreshToken: *newRefreshToken,
	}, nil
}
func (r *mutationResolver) ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*model.Milestone, error) {
	return r.JobsService.ToggleMilestoneCompleted(ctx, milestoneID)
}
func (r *mutationResolver) ToggleJobCompleted(ctx context.Context, jobID string) (*model.Job, error) {
	return r.JobsService.ToggleJobCompleted(ctx, jobID)
}
func (r *mutationResolver) CreateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, milestoneID string) (*model.Review, error) {
	createdReview, err := r.ReviewsService.ReviewAssignedUser(ctx, review.Rating, review.Remark, milestoneID)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*createdReview)
	return gqlReview, nil
}
func (r *mutationResolver) UpdateMilestonePerformanceReview(ctx context.Context, review model.ReviewInput, id string) (*model.Review, error) {
	updatedReview, err := r.ReviewsService.UpdateReview(ctx, review.Rating, review.Remark, id)
	if err != nil {
		return nil, err
	}
	gqlReview := &gqlmodel.Review{}
	gqlReview.MapDbToGql(*updatedReview)
	return gqlReview, nil
}
func (r *mutationResolver) RestoreJobsBackup(ctx context.Context, jobs []*model.CreateJobInput) ([]*model.Job, error) {
	return r.JobsService.CreateJobs(ctx, jobs...)
}
func (r *mutationResolver) MarkAllViewerNotificationsRead(ctx context.Context) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) MarkViewerNotificationsRead(ctx context.Context, ids []string) ([]*model.NotificationItem, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *notificationItemResolver) Recipient(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Recipient.ID)
}
func (r *notificationItemResolver) Sender(ctx context.Context, obj *model.NotificationItem) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Sender.ID)
}
func (r *notificationItemResolver) Job(ctx context.Context, obj *model.NotificationItem) (*model.Job, error) {
	return dataloader.GetJobByJobIdLoader(ctx).Load(obj.Job.ID)
}
func (r *queryResolver) AllJobs(ctx context.Context, filter *model.JobsFilterInput) ([]*model.Job, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobsFromDb, err := r.JobsService.GetAllJobs(ctx, skills, statuses)
	if err != nil {
		return nil, err
	}

	var result []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		result = append(result, &tempJob)
	}
	return result, nil
}
func (r *queryResolver) Job(ctx context.Context, id string) (*model.Job, error) {
	return r.JobsService.GetById(ctx, id)
}
func (r *queryResolver) Jobs(ctx context.Context, filter *model.JobsFilterInput, limit int, after *string) (*model.JobsConnection, error) {
	var skills []string
	var statuses []string

	if filter == nil {
		filter = &gqlmodel.JobsFilterInput{}
	}

	if filter.Skills != nil && len(filter.Skills) != 0 {
		for _, skill := range filter.Skills {
			skills = append(skills, *skill)
		}
	}

	if filter.Status != nil && len(filter.Status) != 0 {
		for _, status := range filter.Status {
			statuses = append(statuses, status.String())
		}
	}

	jobs, err := r.JobsService.GetAllJobsPaginated(ctx, skills, statuses, limit, after)
	if err != nil {
		return connection, err
	}
	var edges []*gqlmodel.JobEdge

	for i, job := range jobs {
		if i < limit {
			edges = append(edges, &gqlmodel.JobEdge{
				Node:   job,
				Cursor: base64.StdEncoding.EncodeToString([]byte(job.ID)),
			})
		}
	}
	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}
	return &gqlmodel.JobsConnection{
		//TODO: Implement
		TotalCount: 10,
		Edges:      edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(jobs) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}
func (r *queryResolver) User(ctx context.Context, id string, jobsStatusFilter *model.JobStatus) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(id)
}
func (r *queryResolver) Skills(ctx context.Context, query string, limit *int) ([]*model.Skill, error) {
	skills, err := r.SkillsService.GetMatchingSkills(query, limit)
	if err != nil {
		return nil, err
	}

	for _, skill := range skills {
		var gqlSkill gqlmodel.Skill
		gqlSkill.MapDbToGql(*skill)
		result = append(result, &gqlSkill)
	}

	return result, nil
}
func (r *queryResolver) Search(ctx context.Context, query string, limit *int) (*model.SearchResult, error) {
	//For fetching jobs with title similar to query string
	jobsFromDb, err := r.JobsService.GetByTitle(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var jobs []*gqlmodel.Job
	for _, dbJob := range jobsFromDb {
		var tempJob gqlmodel.Job
		tempJob.MapDbToGql(dbJob)
		jobs = append(jobs, &tempJob)
	}

	//For fetching users with name similar to query string
	usersFromDb, err := r.UserService.GetByName(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	var users []*gqlmodel.User
	for _, dbUser := range usersFromDb {
		var tempUser gqlmodel.User
		tempUser.MapDbToGql(dbUser)
		users = append(users, &tempUser)
	}

	//Search result with jobs and users
	searchResult := gqlmodel.SearchResult{
		Jobs:  jobs,
		Users: users,
	}

	return &searchResult, nil
}
func (r *queryResolver) ViewerNotifications(ctx context.Context, limit int, after *string) (*model.NotificationConnection, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	var afterString *string
	if after != nil {
		if *after == "" {
			return nil, custom_errors.ErrInvalidCursor
		}
		decodedAfter, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, custom_errors.ErrInvalidCursor
		}
		cursorStr := string(decodedAfter)
		afterString = &cursorStr
	}

	notifications, err := r.NotificationsService.GetAllPaginated(user.Id, afterString, limit+1)
	if err != nil {
		return nil, err
	}

	totalNotificationsForUser, err := r.NotificationsService.GetNotificationCountForReceiver(user.Id)

	if err != nil {
		return nil, err
	}

	notificationEdges := make([]*gqlmodel.NotificationEdge, 0)
	for i, notification := range notifications {
		if i < limit {
			notificationEdges = append(notificationEdges, &gqlmodel.NotificationEdge{
				Node:   notification,
				Cursor: base64.StdEncoding.EncodeToString([]byte(notification.ID)),
			})
		}
	}

	var endCursor *string
	if len(notificationEdges) > 0 {
		endCursor = &notificationEdges[len(notificationEdges)-1].Cursor
	}

	return &gqlmodel.NotificationConnection{
		TotalCount: totalNotificationsForUser,
		Edges:      notificationEdges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(notificationEdges) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}
func (r *reviewResolver) CreatedFor(ctx context.Context, obj *model.Review) (*model.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedFor)
}
