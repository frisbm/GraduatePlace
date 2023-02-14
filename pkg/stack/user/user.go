package user

import (
	"github.com/MatthewFrisby/thesis-pieces/ent"
	userManager "github.com/MatthewFrisby/thesis-pieces/pkg/managers/user"
	userRoutes "github.com/MatthewFrisby/thesis-pieces/pkg/routes/user"
	userStore "github.com/MatthewFrisby/thesis-pieces/pkg/store/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils/auth"
)

type Stack struct {
	Store   *userStore.Store
	Manager *userManager.Manager
	Router  *userRoutes.Router
}

func NewStack(db *ent.Client, auth *auth.AuthService) *Stack {
	store := userStore.NewStore(db)
	manager := userManager.NewManager(store, auth)
	router := userRoutes.NewRouter(manager)
	return &Stack{
		Store:   store,
		Manager: manager,
		Router:  router,
	}
}
