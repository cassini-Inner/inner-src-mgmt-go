package resolver

import "github.com/cassini-Inner/inner-src-mgmt-go/postgres"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	applicationsRepo *postgres.ApplicationsRepo
	discussionsRepo  *postgres.DiscussionsRepo
	jobsRepo         *postgres.JobsRepo
	milestonesRepo   *postgres.MilestonesRepo
	skillsRepo       *postgres.SkillsRepo
	usersRepo        *postgres.UsersRepo
}
