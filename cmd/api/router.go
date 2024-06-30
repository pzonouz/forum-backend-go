package main

import (
	"database/sql"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/models"
	"forum-backend-go/internal/services"
)

type Router interface {
	registerRoutes()
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
func (r *router) registerRoutes() {
	var userService services.Service[models.User] = services.NewUserService(r.db, r.getMux())

	var questionService services.Service[models.Question] = services.NewQuestionService(r.db, r.getMux())

	var answerService services.Service[models.Answer] = services.NewAnswerService(r.db, r.getMux())

	userService.RegisterRoutes()
	questionService.RegisterRoutes()
	answerService.RegisterRoutes()
}

func newRouter(db *sql.DB) *router {
	return &router{
		mux: mux.NewRouter(),
		db:  db,
	}
}
