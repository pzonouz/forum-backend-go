package services

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"forum-backend-go/ineternal/models"
)

type Service[T any] interface {
	RegisterRoutes()
	GetHandler(w http.ResponseWriter, r *http.Request)
	GetHandlerForPlurar(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
	PatchHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	GetAll() ([]T, error)
	GetByID(int) (T, error)
	Create(T) error
	EditByID(T) error
	DeleteByID(int) error
}

func NewUserService(db *sql.DB, router *mux.Router) *UserService {
	return &UserService{
		db: db, router: router,
	}
}

type UserService struct {
	db     *sql.DB
	router *mux.Router
}

// Create implements Service.
func (u *UserService) Create(user models.User) error {
	panic("unimplemented")
}

// DeleteByID implements Service.
func (u *UserService) DeleteByID(int) error {
	panic("unimplemented")
}

// DeleteHandler implements Service.
func (u *UserService) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// EditByID implements Service.
func (u *UserService) EditByID(user models.User) error {
	panic("unimplemented")
}

// GetAll implements Service.
func (u *UserService) GetAll() ([]models.User, error) {
	panic("unimplemented")
}

// GetByID implements Service.
func (u *UserService) GetByID(int) (models.User, error) {
	panic("unimplemented")
}

// GetHandler implements Service.
func (u *UserService) GetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	w.Write([]byte(id))
}

// GetHandlerForPlurar implements Service.
func (u *UserService) GetHandlerForPlurar(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Users"))
}

// PatchHandler implements Service.
func (u *UserService) PatchHandler(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostHandler implements Service.
func (u *UserService) PostHandler(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	API_V1_Router := router.PathPrefix("/api/v1/").Subrouter()
	API_V1_Router.HandleFunc("", u.GetHandlerForPlurar)
	UsersRouter := API_V1_Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/", u.GetHandlerForPlurar)
	UsersRouter.HandleFunc("/{id}", u.GetHandler)
}
