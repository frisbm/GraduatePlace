package user

import (
	"github.com/MatthewFrisby/thesis-pieces/ent"
	userManager "github.com/MatthewFrisby/thesis-pieces/pkg/managers/user"
	userRoutes "github.com/MatthewFrisby/thesis-pieces/pkg/routes/user"
	userStore "github.com/MatthewFrisby/thesis-pieces/pkg/store/user"
)

type Stack struct {
	Store   *userStore.Store
	Manager *userManager.Manager
	Router  *userRoutes.Router
}

func NewStack(db *ent.Client) *Stack {
	store := userStore.NewStore(db)
	manager := userManager.NewManager(store)
	router := userRoutes.NewRouter(manager)
	return &Stack{
		Store:   store,
		Manager: manager,
		Router:  router,
	}
}
