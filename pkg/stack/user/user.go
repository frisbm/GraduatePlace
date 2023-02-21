package user

import (
	userManager "github.com/MatthewFrisby/thesis-pieces/pkg/managers/user"
	userRoutes "github.com/MatthewFrisby/thesis-pieces/pkg/routes/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/services/auth"
	"github.com/MatthewFrisby/thesis-pieces/pkg/services/s3"
	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
	userStore "github.com/MatthewFrisby/thesis-pieces/pkg/store/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/tasks"
)

type Stack struct {
	Store   *userStore.Store
	Tasks   *tasks.TaskManager
	Manager *userManager.Manager
	Router  *userRoutes.Router
}

func NewStack(queries *store.Queries, s3 *s3.S3, tasks *tasks.TaskManager, auth *auth.AuthService) *Stack {
	store := userStore.NewStore(queries)
	manager := userManager.NewManager(store, s3, tasks, auth)
	router := userRoutes.NewRouter(manager)
	return &Stack{
		Store:   store,
		Tasks:   tasks,
		Manager: manager,
		Router:  router,
	}
}
