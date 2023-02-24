package document

import (
	"context"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/pagination"

	"github.com/go-chi/chi/v5"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"
	"github.com/MatthewFrisby/thesis-pieces/pkg/models/document"

	"github.com/MatthewFrisby/thesis-pieces/pkg/routes"
)

type Manager interface {
	UploadDocument(ctx context.Context, uploadDocument document.UploadDocument) error
	SearchDocuments(ctx context.Context, searchDocuments document.SearchDocuments) (*pagination.Pagination[document.SearchDocumentsResult], error)
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
	c.Get("/documents/search", routes.WithContext(r.SearchDocuments))
}

func (r *Router) Admin(c chi.Router) {
}

var decoder = schema.NewDecoder()

func (r *Router) PostDocument(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, constants.MB5)

	var uploadDocument document.UploadDocument
	err := routes.ParseMultiPartFormWithFileAndBody(req, &uploadDocument)
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

func (r *Router) SearchDocuments(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	searchDocuments := document.SearchDocuments{}
	err := decoder.Decode(&searchDocuments, req.URL.Query())
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}

	limit, offset := routes.SanitizePagination(searchDocuments.Limit, searchDocuments.Offset)
	searchDocuments.Limit = limit
	searchDocuments.Offset = offset

	documents, err := r.manager.SearchDocuments(ctx, searchDocuments)
	if err != nil {
		routes.Response(w, http.StatusBadRequest, err.Error())
		return
	}
	routes.Response(w, http.StatusCreated, documents)
}
