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

func (a *Answer) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
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

	answers, err := GetMany[models.Answer](false, "answers", a.db, limit, sortBy, sortDirection, searchField, searchFieldValue, operator, excludedFields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, answers)
}

func (a *Answer) GetHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	answer, err := a.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, answer)
}

func (a *Answer) PostHandler(w http.ResponseWriter, req *http.Request) {
	questionId, err := strconv.Atoi(mux.Vars(req)["question_id"])
	if err != nil {
		http.Error(w, "questions_id is not true", http.StatusBadRequest)

		return
	}

	answer := utils.ReadJSON[models.Answer](w, req)

	if len(answer.Description) < 6 {
		http.Error(w, "At least 6 character for Description", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	answer.QuestionID = int64(questionId)
	answer.UserID = user.ID
	answer.UserName = user.Name
	id, err := a.Create(false, answer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (a *Answer) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	answer := utils.ReadJSON[models.Answer](w, req)

	if answer.Description != "" && len(answer.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = a.EditByID(false, int64(ID), answer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (a *Answer) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	ID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = a.DeleteByID(false, int64(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// Create implements Service.
func (a *Answer) Create(isTest bool, answer models.Answer) (int64, error) {
	var id int64
	id, err := Create[models.Answer](isTest, "answers", answer, a.db)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// DeleteByID implements Service.
func (a *Answer) DeleteByID(isTest bool, id int64) error {
	return Delete[models.Answer](isTest, "answers", a.db, "id", strconv.Itoa(int(id)))
}

// EditByID implements Service.
func (a *Answer) EditByID(isTest bool, id int64, answer models.Answer) error {
	return Edit(isTest, "answers", a.db, "id", strconv.Itoa(int(id)), answer)
}

// GetByID implements Service.
func (a *Answer) GetByID(isTest bool, id int64) (models.Answer, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	answer, err := Get[models.Answer](isTest, "answers", a.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *answer, err
	}

	return *answer, nil
}

// RegisterRoutes implements Service.
func (a *Answer) RegisterRoutes() {
	router := a.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	AnswersRouter := APIV1Router.PathPrefix("/answers/").Subrouter()
	AnswersRouter.HandleFunc("/", a.GetHandlerForPlural).Methods("GET")
	AnswersRouter.HandleFunc("/{id}", a.GetHandler).Methods("GET")
	AnswersRouter.HandleFunc("/{question_id}", middlewares.LoginGuard(a.PostHandler)).Methods("POST")
	AnswersRouter.HandleFunc("/{id}", a.PatchHandler).Methods("PATCH")
	AnswersRouter.HandleFunc("/{id}", a.DeleteHandler).Methods("DELETE")
}
