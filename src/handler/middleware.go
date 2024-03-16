package handler

import (
	"context"
	"dwelt/src/auth"
	"log/slog"
	"net/http"
)

func handlerGETMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func handlerPOSTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func handlerAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 8 {
			slog.Debug("Invalid token length", "token", tokenString)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, valid, err := auth.ValidateToken(tokenString[7:])
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !valid {
			slog.Debug("Validation failed", "token", tokenString)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		slog.Debug("User authenticated", "userId", userId)
		newReq := r.WithContext(context.WithValue(r.Context(), "userId", userId))
		next(w, newReq)
	}
}
