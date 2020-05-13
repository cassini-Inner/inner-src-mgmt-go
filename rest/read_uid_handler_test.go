package rest

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUIDFromCookie(t *testing.T) {
	t.Run("no cookie in request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/read-cookie", nil)
		response := httptest.NewRecorder()
		GetUIDFromCookie(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("cookie exists", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/read-cookie", nil)

		token, _ := mockUser.GenerateAccessToken()

		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    *token,
			HttpOnly: true,
		})
		response := httptest.NewRecorder()
		GetUIDFromCookie(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		user := struct {
			UserId string `json:"user_id"`
		}{}

		err := json.NewDecoder(response.Body).Decode(&user)
		assert.Nil(t, err)
		assert.Equal(t, mockUser.ID, user.UserId)
	})
}
