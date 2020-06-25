package main

import (
	"errors"
	"fmt"
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
	milestonesRepo := impl.NewMilestonesRepoImpl(DB)
	reviewsRepo := impl.NewReviewRepoImpl(DB)
	notificationRepo := impl.NewNotificationsRepo(DB)

	jobsService := impl2.NewJobsService(jobsRepo, skillsRepo, discussionsRepo, applicationsRepo, milestonesRepo, notificationRepo)
	applicationsService := impl2.NewApplicationsService(jobsRepo, applicationsRepo, milestonesRepo, notificationRepo)
	userService := impl2.NewUserProfileService(usersRepo, skillsRepo)
	githubOauthService := impl2.NewGithubOauthService()
	authService := impl2.NewAuthenticationService(usersRepo, githubOauthService)
	skillsService := impl2.NewSkillsService(skillsRepo)
	reviewsService := impl2.NewReviewsService(reviewsRepo, jobsRepo, milestonesRepo)
	notificationsService := impl2.NewNotificationsService(notificationRepo)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &resolver.Resolver{
				ApplicationsRepo:      applicationsRepo,
				DiscussionsRepo:       discussionsRepo,
				JobsRepo:              jobsRepo,
				SkillsRepo:            skillsRepo,
				JobsService:           jobsService,
				ApplicationsService:   applicationsService,
				UserService:           userService,
				AuthenticationService: authService,
				SkillsService:         skillsService,
				ReviewsService:        reviewsService,
				NotificationsService:  notificationsService,
			},
			},
		),
	)

	restAuthHandler := rest.NewAuthenticationHandler(authService)

	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://10.176.20.185", "http://localhost:3000", "http://localhost:8080", "http://10.176.5.190"},
		AllowedMethods: []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodDelete, http.MethodConnect,
			http.MethodTrace, http.MethodHead, http.MethodPatch},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(CustomMiddlewares.AuthMiddleware(usersRepo))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", dataloader.DataloaderMiddleware(DB, srv))
	router.Handle("/authenticate", restAuthHandler)
	router.HandleFunc("/logout", rest.SignoutHandler)
	router.HandleFunc("/read-cookie", rest.GetUIDFromCookie)
	router.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("pong"))
		writer.WriteHeader(http.StatusOK)
	})
	return router, nil
}

func main() {
	DB, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			os.Getenv("host"),
			os.Getenv("port"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_DB"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("sslmode"),
		),
	)
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	port := os.Getenv("SERVER_PORT")
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
