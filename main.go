package main

import (
	"dwelt/src/config"
	"dwelt/src/handler"
	"dwelt/src/model/dao"
	"dwelt/src/service/usrserv"
	"dwelt/src/ws/chat"
	"flag"
	"log/slog"
	"net/http"
)

var port = flag.String("port", ":8080", "port to listen on")

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	flag.Parse()
	config.InitCfg()
	db := dao.InitDB()

	hub := chat.NewHub()
	go hub.Run()

	userDao := dao.NewUserDao(db)
	userService := usrserv.NewUserService(hub, userDao)
	userController := handler.NewUserController(userService)
	userController.InitHandlers(hub)

	userService.StartHandlingMessages()

	server := &http.Server{
		Addr: *port,
	}
	slog.Info("Starting server on port " + *port)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server: ", err)
	}
}
