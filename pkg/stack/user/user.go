package user

import (
	userManager "github.com/frisbm/graduateplace/pkg/managers/user"
	userRoutes "github.com/frisbm/graduateplace/pkg/routes/user"
	"github.com/frisbm/graduateplace/pkg/services/auth"
	"github.com/frisbm/graduateplace/pkg/services/s3"
	"github.com/frisbm/graduateplace/pkg/store"
	userStore "github.com/frisbm/graduateplace/pkg/store/user"
	"github.com/frisbm/graduateplace/pkg/tasks"
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
