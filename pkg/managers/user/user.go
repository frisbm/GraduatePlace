package user

import (
	"context"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/MatthewFrisby/thesis-pieces/pkg/models/user"
	"github.com/MatthewFrisby/thesis-pieces/pkg/utils/auth"
)

type Store interface {
	CreateUser(ctx context.Context, registerUser user.RegisterUser) error
	GetUserFromEmail(ctx context.Context, email string) (*store.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (*store.User, error)
	GetUserByContext(ctx context.Context) (*store.User, error)
	GetUsers(ctx context.Context) ([]*store.User, error)
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
	result, err := m.store.GetUserFromEmail(ctx, loginUser.Email)
	if err != nil {
		return nil, err
	}

	err = auth.ValidatePasswordCorrect(result.Password, loginUser.Password)

	if err != nil {
		return nil, errors.New("invalid login information")
	}

	return m.auth.GenerateTokens(result.Uuid.String())
}

func (m *Manager) RefreshUser(ctx context.Context, refreshUser user.RefreshUser) (*user.AuthTokens, error) {
	uuid, err := m.auth.ParseRefreshToken(refreshUser.RefreshToken)
	if err != nil {
		return nil, err
	}

	result, err := m.store.GetUserByUUID(ctx, *uuid)
	if err != nil {
		return nil, err
	}

	return m.auth.GenerateTokens(result.Uuid.String())
}

func (m *Manager) GetUser(ctx context.Context) (*user.GetUser, error) {
	result, err := m.store.GetUserByContext(ctx)
	if err != nil {
		return nil, err
	}
	return &user.GetUser{
		UUID:      result.Uuid.String(),
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (m *Manager) GetUsers(ctx context.Context) ([]*user.GetUser, error) {
	results, err := m.store.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	var users []*user.GetUser
	for _, result := range results {
		newUser := &user.GetUser{
			UUID:      result.Uuid.String(),
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Email:     result.Email,
			Username:  result.Username,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
		users = append(users, newUser)
	}
	return users, nil
}
