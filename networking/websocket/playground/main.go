package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	handler := NewHandler(logger)

	fs := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /polling", handler.HandlePolling)
	mux.HandleFunc("GET /ws", handler.HandleWS)
	mux.HandleFunc("GET /broadcast", handler.HandleBroadcast)
	mux.HandleFunc("POST /connect", handler.Connect)
	mux.Handle("/", fs)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info("server started at 0.0.0.0:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error while starting server: ", err)
	}
}
