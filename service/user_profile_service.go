package service

import (
	"context"

	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
)

type UserProfileService interface {
	UpdateProfile(ctx context.Context, userDetails *gqlmodel.UpdateUserInput) (*gqlmodel.User, error)
	GetById(ctx context.Context, userId string) (*gqlmodel.User, error)
	GetByName(ctx context.Context, userName string, limit *int) ([]dbmodel.User, error)
}
