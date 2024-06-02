package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	env    *enviroment
}

func NewServer(router *mux.Router) *server {
	return &server{
		env:    NewEnviroment(),
		router: router,
	}
}

func (s *server) serve() {
	srv := http.Server{
		Addr: fmt.Sprintf(":%s", s.env.getEnv("port", "8000")),
	}
	log.Printf("Starting Server on port%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
