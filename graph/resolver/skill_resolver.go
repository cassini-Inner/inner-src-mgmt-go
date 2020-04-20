package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *skillResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Skill) (*gqlmodel.User, error) {
	user, err := r.UsersRepo.GetById(obj.CreatedBy)
	if err != nil {
		return nil, err
	}
	var gqlUser gqlmodel.User
	gqlUser.MapDbToGql(*user)
	return &gqlUser, nil
}
