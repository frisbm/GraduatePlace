package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"entgo.io/ent/entc/integration/ent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"

	"github.com/MatthewFrisby/thesis-pieces/pkg/config"
	"github.com/MatthewFrisby/thesis-pieces/pkg/stack/user"
)

type Route interface {
	Public(c chi.Router)
	Private(c chi.Router)
}

func main() {
	// Create appConfig for further reading configuration variables
	_ = config.NewConfig("config.json")

	// Open sqlite3 db using ent, enable foreign-keys
	db, err := ent.Open("sqlite3", "file:store.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer db.Close()

	// Initialize store for handling interactions with the db
	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	userStack := user.NewStack(db)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	public := r.Group(nil)

	private := r.Group(nil)
	private.Use()

	routes := []Route{
		userStack.Router,
	}

	for _, r := range routes {
		r.Private(private)
		r.Public(public)
	}
	http.ListenAndServe(":8080", r)
}
