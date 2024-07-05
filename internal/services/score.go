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

type data struct {
	Score int64 `json:"score"`
}

func (s *Score) GetHandlerForQuestion(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	intID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	scoreNow, err := utils.GetScoreOfQuestion(s.db, int64(intID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, &data{
		Score: scoreNow,
	})
}

func (s *Score) GetHandlerForAnswer(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	intID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	scoreNow, err := utils.GetScoreOfAnswer(s.db, int64(intID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	utils.WriteJSON(w, &data{
		Score: scoreNow,
	})
}

func (s *Score) PostHandlerForQuestion(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	intID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return
	}

	var questionUserID int64

	query := `SELECT user_id FROM questions WHERE id=$1`
	result := s.db.QueryRow(query, id)

	err = result.Err()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = result.Scan(&questionUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if questionUserID == user.ID {
		http.Error(w, "Same User", http.StatusBadRequest)

		return
	}

	scoreNow, err := utils.GetScoreOfUserToQuestion(s.db, user.ID, int64(intID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	score := utils.ReadJSON[models.Score](w, req)

	if scoreNow >= 1 && strings.Compare(score.Operator, "plus") == 0 || scoreNow <= -1 && strings.Compare(score.Operator, "minus") == 0 {
		http.Error(w, "Voted Before", http.StatusBadRequest)

		return
	}

	if scoreNow == 1 && score.Operator == "minus" || scoreNow == -1 && score.Operator == "plus" {
		err := utils.ResetScoreOfUserToQustion(s.db, user.ID, int64(intID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		return
	}

	query = `INSERT INTO scores (operator,question_id,user_id) VALUES($1,$2,$3);`
	stmt, err := s.db.Prepare(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(score.Operator, id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
}

func (s *Score) PostHandlerForAnswer(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	intID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var answerUserID int64

	query := `SELECT user_id FROM answers WHERE id=$1`
	result := s.db.QueryRow(query, id)

	err = result.Err()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err = result.Scan(&answerUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if answerUserID == user.ID {
		http.Error(w, "Same User", http.StatusBadRequest)

		return
	}

	scoreNow, err := utils.GetScoreOfUserToAnswer(s.db, user.ID, int64(intID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	score := utils.ReadJSON[models.Score](w, req)

	if scoreNow >= 1 && strings.Compare(score.Operator, "plus") == 0 || scoreNow <= -1 && strings.Compare(score.Operator, "minus") == 0 {
		http.Error(w, "Voted Before", http.StatusBadRequest)

		return
	}

	if scoreNow == 1 && score.Operator == "minus" || scoreNow == -1 && score.Operator == "plus" {
		err := utils.ResetScoreOfUserToAnswer(s.db, user.ID, int64(intID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		return
	}

	query = `INSERT INTO scores (operator,answer_id,user_id) VALUES($1,$2,$3);`
	stmt, err := s.db.Prepare(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(score.Operator, id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

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
