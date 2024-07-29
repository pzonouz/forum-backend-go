package models

type File struct {
	ID         int64  `json:"id" sql:"id"`
	Title      string `json:"title" sql:"title"`
	FileName   string `json:"filename" sql:"filename"`
	CreatedAt  string `json:"createdAt" sql:"created_at"`
	UserID     int64  `json:"userId" sql:"user_id"`
	QuestionID int64  `json:"questionId" sql:"question_id"`
	AnswerID   int64  `json:"answerId" sql:"answer_id"`
	FileType   string `json:"fileType" sql:"filetype"`
}
