package resolver

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
	service "github.com/cassini-Inner/inner-src-mgmt-go/service"
	impl "github.com/cassini-Inner/inner-src-mgmt-go/service/impl"
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
	JobsService           *impl.JobsService
	ApplicationsService   *impl.ApplicationsService
	UserService           *impl.UserProfileService
	AuthenticationService service.AuthenticationService
	SkillsService         *impl.SkillsService
}
