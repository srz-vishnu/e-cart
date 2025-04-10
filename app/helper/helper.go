package app

import (
	"context"
	"e-cart/pkg/middleware"
	"errors"
)

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

// Helper function to get username from context
func GetUsernameFromContext(ctx context.Context) (string, error) {
	username, ok := ctx.Value(middleware.UsernameKey).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}
	return username, nil
}
