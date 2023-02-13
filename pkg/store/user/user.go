package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/MatthewFrisby/thesis-pieces/ent"
	user2 "github.com/MatthewFrisby/thesis-pieces/ent/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/models/user"
)

type Store struct {
	db *ent.Client
}

func NewStore(db *ent.Client) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(ctx context.Context, registerUser user.RegisterUser) error {
	_, err := s.db.User.Create().
		SetEmail(registerUser.Email).
		SetPassword(registerUser.Password).
		SetUsername(registerUser.Username).
		SetFirstName(registerUser.FirstName).
		SetLastName(registerUser.LastName).
		SetIsAdmin(false).
		Save(ctx)
	return err
}

func (s *Store) GetUserForLogin(ctx context.Context, email, password string) ([]*ent.User, error) {
	query := s.db.User.Query().Where(
		user2.Email(email),
		user2.Password(password),
	)
	return query.All(ctx)
}

func (s *Store) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*ent.User, error) {
	return s.db.User.Query().
		Where(user2.UUID(uuid)).
		Only(ctx)
}

func (s *Store) GetUsers(ctx context.Context) ([]*ent.User, error) {
	return s.db.User.Query().All(ctx)
}
