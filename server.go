package main

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres"
	"github.com/go-pg/pg/v9"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	DB := postgres.New(&pg.Options{
		User:     "postgres",
		Password: "root",
		Database: "innersource",
	})
	skillsRepo := postgres.NewSkillsRepo(DB)
	usersRepo := postgres.NewUsersRepo(DB)
	milestonesRepo := postgres.NewMilestonesRepo(DB)
	jobsRepo := postgres.NewJobsRepo(DB)
	discussionsRepo := postgres.NewDiscussionsRepo(DB)
	applicationsRepo := postgres.NewApplicationsRepo(DB)


	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		ApplicationsRepo: applicationsRepo,
		DiscussionsRepo:  discussionsRepo,
		JobsRepo:         jobsRepo,
		MilestonesRepo:   milestonesRepo,
		SkillsRepo:       skillsRepo,
		UsersRepo:        usersRepo,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
