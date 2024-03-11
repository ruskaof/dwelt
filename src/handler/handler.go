package handler

import (
	"dwelt/src/auth"
	"dwelt/src/config"
	"dwelt/src/utils"
	"dwelt/src/ws/chat"
	"net/http"
	"strconv"
)

type UserInfo struct {
	Username string `json:"username"`
	Id       int64  `json:"id"`
}

type UserHandlerFunc func(w http.ResponseWriter, r *http.Request, userInfo UserInfo)

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
		next(w, r, UserInfo{Username: username, Id: -1})
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
	utils.WriteJson(w, UserInfo{Username: username, Id: -1}) // todo: send correct id
}

func handlerHelloWorld(w http.ResponseWriter, _ *http.Request, userInfo UserInfo) { // todo remove
	w.WriteHeader(http.StatusOK)
	utils.Must(w.Write([]byte("hello, " + userInfo.Username)))
}

func handleApplicationInfoDashboard(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	utils.Must(w.Write([]byte("Workflow run number: " + strconv.Itoa(config.DweltCfg.WorkflowRunNumber))))
}

func createHandlerWs(hub *chat.Hub) UserHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, userInfo UserInfo) {
		chat.ServeWs(hub, userInfo.Username, w, r)
	}
}
