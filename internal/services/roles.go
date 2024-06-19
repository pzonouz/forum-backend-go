package services

import (
	"database/sql"
	"net/http"
	"strconv"

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
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	role, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, role)
}

func (r *Role) PostHandler(w http.ResponseWriter, req *http.Request) {
	role := utils.ReadJSON[models.Role](w, req)
	id, err := r.Create(false, role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (r *Role) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	role := utils.ReadJSON[models.Role](w, req)
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.EditByID(false, int64(ID), role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *Role) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.DeleteByID(false, int64(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// Create implements Service.
func (r *Role) Create(isTest bool, role models.Role) (int64, error) {
	var id int64
	id, err := Create[models.Role](isTest, "roles", role, r.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (r *Role) DeleteByID(isTest bool, id int64) error {
	return Delete[models.Role](isTest, "roles", r.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (r *Role) EditByID(isTest bool, id int64, role models.Role) error {
	return Edit(isTest, "roles", r.db, "id", strconv.Itoa(int(id)), role)
}

// GetByID implements Service.
func (r *Role) GetByID(isTest bool, id int64) (models.Role, error) {
	var excludedFields []string
	excludedFields = append(excludedFields, "id")
	role, err := Get[models.Role](isTest, "roles", r.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *role, err
	}

	return *role, nil
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
