package rest

import (
	"encoding/json"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	"github.com/cassini-Inner/inner-src-mgmt-go/service"
	_ "io/ioutil"
	"log"
	"net/http"
	"time"
)

type AuthenticationHandler struct {
	authService service.AuthenticationService
}

type authRequestBody struct {
	Code string `json:"code"`
}

func NewAuthenticationHandler(authService service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{authService: authService}
}

// accepts a github code, the id of the user from data after calling github service
// and sets the cookie by generating a jwt token
func (a AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}
	body := &authRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if body.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("invalid code")
		return
	}
	user, err := a.authService.AuthenticateAndGetUser(r.Context(), body.Code)
	if err != nil {
		log.Printf("AuthService Error: %v", err)
		if err == customErrors.ErrCodeExpired {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := user.GenerateAccessToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	cookie := http.Cookie{Name: "token",
		Value:    *token,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 14),
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		Domain:   r.RequestURI,
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte(`{"success": "true"}`))
	w.WriteHeader(http.StatusOK)
}
