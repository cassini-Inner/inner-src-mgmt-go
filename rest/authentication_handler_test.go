package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	mockUser = &gqlmodel.User{
		ID: "1",
	}
	AuthAndGetUser = func() (*gqlmodel.User, error) {
		return mockUser, nil
	}
)

type AuthServiceMock struct {
}

func (a AuthServiceMock) AuthenticateAndGetUser(ctx context.Context, githubCode string) (*gqlmodel.User, error) {
	return AuthAndGetUser()
}

func TestAuthenticationHandler_ServeHTTP(t *testing.T) {
	srv := NewAuthenticationHandler(&AuthServiceMock{})
	t.Run("only accepts post method", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/authenticate", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
	})

	t.Run("rejects other methods", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/authenticate", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
	})

	t.Run("returns 400 on nil body", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/authenticate", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
	})

	t.Run("returns 400 on empty code", func(t *testing.T) {
		body, err := json.Marshal(map[string]string{
			"code": "",
		})

		assert.Nil(t, err)
		assert.NotNil(t, body)
		request := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
	})

	t.Run("add jwt cookie to response", func(t *testing.T) {
		body, err := json.Marshal(map[string]string{
			"code": "12345",
		})

		assert.Nil(t, err)
		assert.NotNil(t, body)
		request := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

		var found bool
		var tokenCookie = &http.Cookie{}
		for _, cookie := range response.Result().Cookies() {
			if cookie.Name == "token" {
				found = true
				tokenCookie = cookie
				break
			}
		}
		assert.Nil(t, err)
		assert.Equal(t, true, found)

		var claims jwt.StandardClaims
		_, err = jwt.ParseWithClaims(tokenCookie.Value, &claims, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		assert.Nil(t, err)
		assert.Equal(t, claims.Id, mockUser.ID)
	})

	t.Run("auth service fails to get create", func(t *testing.T) {
		prev := AuthAndGetUser
		defer func() { AuthAndGetUser = prev }()
		AuthAndGetUser = func() (user *gqlmodel.User, err error) {
			return nil, errors.New("failed to create user")
		}

		body, err := json.Marshal(map[string]string{
			"code": "12345",
		})

		assert.Nil(t, err)
		assert.NotNil(t, body)
		request := httptest.NewRequest(http.MethodPost, "/authenticate", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
