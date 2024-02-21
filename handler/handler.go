package handler

import (
	"dwelt/auth"
	"net/http"
)

func InitHandlers() {
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/hello", handlerAuthMiddleware(handlerHelloWorld))
}

func handlerAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")[7:]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		valid, err := auth.ValidateToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func handlerLogin(w http.ResponseWriter, _ *http.Request) {
	token := auth.GenerateToken()
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func handlerHelloWorld(w http.ResponseWriter, _ *http.Request) { // todo remove
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, world"))
}
