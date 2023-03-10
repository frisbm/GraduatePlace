package utils

import (
	"context"

	"github.com/pkg/errors"

	"github.com/frisbm/graduateplace/pkg/store"
)

func GetUserFromContext(ctx context.Context) (*store.User, error) {
	user := ctx.Value("user").(*store.User)
	if user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}

func Ptr[T any](t T) *T {
	return &t
}
