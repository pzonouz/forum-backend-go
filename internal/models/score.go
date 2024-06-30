package models

type Score struct {
	ID         int64  `json:"id" sql:"id"`
	Operator   string `json:"operator" sql:"operator"`
	CreatedAt  string `json:"createdAt" sql:"created_at"`
	QuestionID int64  `json:"questionId" sql:"question_id"`
	AnswerID   int64  `json:"answerId" sql:"answer_id"`
	UserID     int64  `json:"userId" sql:"user_id"`
}
