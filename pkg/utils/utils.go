package utils

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
)

func GetUserFromContext(ctx context.Context) (*store.User, error) {
	user := ctx.Value("user").(*store.User)
	if user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
