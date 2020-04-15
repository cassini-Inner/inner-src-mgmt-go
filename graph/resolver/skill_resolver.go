package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *skillResolver) CreatedBy(ctx context.Context, obj *model.Skill) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
