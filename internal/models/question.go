package models

type Question struct {
	ID          int64  `json:"id" sql:"id"`
	Title       string `json:"title" sql:"title"`
	Description string `json:"description" sql:"description"`
	CreatedAt   string `json:"createdAt" sql:"created_at"`
	UserName    string `json:"userName" sql:"user_name"`
	UserID      int64  `json:"userId" sql:"user_id"`
	Solved      bool   `json:"solved" sql:"solved"`
}
