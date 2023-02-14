package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"

	"github.com/go-chi/jwtauth/v5"

	"github.com/MatthewFrisby/thesis-pieces/pkg/utils/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"

	"github.com/MatthewFrisby/thesis-pieces/pkg/config"
	"github.com/MatthewFrisby/thesis-pieces/pkg/stack/user"
)

const (
	API_PATH = "/api"
)

type Route interface {
	Public(c chi.Router)
	Private(c chi.Router)
	Admin(c chi.Router)
}

func main() {
	// Create appConfig for further reading configuration variables
	config := config.NewConfig("config.json")

	connectionString := fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBSSLMode,
	)
	// Open postgres
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer database.Close()

	if err := goose.Up(database, "database/migrations"); err != nil {
		log.Fatalf("failed migrating with goose: %v", err)
	}

	db := store.New(database)

	authService := auth.NewAuthService(config.JWTSecretKey)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	public := r.Group(nil)

	authMiddleware := auth.NewAuthMiddleware(db, authService)

	private := r.Group(nil)
	private.Use(jwtauth.Verifier(authService.GetTokenAuth()))
	private.Use(authMiddleware.Private)

	admin := r.Group(nil)
	admin.Use(jwtauth.Verifier(authService.GetTokenAuth()))
	admin.Use(authMiddleware.Admin)

	userStack := user.NewStack(db, authService)

	routes := []Route{
		userStack.Router,
	}

	for _, route := range routes {
		route.Public(public)
		route.Private(private)
		route.Admin(admin)
	}

	router := chi.NewRouter()
	router.Mount(API_PATH, r)
	http.ListenAndServe(":8080", router)
}
