package middleware

import (
	"context"
	"e-cart/pkg/api"
	"e-cart/pkg/jwt"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserIDKey   contextKey = "userid"
	UsernameKey contextKey = "username"
	IsAdminKey  contextKey = "isadmin"
)

// middleware for users routes
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.Fail(w, http.StatusUnauthorized, 401, "Authorization header is missing", "")
			return
		}

		// Extract the token part (removing 'Bearer ' prefix) //have 2 part we wont take bearer part
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			api.Fail(w, http.StatusUnauthorized, 401, "Invalid token format", "")
			return
		}

		// Validate the token
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				api.Fail(w, http.StatusUnauthorized, 401, "Token expired, please log in again", "")
				return
			}

			api.Fail(w, http.StatusUnauthorized, 401, "Invalid token", err.Error())
			return
		}

		// Store userid and username in context,so it can be user on other layers from the ctx
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)

		// updating and Passing the control to the next handler func we have
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// middleware for admin-only routes
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
		if !ok || !isAdmin {
			api.Fail(w, http.StatusForbidden, 403, "Admin access required", "")
			return
		}
		next.ServeHTTP(w, r)
	})
}
