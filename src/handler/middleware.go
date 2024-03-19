package handler

import (
	"context"
	"dwelt/src/auth"
	"log/slog"
	"net/http"
)

func handlerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 8 {
			slog.Debug("Invalid token length", "token", tokenString)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, valid, err := auth.ValidateToken(tokenString[7:])
		if err != nil {
			slog.Error(err.Error(), "method", "ValidateToken")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !valid {
			slog.Debug("Validation failed", "token", tokenString)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		slog.Debug("User authenticated", "userId", userId)
		r = r.WithContext(context.WithValue(r.Context(), "userId", userId))
		next.ServeHTTP(w, r)
	})
}
