package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver"
	CustomMiddlewares "github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	//TODO: Make this more secure
	_ = os.Setenv("jwt_secret", "innersource_jwt_secret_key")
	_ = os.Setenv("client_id", "5a4ff35b849d9cc3cab7")
	_ = os.Setenv("client_secret", "f94c5d74e099ed894f88ac6c75ac19c4c3194427")

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
		AllowedOrigins:   []string{"http://localhost:8081", "http://localhost:8080", "http://localhost:3000"},
		AllowCredentials: true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(CustomMiddlewares.AuthMiddleware(*usersRepo))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query",resolver.DataloaderMiddleware(DB, srv) )

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
