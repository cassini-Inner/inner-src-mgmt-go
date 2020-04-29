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
	applicationsByJobIdLoaderKey = "applicationsByJobIdLoaderKey"
	dataloadersKey               = "dataloadersKey"
	viewerHasAppliedLoaderKey = "viewerHasAppliedLoaderKey"
)

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := NewUserByUserIdLoader(db)
		milestoneByJobIdLoader := NewMilestoneByJobIdLoader(db)
		skillsByMilestoneIdLoader := NewSkillByMilestoneIdLoader(db)
		skillsByJobIdLoader := NewSkillByJobIdLoader(db)
		viewerHasAppliedLoader := NewViewerHasAppliedByUserIdLoader(db)
		applicationsByJobIdLoader := NewApplicationByJobIdLoader(db)

		loaderMap := make(map[string]interface{})
		loaderMap[userLoaderKey] = userloader
		loaderMap[milestoneByJobIdLoaderKey] = milestoneByJobIdLoader
		loaderMap[skillsByMilestoneIdLoaderKey] = skillsByMilestoneIdLoader
		loaderMap[viewerHasAppliedLoaderKey] = viewerHasAppliedLoader
		loaderMap[skillsByJobIdLoaderKey] = skillsByJobIdLoader
		loaderMap[applicationsByJobIdLoaderKey] = applicationsByJobIdLoader


		ctxWithLoaders := context.WithValue(r.Context(), dataloadersKey, loaderMap)

		next.ServeHTTP(w, r.WithContext(ctxWithLoaders))
	})
}

func GetUserByUserIdLoader(ctx context.Context) *generated.UserLoader {
	userLoader := ctx.Value(dataloadersKey).(map[string]interface{})[userLoaderKey].(*generated.UserLoader)
	return userLoader
}

func GetMilestonesByJobIdLoader(ctx context.Context) *generated.MilestoneByJobIdLoader {
	milestoneByJobIdLoader := ctx.Value(dataloadersKey).(map[string]interface{})[milestoneByJobIdLoaderKey].(*generated.MilestoneByJobIdLoader)
	return milestoneByJobIdLoader
}

func GetSkillByMilestoneIdLoader(ctx context.Context) *generated.SkillByMilestoneIdLoader {
	return ctx.Value(dataloadersKey).(map[string]interface{})[skillsByMilestoneIdLoaderKey].(*generated.SkillByMilestoneIdLoader)
}

func GetSkillByJobIdLoader(ctx context.Context) *generated.SkillByJobIdLoader {
	return ctx.Value(dataloadersKey).(map[string]interface{})[skillsByJobIdLoaderKey].(*generated.SkillByJobIdLoader)
}

func GetApplicationsByJobIdLoader(ctx context.Context) *generated.ApplicationsByJobIdLoader {
	return ctx.Value(dataloadersKey).(map[string]interface{})[applicationsByJobIdLoaderKey].(*generated.ApplicationsByJobIdLoader)
}

func GetViewerHasAppliedLoader(ctx context.Context) *generated.ViewerHasAppliedLoader{
	return ctx.Value(dataloadersKey).(map[string]interface{})[viewerHasAppliedLoaderKey].(*generated.ViewerHasAppliedLoader);
}