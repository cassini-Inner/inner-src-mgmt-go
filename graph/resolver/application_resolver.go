package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	return r.UsersRepo.getUserById(obj.ID)
}
