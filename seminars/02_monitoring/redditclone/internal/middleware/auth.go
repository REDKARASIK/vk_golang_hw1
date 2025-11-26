package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey") // In production, use env/config

type contextKey string

const userContextKey = contextKey("user")

type UserClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Auth is a middleware to protect routes that require authentication.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"message": "missing token"}`, http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"message": "invalid token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"message": "invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		userMap, ok := claims["user"].(map[string]interface{})
		if !ok {
			panic("invalid user data in token")
		}

		id, ok := userMap["id"].(string)
		if !ok {
			http.Error(w, `{"message": "invalid user id in token"}`, http.StatusUnauthorized)
			return
		}

		username, ok := userMap["username"].(string)
		if !ok {
			http.Error(w, `{"message": "invalid username in token"}`, http.StatusUnauthorized)
			return
		}

		user := UserClaims{
			ID:       id,
			Username: username,
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (UserClaims, bool) {
	user, ok := ctx.Value(userContextKey).(UserClaims)
	return user, ok
}
