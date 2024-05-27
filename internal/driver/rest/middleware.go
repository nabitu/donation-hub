package rest

import (
	"context"
	"net/http"
	"strings"
)

// authTokenMiddleware middleware for checking Bearer token
func authTokenMiddleware(next http.HandlerFunc, c *Config, isOptional bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ResponseErrorInvalidAccessToken(w, "Invalid access token")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		ctx := context.WithValue(r.Context(), "token", token)
		payload, err := c.AuthTokenService.ValidateToken(token)

		if isOptional {
			// isOptional = true, means the token does not have to be valid (used in the list of projects)
			if err == nil {
				ctx = context.WithValue(ctx, "auth_id", payload.UserID)
				ctx = context.WithValue(ctx, "auth_username", payload.Username)
				ctx = context.WithValue(ctx, "auth_email", payload.Email)
			}
		} else {
			// isOptionalnya = false, then the token must be correct
			if err != nil {
				ResponseErrorInvalidAccessToken(w, "Invalid access token")
				return
			}

			ctx = context.WithValue(ctx, "auth_id", payload.UserID)
			ctx = context.WithValue(ctx, "auth_username", payload.Username)
			ctx = context.WithValue(ctx, "auth_email", payload.Email)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
