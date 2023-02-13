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
	GetUserForLogin(ctx context.Context, email, password string) ([]*ent.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*ent.User, error)
	GetUsers(ctx context.Context) ([]*ent.User, error)
}

type Manager struct {
	store Store
}

func NewManager(store Store) *Manager {
	return &Manager{
		store: store,
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
	hashedPassword, err := auth.HashPassword(loginUser.Password)
	if err != nil {
		// Don't return the actual error here for privacy reasons in-case password present in error
		return nil, errors.New("error hashing password")
	}
	entUser, err := m.store.GetUserForLogin(ctx, loginUser.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	if len(entUser) == 0 {
		return nil, errors.New("invalid login information")
	}

	if len(entUser) > 1 || entUser[0] == nil {
		return nil, errors.New("unexpected error occurred")
	}

	return auth.GenerateTokens(entUser[0].UUID.String())
}

func (m *Manager) RefreshUser(ctx context.Context, refreshUser user.RefreshUser) (*user.AuthTokens, error) {
	uuid, err := auth.ParseRefreshToken(refreshUser.RefreshToken)
	if err != nil {
		return nil, err
	}

	entUser, err := m.store.GetUserByUUID(ctx, *uuid)
	if err != nil {
		return nil, err
	}

	return auth.GenerateTokens(entUser.UUID.String())
}

func (m *Manager) GetUser(ctx context.Context) (*user.GetUser, error) {
	return nil, nil
}

func (m *Manager) GetUsers(ctx context.Context) (*user.GetUsers, error) {
	return nil, nil
}
