package dataloader

import (
	"context"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/jmoiron/sqlx"
	"net/http"
)

const (
	userLoaderKey                = "userloader"
	milestoneByJobIdLoaderKey    = "milestoneByIdLoader"
	skillsByMilestoneIdLoaderKey = "skillsByMilestoneIdLoader"
)

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := NewUserByUserIdLoader(db)
		milestoneByJobIdLoader := NewMilestoneByJobIdLoader(db)
		skillsByMilestoneIdLoader := NewSkillByMilestoneIdLoader(db)

		ctxWithuserloader := context.WithValue(r.Context(), userLoaderKey, userloader)
		ctxWithmilestonebyidloader := context.WithValue(ctxWithuserloader, milestoneByJobIdLoaderKey, milestoneByJobIdLoader)
		ctxWithSkillByMilestoneIdLoader := context.WithValue(ctxWithmilestonebyidloader, skillsByMilestoneIdLoaderKey, skillsByMilestoneIdLoader)

		next.ServeHTTP(w, r.WithContext(ctxWithSkillByMilestoneIdLoader))
	})
}

func GetUserByUserIdLoader(ctx context.Context) *generated.UserLoader {
	userLoader := ctx.Value(userLoaderKey).(*generated.UserLoader)
	return userLoader
}

func GetMilestonesByJobIdLoader(ctx context.Context) *generated.MilestoneByJobIdLoader {
	milestoneByJobIdLoader := ctx.Value(milestoneByJobIdLoaderKey).(*generated.MilestoneByJobIdLoader)
	return milestoneByJobIdLoader
}

func GetSkillByMilestoneIdLoader(ctx context.Context) *generated.SkillByMilestoneIdLoader {
	return ctx.Value(skillsByMilestoneIdLoaderKey).(*generated.SkillByMilestoneIdLoader)
}
