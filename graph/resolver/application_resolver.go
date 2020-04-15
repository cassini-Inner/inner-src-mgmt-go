package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-src-mgmt-go/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
