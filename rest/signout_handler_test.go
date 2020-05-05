package rest

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignoutHandler(t *testing.T) {
	t.Run("only accepts post method", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		SignoutHandler(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run(`return 400 on a request without "token" cookie`, func(t *testing.T) {
		request, _ := http.NewRequest("POST", "/", nil)
		response := httptest.NewRecorder()

		SignoutHandler(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run(`get "token" cookie from request and expire the cookie`, func(t *testing.T) {
		request, _ := http.NewRequest("POST", "/", nil)
		request.AddCookie(&http.Cookie{
			Name:     "token",
			Value:    "thisisatestcookie",
			Path:     "/",
			Expires:  time.Now().AddDate(0, 0, 14),
			MaxAge:   86400,
			HttpOnly: true,
			Secure:   false,
			Domain:   "localhost",
		})
		response := httptest.NewRecorder()
		SignoutHandler(response, request)
		assert.Equal(t, http.StatusOK, response.Code)

		cookieExpired := false
		for _, cookie := range response.Result().Cookies() {
			if cookie.Name == "token" && cookie.MaxAge < 0 {
				cookieExpired = true
			}
		}

		assert.Equal(t, true, cookieExpired)
	})
}
