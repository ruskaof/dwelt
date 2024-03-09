package handler

import (
	"dwelt/auth"
	"dwelt/ws/chat"
	"flag"
	"net/http"
	"strconv"
)

var workflowRunNumber = flag.Int("wflrn", 0, "workflow run number")

type UserHandlerFunc func(w http.ResponseWriter, r *http.Request, username string)

func InitHandlers(hub *chat.Hub) {
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/hello", handlerAuthMiddleware(handlerHelloWorld))
	http.HandleFunc("/ws", handlerAuthMiddleware(createHandlerWs(hub)))
	http.HandleFunc("/info", handleApplicationInfoDashboard)
}

func handlerAuthMiddleware(next UserHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 8 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, valid, err := auth.ValidateToken(tokenString[7:])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, username)
	}
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	// get username and password from request
	username, _, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// todo: validate username and password

	token := auth.GenerateToken(username)
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func handlerHelloWorld(w http.ResponseWriter, _ *http.Request, username string) { // todo remove
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello, " + username))
}

func handleApplicationInfoDashboard(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Workflow run number: " + strconv.Itoa(*workflowRunNumber)))
}

func createHandlerWs(hub *chat.Hub) UserHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, username string) {
		chat.ServeWs(hub, username, w, r)
	}
}
