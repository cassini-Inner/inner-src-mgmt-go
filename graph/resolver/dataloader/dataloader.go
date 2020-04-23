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
	skillsByJobIdLoaderKey       = "skillsByJobIdLoader"
)

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := NewUserByUserIdLoader(db)
		milestoneByJobIdLoader := NewMilestoneByJobIdLoader(db)
		skillsByMilestoneIdLoader := NewSkillByMilestoneIdLoader(db)
		skillsByJobIdLoader := NewSkillByJobIdLoader(db)

		ctxWithuserloader := context.WithValue(r.Context(), userLoaderKey, userloader)
		ctxWithmilestonebyidloader := context.WithValue(ctxWithuserloader, milestoneByJobIdLoaderKey, milestoneByJobIdLoader)
		ctxWithSkillByMilestoneIdLoader := context.WithValue(ctxWithmilestonebyidloader, skillsByMilestoneIdLoaderKey, skillsByMilestoneIdLoader)
		ctxWithSkillByJobIdLoader := context.WithValue(ctxWithSkillByMilestoneIdLoader, skillsByJobIdLoaderKey, skillsByJobIdLoader)

		next.ServeHTTP(w, r.WithContext(ctxWithSkillByJobIdLoader))
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
func GetSkillByJobIdLoader(ctx context.Context) *generated.SkillByJobIdLoader {
	return ctx.Value(skillsByJobIdLoaderKey).(*generated.SkillByJobIdLoader)
}
