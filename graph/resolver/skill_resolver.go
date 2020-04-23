package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
)

func (r *skillResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Skill) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.CreatedBy)
}
