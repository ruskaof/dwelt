package handler

import (
	"dwelt/assets"
	"dwelt/src/auth"
	"dwelt/src/config"
	"dwelt/src/dto"
	"dwelt/src/service/usrserv"
	"dwelt/src/utils"
	"dwelt/src/ws/chat"
	"github.com/flowchartsman/swaggerui"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	router.Use(handlerMetricsMiddleware)

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
	http.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(assets.SwaggerYaml)))
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", router)
}

// @Summary		Authenticate user
// @Description	Get a JWT token for the user using basic auth
// @Tags			Auth
// @Accept			json
// @Param			Authorization	header	string	true	"Basic auth"
// @Produce		json
// @Success		200	{object}	dto.UserInfo
// @Failure		401
// @Failure		500
// @Header			200	{string}	Authorization	"Bearer <token>"
// @Router			/login [get]
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

// @Summary		Register user
// @Description	Register a new user using basic auth
// @Tags			Auth
// @Accept			json
// @Param			Authorization	header	string	true	"Basic auth"
// @Produce		json
// @Success		200	{object}	dto.UserInfo
// @Failure		401
// @Failure		409
// @Failure		500
// @Header			200	{string}	Authorization
// @Router			/register [post]
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

// @Summary		Search users
// @Description	Search for users by prefix
// @Tags			User
// @Accept			json
// @Param			prefix			query	string	true	"Prefix to search for"
// @Param			Authorization	header	string	true	"Bearer <token>"
// @Produce		json
// @Success		200	{object}	[]dto.UserInfo
// @Failure		500
// @Router			/users/search [get]
func (uc UserController) handlerSearchUsers(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	users, err := uc.userService.SearchUsers(prefix, defaultLimit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.WriteJson(w, users)
}

// @Summary		Open direct chat
// @Description	Creates a chat if it does not exist and returns the chat id and previous messages
// @Tags			Chats
// @Accept			json
// @Param			directToUid		path	int64	true	"User id to open chat with"
// @Param			Authorization	header	string	true	"Bearer <token>"
// @Produce		json
// @Success		200	{object}	dto.OpenDirectChatResponse
// @Failure		500
// @Router			/chats/direct/{directToUid} [get]
func (uc UserController) handlerFindDirectChat(w http.ResponseWriter, r *http.Request) {
	requesterUid := retrieveUserId(r)
	userId, err := strconv.ParseInt(mux.Vars(r)["directToUid"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := uc.userService.OpenDirectChat(requesterUid, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, res)
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

// @Summary		Starts websocket connection
// @Description	Connects to the websocket
// @Description	<br>Messages from the server are in the following format:
// @Description	```json
// @Description	{
// @Description	"chatId": "integer",
// @Description	"userId": "integer",
// @Description	"username": "string",
// @Description	"message": "string",
// @Description	"createdAt": "date"
// @Description	}
// @Description	```
// @Description	<br>Messages to the server are in the following format:
// @Description	```json
// @Description	{
// @Description	"chatId": "integer",
// @Description	"message": "string"
// @Description	}
// @Description	```
// @Tags			Ws
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500
// @Router			/ws [get]
func (uc UserController) createHandlerWs(hub *chat.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, retrieveUserId(r), w, r)
	}
}
