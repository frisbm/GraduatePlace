package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/routes"
)

type Manager interface {
	RegisterUser(ctx context.Context, registerUser user.RegisterUser) error
	LoginUser(ctx context.Context, loginUser user.LoginUser) (*user.AuthTokens, error)
	RefreshUser(ctx context.Context, refreshUser user.RefreshUser) (*user.AuthTokens, error)
	GetUser(ctx context.Context) (*user.GetUser, error)
	GetUsers(ctx context.Context) (*user.GetUsers, error)
}

type Router struct {
	manager Manager
}

func NewRouter(manager Manager) *Router {
	return &Router{
		manager: manager,
	}
}

func (r *Router) Public(c chi.Router) {
	c.Post("/user/register", routes.WithContext(r.Register))
	c.Post("/user/login", routes.WithContext(r.Login))
}

func (r *Router) Private(c chi.Router) {
	c.Get("/user", routes.WithContext(r.GetUser))
}

func (r *Router) Admin(c chi.Router) {
	c.Get("/users", routes.WithContext(r.GetUsers))
}

func (r *Router) Register(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerUser user.RegisterUser
	err := json.NewDecoder(req.Body).Decode(&registerUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = r.manager.RegisterUser(ctx, registerUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Success")
}

func (r *Router) Login(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginUser user.LoginUser
	err := json.NewDecoder(req.Body).Decode(&loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	token, err := r.manager.LoginUser(ctx, loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (r *Router) GetUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	user, err := r.manager.GetUser(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (r *Router) GetUsers(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	users, err := r.manager.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
