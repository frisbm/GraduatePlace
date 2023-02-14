package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/MatthewFrisby/thesis-pieces/ent"
	"github.com/MatthewFrisby/thesis-pieces/pkg/models/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils/auth"
)

type Store interface {
	CreateUser(ctx context.Context, registerUser user.RegisterUser) error
	GetUserForLogin(ctx context.Context, email string) ([]*ent.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	GetUserByContext(ctx context.Context) (*ent.User, error)
	GetUsers(ctx context.Context) ([]*ent.User, error)
}

type Auth interface {
	ParseRefreshToken(refreshToken string) (*uuid.UUID, error)
	GenerateTokens(uuid string) (*user.AuthTokens, error)
}

type Manager struct {
	store Store
	auth  Auth
}

func NewManager(store Store, auth Auth) *Manager {
	return &Manager{
		store: store,
		auth:  auth,
	}
}

func (m *Manager) RegisterUser(ctx context.Context, registerUser user.RegisterUser) error {
	err := auth.ValidatePassword(registerUser.Password)
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(registerUser.Password)
	if err != nil {
		// Don't return the actual error here for privacy reasons in-case password present in error
		return errors.New("error hashing password")
	}
	registerUser.Password = hashedPassword
	return m.store.CreateUser(ctx, registerUser)
}

func (m *Manager) LoginUser(ctx context.Context, loginUser user.LoginUser) (*user.AuthTokens, error) {
	entUsers, err := m.store.GetUserForLogin(ctx, loginUser.Email)
	if err != nil {
		return nil, err
	}

	if len(entUsers) == 0 {
		return nil, errors.New("invalid login information")
	}

	if len(entUsers) > 1 || entUsers[0] == nil {
		return nil, errors.New("unexpected error occurred")
	}

	entUser := entUsers[0]

	err = auth.ValidatePasswordCorrect(entUser.Password, loginUser.Password)

	if err != nil {
		return nil, errors.New("invalid login information")
	}

	return m.auth.GenerateTokens(entUser.UUID.String())
}

func (m *Manager) RefreshUser(ctx context.Context, refreshUser user.RefreshUser) (*user.AuthTokens, error) {
	uuid, err := m.auth.ParseRefreshToken(refreshUser.RefreshToken)
	if err != nil {
		return nil, err
	}

	entUser, err := m.store.GetUserByUUID(ctx, *uuid)
	if err != nil {
		return nil, err
	}

	return m.auth.GenerateTokens(entUser.UUID.String())
}

func (m *Manager) GetUser(ctx context.Context) (*user.GetUser, error) {
	entUser, err := m.store.GetUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	return &user.GetUser{
		UUID:      entUser.UUID.String(),
		FirstName: entUser.FirstName,
		LastName:  entUser.LastName,
		Email:     entUser.Email,
		Username:  entUser.Username,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
	}, nil
}

func (m *Manager) GetUsers(ctx context.Context) ([]*user.GetUser, error) {
	entUsers, err := m.store.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	var users []*user.GetUser
	for _, entUser := range entUsers {
		newUser := &user.GetUser{
			UUID:      entUser.UUID.String(),
			FirstName: entUser.FirstName,
			LastName:  entUser.LastName,
			Email:     entUser.Email,
			Username:  entUser.Username,
			CreatedAt: entUser.CreatedAt,
			UpdatedAt: entUser.UpdatedAt,
		}
		users = append(users, newUser)
	}
	return users, nil
}
