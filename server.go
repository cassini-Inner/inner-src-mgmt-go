package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	DB, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres dbname=innersource password=root sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer DB.Close()

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
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081","http://localhost:8080"},
		AllowCredentials: true,
	}).Handler)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}


