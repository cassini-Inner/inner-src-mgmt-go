package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-source-mgmt-srv/graph/model"
)

func (r *commentResolver) Job(ctx context.Context, obj *model.Comment) (*model.Job, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commentResolver) CreatedBy(ctx context.Context, obj *model.Comment) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
