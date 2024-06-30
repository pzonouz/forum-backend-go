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

func NewScoreService(db *sql.DB, router *mux.Router) *Score {
	return &Score{
		db:     db,
		router: router,
	}
}

type Score struct {
	db     *sql.DB
	router *mux.Router
}

func (s *Score) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
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

	scores, err := GetMany[models.Score](false, "scores", s.db, limit, sortBy, sortDirection, searchField, searchFieldValue, operator, excludedFields)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, scores)
}

func (s *Score) GetHandlerForQuestion(w http.ResponseWriter, req *http.Request) {
	// params := mux.Vars(req)
	// id, err := strconv.Atoi(params["id"])
	//
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// score, err := r.GetByID(false, int64(id))
	//
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	//
	// 	return
	// }
	//
	// utils.WriteJSON(w, score)
}

func (s *Score) GetHandlerForAnswer(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	score, err := s.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, score)
}

func (s *Score) PostHandlerForQuestion(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	score := utils.ReadJSON[models.Score](w, req)
	user := utils.GetUserFromRequest(req, w)
	query := `INSERT INTO scores (operator,question_id,user_id) VALUES($1,$2,$3);`
	stmt, err := s.db.Prepare(query)

	if err != nil {
		http.Error(w, err.Error(), 400)

		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(score.Operator, id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)

		return
	}
}

func (s *Score) PostHandlerForAnswer(w http.ResponseWriter, req *http.Request) {

	id := mux.Vars(req)["id"]
	score := utils.ReadJSON[models.Score](w, req)
	user := utils.GetUserFromRequest(req, w)
	query := `INSERT INTO scores (operator,answer_id,user_id) VALUES($1,$2,$3);`
	stmt, err := s.db.Prepare(query)

	if err != nil {
		http.Error(w, err.Error(), 400)

		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(score.Operator, id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)

		return
	}
}

// GetByID implements Service.
func (s *Score) GetByID(isTest bool, id int64) (models.Score, error) {
	var excludedFields []string
	// excludedFields = append(excludedFields, "id")
	score, err := Get[models.Score](isTest, "scores", s.db, "id", strconv.Itoa(int(id)), excludedFields)

	if err != nil {
		return *score, err
	}

	return *score, nil
}

// RegisterRoutes implements Service.
func (s *Score) RegisterRoutes() {
	router := s.router
	APIV1Router := router.PathPrefix("/api/v1/").Subrouter()
	ScoresRouter := APIV1Router.PathPrefix("/scores/").Subrouter()
	// ScoresRouter.HandleFunc("/", r.GetHandlerForPlural).Methods("GET")
	ScoresRouter.HandleFunc("/questions/{id}", s.GetHandlerForQuestion).Methods("GET")
	ScoresRouter.HandleFunc("/answers/{id}", s.GetHandlerForAnswer).Methods("GET")
	ScoresRouter.HandleFunc("/questions/{id}", middlewares.LoginGuard(s.PostHandlerForQuestion)).Methods("POST")
	ScoresRouter.HandleFunc("/answers/{id}", middlewares.LoginGuard(s.PostHandlerForAnswer)).Methods("POST")
}
