package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *commentResolver) CreatedBy(ctx context.Context, obj *model.Comment) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
