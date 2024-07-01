package models

type Answer struct {
	ID          int64  `json:"id" sql:"id"`
	Description string `json:"description" sql:"description"`
	CreatedAt   string `json:"createdAt" sql:"created_at"`
	UserName    string `json:"userName" sql:"user_name"`
	UserID      int64  `json:"userId" sql:"user_id"`
	QuestionID  int64  `json:"questionId" sql:"question_id"`
}
