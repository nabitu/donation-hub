package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func RecoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// print error
				fmt.Println("panic error", err)
				ResponseErrorInternalServerError(w, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// authTokenMiddleware middleware for checking Bearer token
func authTokenMiddleware(next http.HandlerFunc, c *Config, isOptional bool, roles []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		ctx := context.WithValue(r.Context(), "withAuth", false)

		// if route authorization is optional
		if isOptional && authHeader == "" {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		token, err := getToken(authHeader)
		if err != nil {
			return
		}

		payload, err := c.AuthTokenService.ValidateToken(token)

		// make sure access token is valid if present
		if authHeader != "" && err != nil {
			ResponseErrorInvalidAccessToken(w, "invalid access token")
			return
		}

		ctx = context.WithValue(ctx, "withAuth", true)
		ctx = context.WithValue(ctx, "payload", payload)
		ctx = context.WithValue(ctx, "auth_id", payload.UserID)
		ctx = context.WithValue(ctx, "auth_username", payload.Username)
		ctx = context.WithValue(ctx, "auth_email", payload.Email)
		ctx = context.WithValue(ctx, "auth_roles", payload.Role)

		if len(roles) > 0 && validRole(roles, payload.Role) {
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			ResponseErrorForbiddenAccess(w, "user doesn't have enough authorization")
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getToken(tokenString string) (string, error) {
	token := strings.Split(tokenString, " ")
	if len(token) != 2 {
		return "", fmt.Errorf("invalid token format")
	}
	if strings.ToLower(token[0]) != "bearer" {
		return "", fmt.Errorf("invalid token format")
	}
	return token[1], nil
}

func validRole(roles []string, userRole []string) bool {
	for _, role := range roles {
		for _, uRole := range userRole {
			if role == uRole {
				return true
			}
		}
	}
	return false
}
