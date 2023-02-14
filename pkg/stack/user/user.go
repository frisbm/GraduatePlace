package user

import (
	userManager "github.com/MatthewFrisby/thesis-pieces/pkg/managers/user"
	userRoutes "github.com/MatthewFrisby/thesis-pieces/pkg/routes/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
	userStore "github.com/MatthewFrisby/thesis-pieces/pkg/store/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils/auth"
)

type Stack struct {
	Store   *userStore.Store
	Manager *userManager.Manager
	Router  *userRoutes.Router
}

func NewStack(queries *store.Queries, auth *auth.AuthService) *Stack {
	store := userStore.NewStore(queries)
	manager := userManager.NewManager(store, auth)
	router := userRoutes.NewRouter(manager)
	return &Stack{
		Store:   store,
		Manager: manager,
		Router:  router,
	}
}
