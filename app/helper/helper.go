package helper

import (
	"context"
	"e-cart/pkg/middleware"
	"errors"
)

type ContextHelper interface {
	GetUserID(ctx context.Context) (int64, error)
	GetUsername(ctx context.Context) (string, error)
}

type contextHelperImpl struct{}

func NewContextHelper() ContextHelper {
	return &contextHelperImpl{}
}

func (h *contextHelperImpl) GetUserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

func (h *contextHelperImpl) GetUsername(ctx context.Context) (string, error) {
	username, ok := ctx.Value(middleware.UsernameKey).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}
