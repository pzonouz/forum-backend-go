package models

type View struct {
	ID         int64 `json:"id" sql:"id"`
	QuestionID int64 `json:"questionId" sql:"question_id"`
	UserID     int64 `json:"userId" sql:"user_id"`
}
