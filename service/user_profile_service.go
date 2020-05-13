package service

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

type UserProfileService interface {
	UpdateProfile(ctx context.Context, userDetails *gqlmodel.UpdateUserInput) (*gqlmodel.User, error)
	GetById(ctx context.Context, userId string) (*gqlmodel.User, error)
}
