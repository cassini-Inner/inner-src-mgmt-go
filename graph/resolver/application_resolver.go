package resolver

import (
	"context"
	"fmt"
	"github.com/cassini-inner/inner-source-mgmt-srv/graph/model"
)

func (r *applicationResolver) Applicant(ctx context.Context, obj *model.Application) (*model.User, error) {
	return r.UserRepo.getApplicantById(obj.ID)
}
