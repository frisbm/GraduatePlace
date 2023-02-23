package document

import (
	documentManager "github.com/MatthewFrisby/thesis-pieces/pkg/managers/document"
	documentRoutes "github.com/MatthewFrisby/thesis-pieces/pkg/routes/document"
	"github.com/MatthewFrisby/thesis-pieces/pkg/services/s3"
	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
	documentStore "github.com/MatthewFrisby/thesis-pieces/pkg/store/document"
	"github.com/MatthewFrisby/thesis-pieces/pkg/tasks"
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
