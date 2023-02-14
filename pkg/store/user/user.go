package user

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/user"
)

type Store struct {
	db store.Querier
}

func NewStore(db store.Querier) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context, registerUser user.RegisterUser) error {
	_, err := s.db.CreateUser(ctx, store.CreateUserParams{
		Username:  registerUser.Username,
		Email:     registerUser.Email,
		Password:  registerUser.Password,
		FirstName: registerUser.FirstName,
		LastName:  registerUser.LastName,
	})
	return err
}

func (s *Store) GetUserFromEmail(ctx context.Context, email string) (*store.User, error) {
	return s.db.GetUserFromEmail(ctx, email)
}

func (s *Store) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*store.User, error) {
	return s.db.GetUserFromUUID(ctx, uuid)
}

func (s *Store) GetUserByContext(ctx context.Context) (*store.User, error) {
	user := ctx.Value("user").(*store.User)
	if user == nil {
		return nil, errors.New("no user in context")
	}
	return s.db.GetUserFromUUID(ctx, user.Uuid)
}

func (s *Store) GetUsers(ctx context.Context) ([]*store.User, error) {
	return s.db.GetUsers(ctx)
}
