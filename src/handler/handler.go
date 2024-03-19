package handler

import (
	"dwelt/src/auth"
	"dwelt/src/config"
	"dwelt/src/service/usrserv"
	"dwelt/src/utils"
	"dwelt/src/ws/chat"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func InitHandlers(hub *chat.Hub) {
	router := mux.NewRouter()

	authenticatedRouter := router.PathPrefix("/").Subrouter()
	authenticatedRouter.Use(handlerAuthMiddleware)
	authenticatedRouter.HandleFunc("/hello", handlerHelloWorld).Methods(http.MethodGet)
	authenticatedRouter.HandleFunc("/ws", createHandlerWs(hub)).Methods(http.MethodGet)
	authenticatedRouter.HandleFunc("/info", handleApplicationInfoDashboard).Methods(http.MethodGet)

	noAuthRouter := router.PathPrefix("/").Subrouter()
	noAuthRouter.HandleFunc("/register", handlerRegister).Methods(http.MethodPost)
	noAuthRouter.HandleFunc("/login", handlerLogin).Methods(http.MethodGet)

	http.Handle("/", router)
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, valid, err := usrserv.ValidateUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := auth.GenerateToken(userId)
	w.Header().Set("Authorization", "Bearer "+token)
	utils.WriteJson(w, userInfo{UserId: userId})
}

func handlerRegister(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, duplicate, err := usrserv.RegisterUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if duplicate {
		w.WriteHeader(http.StatusConflict)
		return
	}

	token := auth.GenerateToken(userId)
	w.Header().Set("Authorization", "Bearer "+token)
	utils.WriteJson(w, userInfo{UserId: userId})
}

func handlerHelloWorld(w http.ResponseWriter, r *http.Request) { // todo remove
	w.WriteHeader(http.StatusOK)
	userId, _ := r.Context().Value("userId").(int64)
	utils.Must(w.Write([]byte("hello, " + strconv.FormatInt(userId, 10))))
}

func handleApplicationInfoDashboard(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	utils.Must(w.Write([]byte("Workflow run number: " + strconv.Itoa(config.DweltCfg.WorkflowRunNumber))))
}

func createHandlerWs(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, retrieveUserId(r), w, r)
	}
}
