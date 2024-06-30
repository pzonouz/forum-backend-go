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

func NewQuestionService(db *sql.DB, router *mux.Router) *Question {
	return &Question{
		db:     db,
		router: router,
	}
}

type Question struct {
	db     *sql.DB
	router *mux.Router
}

func (r *Question) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
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

	questions, err := GetMany[models.Question](false, "questions", r.db, limit, sortBy, sortDirection, searchField, searchFieldValue, operator, excludedFields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, questions)
}

func (r *Question) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	question, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, question)
}

func (r *Question) PostHandler(w http.ResponseWriter, req *http.Request) {
	question := utils.ReadJSON[models.Question](w, req)
	if len(question.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if len(question.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	user := utils.GetUserFromRequest(req, w)
	question.UserID = user.ID
	question.UserName = user.Name
	id, err := r.Create(false, question)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (r *Question) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	question := utils.ReadJSON[models.Question](w, req)

	if question.Title != "" && len(question.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if question.Description != "" && len(question.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = r.EditByID(false, int64(ID), question)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *Question) DeleteHandler(w http.ResponseWriter, req *http.Request) {
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
func (r *Question) Create(isTest bool, question models.Question) (int64, error) {
	var id int64
	id, err := Create[models.Question](isTest, "questions", question, r.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (r *Question) DeleteByID(isTest bool, id int64) error {
	return Delete[models.Question](isTest, "questions", r.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (r *Question) EditByID(isTest bool, id int64, question models.Question) error {
	return Edit(isTest, "questions", r.db, "id", strconv.Itoa(int(id)), question)
}

// GetByID implements Service.
func (r *Question) GetByID(isTest bool, id int64) (models.Question, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	question, err := Get[models.Question](isTest, "questions", r.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *question, err
	}

	return *question, nil
}

func (r *Question) ScorePostHandler(w http.ResponseWriter, req *http.Request) {
}

// RegisterRoutes implements Service.
func (r *Question) RegisterRoutes() {
	router := r.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	QuestionsRouter := APIV1Router.PathPrefix("/questions/").Subrouter()
	QuestionsRouter.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	QuestionsRouter.HandleFunc("/{id}", r.GetHandler).Methods("GET")
	QuestionsRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	QuestionsRouter.HandleFunc("/{id}", r.PatchHandler).Methods("PATCH")
	QuestionsRouter.HandleFunc("/{id}", r.DeleteHandler).Methods("DELETE")
}
