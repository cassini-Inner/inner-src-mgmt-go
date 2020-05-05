package impl

import (
	"encoding/json"
	"fmt"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GithubOauthService struct {
	accessToken string
}

func NewGithubOauthService() *GithubOauthService {
	return &GithubOauthService{}
}

func (g *GithubOauthService) Authenticate(code string) (string, error) {
	urlStr := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%v&client_secret=%v&code=%v", os.Getenv("client_id"),
		os.Getenv("client_secret"),
		code,
	)

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(""))
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	data, err := g.parseJsonFromResponse(response.Body)
	if err != nil {
		return "", err
	}
	//TODO: Error message is ambiguous, use the error message provided by github
	if response.StatusCode != 200 {
		return "", customErrors.ErrCodeExpired
	}
	accessToken, ok := data["access_token"].(string)
	if !ok {
		return "", customErrors.ErrInvalidAuthResponse
	}
	g.accessToken = accessToken
	return accessToken, nil
}

func (g GithubOauthService) GetUserInfo() (*dbmodel.User, error) {
	if g.accessToken == "" {
		return nil, customErrors.ErrOauthClientNotAuthenticated
	}

	client := http.Client{}
	request, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", strings.NewReader(""))
	request.Header.Set("Authorization", fmt.Sprintf(
		"token %v", g.accessToken))

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return nil, err
	}

	result, err := g.parseJsonFromResponse(response.Body)
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

func (g GithubOauthService) parseJsonFromResponse(responseBody io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(responseBody)
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
