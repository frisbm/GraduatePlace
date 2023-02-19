package document

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"
	"github.com/MatthewFrisby/thesis-pieces/pkg/models/document"

	"github.com/MatthewFrisby/thesis-pieces/pkg/routes"
)

type Manager interface {
	UploadDocument(ctx context.Context, uploadDocument document.UploadDocument) error
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
}

func (r *Router) Private(c chi.Router) {
	c.Post("/document", routes.WithContext(r.PostDocument))
}

func (r *Router) Admin(c chi.Router) {
}

func (r *Router) PostDocument(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, constants.MB5)

	var uploadDocument document.UploadDocument
	err := routes.ParseMultiPartFormWithFileAndBody(req, &uploadDocument)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	err = r.manager.UploadDocument(ctx, uploadDocument)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusCreated, "Document Uploaded Successfully")
}
