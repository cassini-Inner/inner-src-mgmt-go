package resolver

import "github.com/cassini-Inner/inner-src-mgmt-go/postgres"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ApplicationsRepo *postgres.ApplicationsRepo
	DiscussionsRepo  *postgres.DiscussionsRepo
	JobsRepo         *postgres.JobsRepo
	MilestonesRepo   *postgres.MilestonesRepo
	SkillsRepo       *postgres.SkillsRepo
	UsersRepo        *postgres.UsersRepo
}
