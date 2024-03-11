package main

import (
	"dwelt/src/config"
	"dwelt/src/handler"
	"dwelt/src/model/dao"
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
	dao.InitDB()

	hub := chat.NewHub()
	go hub.Run()

	handler.InitHandlers(hub)

	server := &http.Server{
		Addr: *port,
	}
	slog.Info("Starting server on port " + *port)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server: ", err)
	}
}
