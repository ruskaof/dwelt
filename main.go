package main

import (
	"dwelt/handler"
	"dwelt/ws/chat"
	"flag"
	"log/slog"
	"net/http"
)

var port = flag.String("port", ":8080", "port to listen on")

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	flag.Parse()

	hub := chat.NewHub()
	go hub.Run()

	handler.InitHandlers(hub)

	server := &http.Server{
		Addr: *port,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server: ", err)
	}
}
