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
	GetUsers(ctx context.Context) ([]*user.GetUser, error)
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
	c.Post("/user/register", routes.WithContext(r.RegisterUser))
	c.Post("/user/login", routes.WithContext(r.LoginUser))
	c.Post("/user/refresh", routes.WithContext(r.RefreshUser))
}

func (r *Router) Private(c chi.Router) {
	c.Get("/user", routes.WithContext(r.GetUser))
}

func (r *Router) Admin(c chi.Router) {
	c.Get("/users", routes.WithContext(r.GetUsers))
}

func (r *Router) RegisterUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var registerUser user.RegisterUser
	err := json.NewDecoder(req.Body).Decode(&registerUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	err = r.manager.RegisterUser(ctx, registerUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusCreated, "Registered Successfully")
}

func (r *Router) LoginUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var loginUser user.LoginUser
	err := json.NewDecoder(req.Body).Decode(&loginUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := r.manager.LoginUser(ctx, loginUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusOK, token)
}

func (r *Router) RefreshUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var refreshUser user.RefreshUser
	err := json.NewDecoder(req.Body).Decode(&refreshUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := r.manager.RefreshUser(ctx, refreshUser)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusOK, token)
}

func (r *Router) GetUser(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	user, err := r.manager.GetUser(ctx)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusOK, user)
}

func (r *Router) GetUsers(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	users, err := r.manager.GetUsers(ctx)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusOK, users)
}
