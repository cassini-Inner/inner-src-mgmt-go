package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
	CustomMiddlewares "github.com/cassini-Inner/inner-src-mgmt-go/middleware"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository/impl"
	"github.com/cassini-Inner/inner-src-mgmt-go/rest"
	impl2 "github.com/cassini-Inner/inner-src-mgmt-go/service/impl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const defaultPort = "8080"

var (
	ErrNullDB = errors.New("no db supplied")
)

func SetupRouter(DB *sqlx.DB) (*chi.Mux, error) {
	if DB == nil {
		return nil, ErrNullDB
	}

	skillsRepo := impl.NewSkillsRepo(DB)
	usersRepo := impl.NewUsersRepo(DB)
	jobsRepo := impl.NewJobsRepo(DB)
	discussionsRepo := impl.NewDiscussionsRepo(DB)
	applicationsRepo := impl.NewApplicationsRepo(DB)

	jobsService := impl2.NewJobsService(jobsRepo, skillsRepo, discussionsRepo, applicationsRepo)
	applicationsService := impl2.NewApplicationsService(jobsRepo, applicationsRepo)
	userService := impl2.NewUserProfileService(usersRepo, skillsRepo)
	githubOauthService := impl2.NewGithubOauthService()
	authService := impl2.NewAuthenticationService(usersRepo, githubOauthService)
	skillsService := impl2.NewSkillsService(skillsRepo)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		ApplicationsRepo:      applicationsRepo,
		DiscussionsRepo:       discussionsRepo,
		JobsRepo:              jobsRepo,
		SkillsRepo:            skillsRepo,
		JobsService:           jobsService,
		ApplicationsService:   applicationsService,
		UserService:           userService,
		AuthenticationService: authService,
		SkillsService:         skillsService,
	}}))

	restAuthHandler := rest.NewAuthenticationHandler(authService)

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081", "http://localhost:8080", "http://localhost:3000"},
		AllowedMethods:   []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodDelete, http.MethodConnect},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(CustomMiddlewares.AuthMiddleware(usersRepo))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", dataloader.DataloaderMiddleware(DB, srv))
	router.Handle("/authenticate", restAuthHandler)
	router.HandleFunc("/logout", rest.SignoutHandler)
	router.HandleFunc("/read-cookie", rest.GetUIDFromCookie)
	return router, nil
}

func main() {
	//TODO: Make this more secure
	_ = os.Setenv("jwt_secret", "innersource_jwt_secret_key")
	_ = os.Setenv("client_id", "5a4ff35b849d9cc3cab7")
	_ = os.Setenv("client_secret", "f94c5d74e099ed894f88ac6c75ac19c4c3194427")
	_ = os.Setenv("db_conn_string", "host=localhost port=5432 user=postgres dbname=innersource password=a sslmode=disable")

	DB, err := sqlx.Connect("postgres", os.Getenv("db_conn_string"))
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router, err := SetupRouter(DB)
	if err != nil {
		panic(err)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
