package main

import (
	"github.com/gorilla/mux"
)

type Router interface {
}

type router struct {
	mux *mux.Router
}

func NewRouter() *router {
	return &router{
		mux: mux.NewRouter(),
	}

}
