package resolver

import (
	"context"
	"errors"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, user *gqlmodel.UpdateUserInput) (*gqlmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUserProfile(ctx context.Context, user *gqlmodel.CreateUserInput) (*gqlmodel.User, error) {
	var dbuser *dbmodel.User
	var gqluser gqlmodel.User
	dbuser, err := r.UsersRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	gqluser.MapDbToGql(*dbuser)
	return &gqluser, err
}

func (r *mutationResolver) CreateJob(ctx context.Context, job *gqlmodel.CreateJobInput) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJob(ctx context.Context, job *gqlmodel.UpdateJobInput) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJob(ctx context.Context, jobID string) (*gqlmodel.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddCommentToJob(ctx context.Context, comment string, jobID string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateComment(ctx context.Context, id string, comment string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteCommment(ctx context.Context, id string) (*gqlmodel.Comment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateJobApplication(ctx context.Context, jobID string) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteJobApplication(ctx context.Context, jobID string) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateJobApplication(ctx context.Context, applicantID string, jobID string, status *gqlmodel.ApplicationStatus) (*gqlmodel.Application, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Authenticate(ctx context.Context, githubCode string) (*gqlmodel.UserAuthenticationPayload, error) {
	// authenticate the user with github and store them in db
	user, err := r.UsersRepo.AuthenticateAndGetUser(githubCode)
	if err != nil {
		return nil, err
	}

	// map db user to graphql model
	var resultUser gqlmodel.User
	resultUser.MapDbToGql(*user)
	//generate a token for the user and return
	authToken, err := resultUser.GenerateJwtToken()
	if err != nil {
		log.Println(err)
		return nil, errors.New("something went wrong")
	}
	resultPayload := &gqlmodel.UserAuthenticationPayload{
		Profile: &resultUser,
		Token:   *authToken,
	}
	return resultPayload, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (*gqlmodel.UserAuthenticationPayload, error) {
	// get the claims for the user
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("error while refreshing token %v", token)
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("something went wrong")
	}
	if !tkn.Valid {
		return nil, errors.New("token is not valid")
	}
	// only refresh the token if it's expiring in 2 minutes
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > (time.Minute * 2) {
		return nil, errors.New("token can only be refreshed 2 minutes from expiry time")
	}
	//generate a new token for the user
	user, err := r.UsersRepo.GetById(claims.Id)
	if err != nil {
		log.Printf("error getting user from claims for user id %v", claims.Id)
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	newToken, err := gqlUser.GenerateJwtToken()
	if err != nil {
		log.Printf("error generating token for user %+v", gqlUser)
		return nil, err
	}
	return &gqlmodel.UserAuthenticationPayload{
		Profile: &gqlUser,
		Token:   *newToken,
	}, nil
}
