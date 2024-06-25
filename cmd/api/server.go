package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/middlewares"
	"forum-backend-go/internal/utils"
)

type server struct {
	router *mux.Router
}

func newServer(router *mux.Router) *server {
	router.Use(middlewares.Logging)

	return &server{
		router: router,
	}
}

func (s *server) serve() {
	srv := http.Server{
		Addr:              ":" + utils.GetEnv("port", "8000"),
		Handler:           s.router,
		ReadHeaderTimeout: time.Second}
	log.Printf("Starting Server on port%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
