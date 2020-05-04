package impl

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

type AuthenticationService interface {
	AuthenticateAndGetUser(ctx context.Context, githubCode string) (*gqlmodel.User, error)
}