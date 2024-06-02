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

		// if route is must present token, but the token is empty. Problem: always return 200 if token not present, but is optional = false
		if !isOptional && authHeader == "" {
			ResponseErrorInvalidAccessToken(w, "invalid access token")
			return
		}

		token, err := getToken(authHeader)
		// if header Authorization: Bearer preset with empty token
		if err != nil {
			ResponseErrorInvalidAccessToken(w, "invalid access token")
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
		}

		return
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

// corsMiddleware is a middleware that sets CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
