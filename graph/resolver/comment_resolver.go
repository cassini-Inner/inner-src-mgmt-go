package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *commentResolver) CreatedBy(ctx context.Context, obj *gqlmodel.Comment) (*gqlmodel.User, error) {
		return getUserLoader(ctx).Load(obj.CreatedBy)
}
