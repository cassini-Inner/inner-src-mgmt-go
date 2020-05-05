package impl

import (
	"context"
	"errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/service"
)

type AuthenticationService struct {
	usersRepo    repository.UsersRepo
	oauthService service.OAuthService
}

func NewAuthenticationService(usersRepo repository.UsersRepo, oauthService service.OAuthService) *AuthenticationService {
	return &AuthenticationService{usersRepo: usersRepo, oauthService: oauthService}
}

func (s *AuthenticationService) AuthenticateAndGetUser(ctx context.Context, code string) (*gqlmodel.User, error) {
	_, err := s.oauthService.Authenticate(code)
	if err != nil {
		return nil, err
	}

	fetchedUser, err := s.oauthService.GetUserInfo()
	if err != nil {
		return nil, err
	}

	tx, err := s.usersRepo.BeginTx(ctx)
	// check if the user is signing up for the first time
	usersCount, err := s.usersRepo.CountUsersByGithubId(tx, fetchedUser.GithubId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var gqlUser gqlmodel.User
	var user *dbmodel.User
	switch usersCount {
	case 0:
		user, err = s.usersRepo.CreateNewUser(tx, fetchedUser)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	case 1:
		user, err = s.usersRepo.GetByGithubId(fetchedUser.GithubId.String)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		gqlUser.MapDbToGql(*user)
		return &gqlUser, nil
	default:
		_ = tx.Rollback()
		return nil, errors.New("multiple users by the same name exist in database")
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
