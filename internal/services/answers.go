package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"forum-backend-go/internal/middlewares"
	"forum-backend-go/internal/models"
	"forum-backend-go/internal/utils"
)

func NewAnswerService(db *sql.DB, router *mux.Router) *Answer {
	return &Answer{
		db:     db,
		router: router,
	}
}

type Answer struct {
	db     *sql.DB
	router *mux.Router
}

func (r *Answer) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
	var excludedFields []string

	requestQuery := req.URL.Query()
	sortBy := requestQuery.Get("sort_by")
	sortDirection := requestQuery.Get("sort_direction")
	searchField := requestQuery.Get("search_field")
	searchFieldValue := requestQuery.Get("search_field_value")
	operator := requestQuery.Get("operator")
	limit := requestQuery.Get("limit")

	if _, err := strconv.Atoi(limit); err != nil && len(limit) != 0 {
		http.Error(w, "Limit is not number", http.StatusBadRequest)

		return
	}

	if len(sortDirection) != 0 && strings.Compare(sortDirection, "ASC") != 0 && strings.Compare(sortDirection, "DESC") != 0 {
		http.Error(w, "Sort direction in not ASC or DESC", http.StatusBadRequest)

		return
	}

	if len(sortDirection) == 0 {
		sortDirection = "ASC"
	}

	answers, err := GetMany[models.Answer](false, "answers", r.db, limit, sortBy, sortDirection, searchField, searchFieldValue, operator, excludedFields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, answers)
}

func (r *Answer) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	answer, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, answer)
}

func (r *Answer) PostHandler(w http.ResponseWriter, req *http.Request) {
	answer := utils.ReadJSON[models.Answer](w, req)
	if len(answer.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if len(answer.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	user := utils.GetUserFromRequest(req, w)
	answer.UserID = user.ID
	answer.UserName = user.Name
	id, err := r.Create(false, answer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (r *Answer) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	answer := utils.ReadJSON[models.Answer](w, req)

	if answer.Title != "" && len(answer.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if answer.Description != "" && len(answer.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.EditByID(false, int64(ID), answer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *Answer) DeleteHandler(w http.ResponseWriter, req *http.Request) {
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
func (r *Answer) Create(isTest bool, answer models.Answer) (int64, error) {
	var id int64
	id, err := Create[models.Answer](isTest, "answers", answer, r.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (r *Answer) DeleteByID(isTest bool, id int64) error {
	return Delete[models.Answer](isTest, "answers", r.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (r *Answer) EditByID(isTest bool, id int64, answer models.Answer) error {
	return Edit(isTest, "answers", r.db, "id", strconv.Itoa(int(id)), answer)
}

// GetByID implements Service.
func (r *Answer) GetByID(isTest bool, id int64) (models.Answer, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	answer, err := Get[models.Answer](isTest, "answers", r.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *answer, err
	}

	return *answer, nil
}

// RegisterRoutes implements Service.
func (r *Answer) RegisterRoutes() {
	router := r.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	AnswersRouter := APIV1Router.PathPrefix("/answers/").Subrouter()
	AnswersRouter.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	AnswersRouter.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	AnswersRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	AnswersRouter.HandleFunc("/{id}", r.PatchHandler).Methods("PATCH")
	AnswersRouter.HandleFunc("/{id}", r.DeleteHandler).Methods("DELETE")
}
