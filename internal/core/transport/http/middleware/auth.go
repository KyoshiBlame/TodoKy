package core_http_middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIdKey contextKey = "user_id"

type TokenValidater interface {
	ValidateToken(ctx context.Context, token string) (userID int64, err error)
}

func Auth(validator TokenValidater) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Redirect(w, r, "/auth.html", http.StatusTemporaryRedirect)
				return
			}

			userID, err := validator.ValidateToken(r.Context(), cookie.Value)
			if err != nil {
				http.Redirect(w, r, "/auth.html", http.StatusTemporaryRedirect)
			}

			ctx := context.WithValue(r.Context(), UserIdKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIdKey).(int64)
	return userID, ok
}
