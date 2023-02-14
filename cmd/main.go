package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"entgo.io/ent/dialect"

	"github.com/MatthewFrisby/thesis-pieces/ent"

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
	db, err := ent.Open(dialect.Postgres, connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer db.Close()

	// Initialize store for handling interactions with the db
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

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
