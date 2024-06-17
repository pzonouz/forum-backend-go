package services

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/models"
	"forum-backend-go/internal/utils"
)

func NewRoleService(db *sql.DB, router *mux.Router) *Role {
	return &Role{
		db:     db,
		router: router,
	}
}

type Role struct {
	db     *sql.DB
	router *mux.Router
}

func (r *Role) GetHandlerForPlural(w http.ResponseWriter, _ *http.Request) {
	var excludedFields []string
	excludedFields = append(excludedFields, "id")
	roles, err := GetAll[models.Role](false, "roles", r.db, "20", "", "", excludedFields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, roles)
}

func (r *Role) GetHandler(w http.ResponseWriter, req *http.Request) {
	panic("unimplemented")
}

func (r *Role) PostHandler(w http.ResponseWriter, req *http.Request) {
	panic("unimplemented")
}

func (r *Role) PatchHandler(w http.ResponseWriter, req *http.Request) {
	panic("unimplemented")
}

func (r *Role) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	panic("unimplemented")
}

// Create implements Service.
func (r *Role) Create(isTest bool, role models.Role) (int64, error) {
	panic("unimplemented")
}

// DeleteByID implements Service.
func (r *Role) DeleteByID(isTest bool, id int64) error {
	panic("unimplemented")
}

// EditByID implements Service.
func (r *Role) EditByID(isTest bool, id int64, role models.Role) error {
	panic("unimplemented")
}

// GetAll implements Service.
func (r *Role) GetAll() ([]*models.Role, error) {
	panic("unimplemented")
}

// GetByID implements Service.
func (r *Role) GetByID(isTest bool, id int64) (models.Role, error) {
	panic("unimplemented")
}

// RegisterRoutes implements Service.
func (r *Role) RegisterRoutes() {
	router := r.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	RolesRouter := APIV1Router.PathPrefix("/roles/").Subrouter()
	RolesRouter.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	RolesRouter.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	RolesRouter.HandleFunc("/", r.PostHandler).Methods("POST")
	RolesRouter.HandleFunc("/{id}", r.PatchHandler).Methods("PATCH")
	RolesRouter.HandleFunc("/{id}", r.DeleteHandler).Methods("DELETE")
}
