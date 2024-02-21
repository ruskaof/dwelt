package main

import (
	"dwelt/handler"
	"flag"
	"log/slog"
	"net/http"
)

var port = flag.String("port", ":8080", "port to listen on")

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	flag.Parse()

	handler.InitHandlers()

	//hub := chat.NewHub()
	//go hub.Run()
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	chat.ServeWs(hub, w, r)
	//})
	server := &http.Server{
		Addr: *port,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("error starting server: ", err)
	}
}
