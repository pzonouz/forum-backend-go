package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"forum-backend-go/ineternal/models"
	"forum-backend-go/ineternal/utils"
)

type Service[T any] interface {
	RegisterRoutes()
	GetHandler(w http.ResponseWriter, r *http.Request)
	GetHandlerForPlurar(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
	PatchHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	GetAll() ([]T, error)
	GetByID(id int) (T, error)
	Create(user T) error
	EditByID(user T) error
	DeleteByID(id int) error
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
	return errors.New("")
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
	// params := mux.Vars(r)
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
	user := utils.ReadJSON[models.User](w, r)
	err := u.Create(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

// registerRoutes implements Service.
func (u *UserService) RegisterRoutes() {
	router := u.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	APIV1Router.HandleFunc("", u.GetHandlerForPlurar)
	UsersRouter := APIV1Router.PathPrefix("/users/").Subrouter()
	UsersRouter.HandleFunc("/", u.GetHandlerForPlurar).Methods("GET")
	UsersRouter.HandleFunc("/{id}", u.GetHandler).Methods("GET")
	UsersRouter.HandleFunc("/", u.PostHandler).Methods("POST")
	UsersRouter.HandleFunc("/{id}", u.PatchHandler).Methods("PATCH")
	UsersRouter.HandleFunc("/{id}", u.DeleteHandler).Methods("DELETE")
}
