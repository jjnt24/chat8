package auth

import (
	"context"
	"net/http"
)

type contextKey string

const userKey contextKey = "user"

// Middleware untuk memvalidasi session cookie
func AuthMiddleware(store *SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := GetSessionData(r, store)
			if err != nil {
				// http.Error(w, "unauthorized", http.StatusUnauthorized)
				http.Redirect(w, r, "/login", http.StatusSeeOther)

				return
			}

			ctx := context.WithValue(r.Context(), userKey, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSessionData(r *http.Request, store *SessionStore) (*SessionData, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	data, err := store.Get(context.Background(), cookie.Value)
	if err != nil || data == nil {
		return nil, err
	}

	return data, err
}

func GetUserFromContext(ctx context.Context) *SessionData {
	if val, ok := ctx.Value(userKey).(*SessionData); ok {
		return val
	}
	return nil
}
