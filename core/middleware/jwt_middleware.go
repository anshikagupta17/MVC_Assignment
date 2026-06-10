package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/anshikagupta17/MVC_Assignment/core/auth"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth_header := r.Header.Get("Authorization")

		if auth_header == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth_header, "Bearer ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := auth_header[len("Bearer "):]
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "username", claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
