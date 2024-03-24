package handler

import (
	"dwelt/src/auth"
	"dwelt/src/config"
	"dwelt/src/dto"
	"dwelt/src/service/usrserv"
	"dwelt/src/utils"
	"dwelt/src/ws/chat"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const defaultLimit = 100

type UserController struct {
	userService *usrserv.UserService
}

func NewUserController(userService *usrserv.UserService) *UserController {
	return &UserController{userService}
}

func (uc UserController) InitHandlers(hub *chat.Hub) {
	router := mux.NewRouter()

	authenticatedRouter := router.PathPrefix("/").Subrouter()
	authenticatedRouter.Use(handlerAuthMiddleware)
	authenticatedRouter.HandleFunc("/hello", uc.handlerHelloWorld).Methods(http.MethodGet)
	authenticatedRouter.HandleFunc("/ws", uc.createHandlerWs(hub)).Methods(http.MethodGet)

	usersRouter := authenticatedRouter.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/search", uc.handlerSearchUsers).Methods(http.MethodGet)

	chatsRouter := authenticatedRouter.PathPrefix("/chats").Subrouter()
	chatsRouter.HandleFunc("/direct/{directToUid}", uc.handlerFindDirectChat).Methods(http.MethodGet)

	noAuthRouter := router.PathPrefix("/").Subrouter()
	noAuthRouter.HandleFunc("/register", uc.handlerRegister).Methods(http.MethodPost)
	noAuthRouter.HandleFunc("/login", uc.handlerLogin).Methods(http.MethodGet)
	noAuthRouter.HandleFunc("/info", uc.handleApplicationInfoDashboard).Methods(http.MethodGet)

	http.Handle("/", router)
}

func (uc UserController) handlerLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, valid, err := uc.userService.ValidateUser(username, password)
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
	utils.WriteJson(w, dto.UserInfo{UserId: userId})
}

func (uc UserController) handlerRegister(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, duplicate, err := uc.userService.RegisterUser(username, password)
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
	utils.WriteJson(w, dto.UserInfo{UserId: userId})
}

func (uc UserController) handlerSearchUsers(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	users, err := uc.userService.SearchUsers(prefix, defaultLimit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, users)
}

func (uc UserController) handlerFindDirectChat(w http.ResponseWriter, r *http.Request) {
	requesterUid := retrieveUserId(r)
	userId, err := strconv.ParseInt(mux.Vars(r)["directToUid"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, badUsers, err := uc.userService.FindDirectChat(requesterUid, userId) // todo: show old messages
	if badUsers {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, chatId)
}

func (uc UserController) handlerHelloWorld(w http.ResponseWriter, r *http.Request) { // todo remove
	w.WriteHeader(http.StatusOK)
	userId, _ := r.Context().Value("userId").(int64)
	utils.Must(w.Write([]byte("hello, " + strconv.FormatInt(userId, 10))))
}

func (uc UserController) handleApplicationInfoDashboard(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	utils.Must(w.Write([]byte("Workflow run number: " + strconv.Itoa(config.DweltCfg.WorkflowRunNumber))))
}

func (uc UserController) createHandlerWs(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, retrieveUserId(r), w, r)
	}
}
