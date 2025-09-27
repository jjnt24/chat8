package auth

import (
	"context"
	"net/http"
	"strings"
)

func (s *JWTService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := s.ParseAccess(token)
		if err != nil {
			if err == ErrTokenExpired {
				http.Error(w, "Token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			return
		}

		// store claims in context for handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
