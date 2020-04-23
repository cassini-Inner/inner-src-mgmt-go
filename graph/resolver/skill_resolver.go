package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *skillResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Skill) (*gqlmodel.User, error) {
	return getUserLoader(ctx).Load(obj.CreatedBy)
}
