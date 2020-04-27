package resolver

import (
	"context"
	"errors"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var (
	ErrUserNotOwner                   = errors.New("current user is not owner of this entity, and hence cannot modify it")
	ErrNoEntityMatchingId             = errors.New("no entity found that matches given id")
	ErrOwnerApplyToOwnJob             = errors.New("owner cannot apply to their job")
	ErrApplicationWithdrawnOrRejected = errors.New("owner cannot modify applications with withdrawn status")
	ErrInvalidNewApplicationState     = errors.New("owner cannot move application status to withdrawn or pending")
	ErrJobAlreadyCompleted            = errors.New("job is already completed")
	ErrEntityDeleted                  = errors.New("entity was deleted")
	ErrUserNotAuthenticated           = errors.New("unauthorized request")
)

func (r *mutationResolver) UpdateProfile(ctx context.Context, updatedUserDetails *gqlmodel.UpdateUserInput) (*gqlmodel.User, error) {
	return r.UserService.UpdateProfile(ctx, updatedUserDetails)
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *gqlmodel.CreateJobInput) (*gqlmodel.Job, error) {
	return r.JobsService.CreateJob(ctx, job)
}

func (r *mutationResolver) ToggleMilestoneCompleted(ctx context.Context, milestoneID string) (*gqlmodel.Milestone, error) {
	return r.JobsService.ToggleMilestoneCompleted(ctx, milestoneID)
}

func (r *mutationResolver) ToggleJobCompleted(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	return r.JobsService.ToggleJobCompleted(ctx, jobID)
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *gqlmodel.UpdateJobInput) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	return r.JobsService.DeleteJob(ctx, jobID)
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*gqlmodel.Comment, error) {
	return r.JobsService.AddDiscussionToJob(ctx, comment, jobID)
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*gqlmodel.Comment, error) {
	return r.JobsService.UpdateJobDiscussion(ctx, id, comment)
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (*gqlmodel.Comment, error) {
	return r.JobsService.DeleteJobDiscussion(ctx, id)
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) ([]*gqlmodel.Application, error) {
	return r.ApplicationsService.CreateUserJobApplication(ctx, jobID)
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) ([]*gqlmodel.Application, error) {
	return r.ApplicationsService.DeleteUserJobApplication(ctx, jobID)
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *gqlmodel.ApplicationStatus, note *string) (result []*gqlmodel.Application, err error) {
	return r.ApplicationsService.UpdateJobApplicationStatus(ctx, applicantID, jobID, status, note)
}

func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*gqlmodel.UserAuthenticationPayload, error) {
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

func (r *mutationResolver) RefreshToken(ctx context.Context, refreshToken string) (*gqlmodel.UserAuthenticationPayload, error) {
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
