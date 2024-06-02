package main

import (
	"database/sql"

	"github.com/gorilla/mux"

	"forum-backend-go/ineternal/models"
	"forum-backend-go/ineternal/services"
)

type Router interface {
	registerRoutes(db *sql.DB)
	getMux() *mux.Router
}

type router struct {
	mux *mux.Router
	db  *sql.DB
}

// getMux implements Router.
func (r *router) getMux() *mux.Router {
	return r.mux
}

// registerRoutes implements Router.
func (r *router) registerRoutes(db *sql.DB) {
	var userService services.Service[models.User] = services.NewUserService(db, r.getMux())

	userService.RegisterRoutes()
}

func newRouter(db *sql.DB) *router {
	return &router{
		mux: mux.NewRouter(),
		db:  db,
	}
}
