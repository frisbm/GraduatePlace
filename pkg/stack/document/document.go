package document

import (
	documentManager "github.com/frisbm/graduateplace/pkg/managers/document"
	documentRoutes "github.com/frisbm/graduateplace/pkg/routes/document"
	"github.com/frisbm/graduateplace/pkg/services/s3"
	"github.com/frisbm/graduateplace/pkg/store"
	documentStore "github.com/frisbm/graduateplace/pkg/store/document"
	"github.com/frisbm/graduateplace/pkg/tasks"
)

type Stack struct {
	Store   *documentStore.Store
	Tasks   *tasks.TaskManager
	Manager *documentManager.Manager
	Router  *documentRoutes.Router
}

func NewStack(queries *store.Queries, s3 *s3.S3, tasks *tasks.TaskManager) *Stack {
	store := documentStore.NewStore(queries)
	manager := documentManager.NewManager(store, s3, tasks)
	router := documentRoutes.NewRouter(manager)
	return &Stack{
		Store:   store,
		Manager: manager,
		Router:  router,
	}
}
