package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
}

func NewServer(router *mux.Router) *server {
	return &server{
		router: router,
	}
}

func (s *server) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", GetEnv("port", "8000")),
		Handler: s.router,
	}
	log.Printf("Starting Server on port%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
