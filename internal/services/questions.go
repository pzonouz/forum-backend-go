package services

import (
	"context"
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

type QuestionModel struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	ViewCount   int64  `json:"view"`
	UserName    string `json:"userName"`
	UserID      int64  `json:"userId"`
	ScoreCount  int64  `json:"scoreCount"`
	AnswerCount int64  `json:"answerCount"`
	Solved      bool   `json:"solved"`
}

func (r *Question) GetHandlerForPlural(w http.ResponseWriter, req *http.Request) {
	orderBy := req.URL.Query().Get("order_by")
	orderDirection := req.URL.Query().Get("order_direction")
	searchFiled := req.URL.Query().Get("search_field")
	searchFiledValue := req.URL.Query().Get("search_field_value")
	query := `SELECT qs.id,qs.title,qs.description,qs.created_at,COUNT(DISTINCT CASE WHEN ans.solved THEN 1 ELSE NULL END) as solved,us.nickname,us.id,COUNT(DISTINCT vw.id) as view_count,COUNT(DISTINCT ans.id) as answer_count,SUM(DISTINCT CASE WHEN sc.operator='plus' THEN 1 ELSE -1 END) as score FROM questions as qs LEFT JOIN "views" as vw ON vw.question_id=qs.id LEFT JOIN users as us ON qs.user_id=us.id LEFT JOIN answers as ans ON ans.question_id=qs.id LEFT JOIN scores as sc ON sc.question_id=qs.id GROUP BY qs.id,us.nickname,us.id`

	if strings.Compare(searchFiled, "") != 0 {
		query = `SELECT * FROM (` + query
		query = query + `) WHERE ` + searchFiled + searchFiledValue
	}

	if strings.Compare(orderBy, "") != 0 {
		query = query + ` ORDER BY ` + orderBy
		if strings.Compare(orderDirection, "") == 0 {
			query = query + ` ASC`
		}

		if strings.Compare(orderDirection, "DESC") == 0 {
			query = query + ` DESC`
		}
	}

	rows, err := r.db.Query(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	defer rows.Close()

	var questions []QuestionModel

	for rows.Next() {
		question := QuestionModel{}
		_ = rows.Scan(&question.ID, &question.Title, &question.Description, &question.CreatedAt, &question.Solved, &question.UserName, &question.UserID, &question.ViewCount, &question.AnswerCount, &question.ScoreCount)
		questions = append(questions, question)
	}

	utils.WriteJSON(w, questions)
}

func (r *Question) GetHandlerForPluralOfQuestions(w http.ResponseWriter, req *http.Request) {
	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	query := `SELECT * FROM questions WHERE user_id=$1`
	stmt, err := r.db.Prepare(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	questions, err := QueryRowsToStruct[models.Question](stmt, nil, user.ID)
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

func (r *Question) GetViewUpHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, _ := utils.GetUserFromRequest(req, w)
	query := `INSERT INTO "views" (question_id,user_id) VALUES ($1,$2)`
	_, _ = r.db.Exec(query, id, user.ID)

	return
}

func (r *Question) PostHandler(w http.ResponseWriter, req *http.Request) {
	type Question struct {
		ID          int64         `json:"id" sql:"id"`
		Title       string        `json:"title" sql:"title"`
		Description string        `json:"description" sql:"description"`
		CreatedAt   string        `json:"createdAt" sql:"created_at"`
		UserName    string        `json:"userName" sql:"user_name"`
		UserID      int64         `json:"userId" sql:"user_id"`
		Files       []models.File `json:"files"`
	}

	question := utils.ReadJSON[Question](w, req)
	files := question.Files
	if len(question.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if len(question.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	question.UserID = user.ID
	question.UserName = user.NickName

	question2 := &models.Question{
		Title:       question.Title,
		Description: question.Description,
		UserID:      question.UserID,
		UserName:    question.UserName,
	}

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer tx.Rollback()

	var id int64
	query := `INSERT INTO questions(title,description,user_name,user_id) VALUES($1,$2,$3,$4) RETURNING "id"`
	err = tx.QueryRow(query, question2.Title, question2.Description, question2.UserName, question2.UserID).Scan(&id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query = `UPDATE files SET question_id=$1 WHERE id=$2`
	for _, value := range files {
		result, err := tx.Exec(query, id, value.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rows, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if rows == 0 {
			http.Error(w, "0 rows affected", http.StatusBadRequest)

			return
		}
	}
	tx.Commit()
	type data struct {
		ID int64 `json:"id"`
	}

	utils.WriteJSON(w, &data{ID: id})
}

func (r *Question) PatchHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}
	questionPartial := utils.ReadJSON[models.Question](w, req)
	question, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, "not Found", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)

	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if user.ID != question.UserID && user.Role != "admin" {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if questionPartial.Title != "" && len(questionPartial.Title) < 11 {
		http.Error(w, "At least 10 character for title", http.StatusBadRequest)

		return
	}

	if questionPartial.Description != "" && len(questionPartial.Description) < 21 {
		http.Error(w, "At least 20 character for Description", http.StatusBadRequest)

		return
	}

	err = r.EditByID(false, int64(id), question)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (r *Question) DeleteHandler(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "ID is not integer", http.StatusBadRequest)

		return
	}

	question, err := r.GetByID(false, int64(id))

	if err != nil {
		http.Error(w, "not Found", http.StatusBadRequest)

		return
	}

	user, err := utils.GetUserFromRequest(req, w)

	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	if user.ID != question.UserID && user.Role != "admin" {
		http.Error(w, "", http.StatusUnauthorized)

		return
	}

	err = r.DeleteByID(false, int64(id))

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
	QuestionsRouter.HandleFunc("/current_user/", middlewares.LoginGuard(r.GetHandlerForPluralOfQuestions)).Methods("GET")
	QuestionsRouter.HandleFunc("/{id}/", r.GetHandler).Methods("GET")
	QuestionsRouter.HandleFunc("/{id}/view_up", middlewares.LoginGuard(r.GetViewUpHandler)).Methods("GET")
	QuestionsRouter.HandleFunc("/", middlewares.LoginGuard(r.PostHandler)).Methods("POST")
	QuestionsRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.PatchHandler)).Methods("PATCH")
	QuestionsRouter.HandleFunc("/{id}", middlewares.LoginGuard(r.DeleteHandler)).Methods("DELETE")
}
