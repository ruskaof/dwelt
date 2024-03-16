package handler

import (
	"context"
	"dwelt/src/auth"
	"fmt"
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
			slog.Debug(fmt.Sprintf("Invalid token length: %d", len(tokenString)))
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
			slog.Debug(fmt.Sprintf("Validation failed for token: %s", tokenString))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		slog.Debug(fmt.Sprintf("User %d authenticated", userId))
		newReq := r.WithContext(context.WithValue(r.Context(), "userId", userId))
		next(w, newReq)
	}
}
