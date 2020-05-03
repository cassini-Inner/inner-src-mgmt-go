package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AuthenticationService struct {
	usersRepo repository.UsersRepo
}

func NewAuthenticationService( usersRepo repository.UsersRepo) *AuthenticationService {
	return &AuthenticationService{ usersRepo: usersRepo}
}

func (s *AuthenticationService) AuthenticateAndGetUser(ctx context.Context, githubCode string) (*gqlmodel.User, error) {
	accessToken, err := s.getAccessTokenFromCode(githubCode)
	if err != nil {
		return nil, err
	}

	fetchedUser, err := s.getUserInformationFromToken(accessToken)
	if err != nil {
		return nil, err
	}

	tx, err := s.usersRepo.BeginTx(ctx)
	// check if the user is signing up for the first time
	usersCount, err := s.usersRepo.CountUsersByGithubId(fetchedUser.GithubId, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var gqlUser gqlmodel.User
	var user *dbmodel.User
	switch usersCount {
	case 0:
		user, err = s.usersRepo.CreateNewUser(fetchedUser, tx)
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

func (s *AuthenticationService) getAccessTokenFromCode(githubCode string) (string, error) {
	urlStr := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v", os.Getenv("client_id"),
		os.Getenv("client_secret"),
		githubCode,
	)

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(""))
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	data, err := s.parseJsonFromResponse(response.Body)
	if err != nil {
		return "", err
	}
	//TODO: Error message is ambiguous, use the error message provided by github
	if response.StatusCode != 200 {
		return "", errors.New("could not authenticate with github")
	}
	accessToken, ok := data["access_token"].(string)
	fmt.Println(accessToken)
	if !ok {
		return "", errors.New("could not get access_token from github auth response, token expired or invalid")
	}

	return accessToken, nil
}

func (s *AuthenticationService) getUserInformationFromToken(accessToken string) (*dbmodel.User, error) {
	client := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", strings.NewReader(""))
	request.Header.Set("Authorization", fmt.Sprintf(
		"token %v", accessToken))

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return nil, err
	}

	result, err := s.parseJsonFromResponse(response.Body)
	if err != nil {
		return nil, err
	}

	user := &dbmodel.User{
		Email:      dbmodel.ToNullString(result["email"]),
		Name:       dbmodel.ToNullString(result["name"]),
		Bio:        dbmodel.ToNullString(result["bio"]),
		PhotoUrl:   dbmodel.ToNullString(result["avatar_url"]),
		GithubUrl:  dbmodel.ToNullString(result["html_url"]),
		GithubId:   dbmodel.ToNullString(result["id"]),
		GithubName: dbmodel.ToNullString(result["login"]),
	}
	return user, nil
}

func (s *AuthenticationService) parseJsonFromResponse(responseBody io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(responseBody)
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
