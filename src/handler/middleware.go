package handler

import (
	"context"
	"dwelt/src/auth"
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userId, valid, err := auth.ValidateToken(tokenString[7:])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		newReq := r.WithContext(context.WithValue(r.Context(), "userId", userId))
		next(w, newReq)
	}
}
