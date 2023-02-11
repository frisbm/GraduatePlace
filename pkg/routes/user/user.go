package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Manager interface{}

type Router struct {
	manager Manager
}

func NewRouter(manager Manager) *Router {
	return &Router{
		manager: manager,
	}
}

func (r *Router) Public(c chi.Router) {
	c.Get("/users", r.GetUsers)
}

func (r *Router) Private(c chi.Router) {

}

func (r *Router) GetUsers(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hi"))
}
