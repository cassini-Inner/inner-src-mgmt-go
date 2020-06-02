package resolver

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	service "github.com/cassini-Inner/inner-src-mgmt-go/service"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ApplicationsRepo      repository.ApplicationsRepo
	DiscussionsRepo       repository.DiscussionsRepo
	JobsRepo              repository.JobsRepo
	SkillsRepo            repository.SkillsRepo
	MilestonesRepo        repository.MilestonesRepo
	JobsService           service.JobsService
	ApplicationsService   service.ApplicationsService
	UserService           service.UserProfileService
	AuthenticationService service.AuthenticationService
	SkillsService         service.SkillsService
}
