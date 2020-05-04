package middleware

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

const (
	CurrentUserKey = "currentUserKey"
)

func AuthMiddleware(repo repository.UsersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := parseToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			claims, ok := token.Claims.(jwt.StandardClaims)
			//TODO: Change this to block request if valid token is not supplied
			fmt.Println(claims)
			if !ok || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}
			user, err := repo.GetById(claims.Id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), CurrentUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseToken(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	fmt.Println(cookie.Value)
	var claims jwt.StandardClaims
	token, err := jwt.ParseWithClaims(cookie.Value, &claims,func(token *jwt.Token) (interface{}, error) {
		t := []byte(os.Getenv("JWT_SECRET"))
		return t, nil
	},)
	if err != nil {
		return nil, err
	}
	token.Claims = claims
	return token, errors.Wrap(err, "parseToken error: ")
}

func GetCurrentUserFromContext(ctx context.Context) (*dbmodel.User, error) {
	NoUserInContextError := errors.New("no user in context")

	if ctx.Value(CurrentUserKey) == nil {
		return nil, NoUserInContextError
	}

	user, ok := ctx.Value(CurrentUserKey).(*dbmodel.User)
	if !ok || user.Id == "" {
		return nil, NoUserInContextError
	}
	return user, nil
}
