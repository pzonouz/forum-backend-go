package models

type Question struct {
	ID          int64  `json:"id" sql:"id"`
	Title       string `json:"title" sql:"title" validate:"required"`
	Description string `json:"description" sql:"description" validate:"required"`
	CreateAt    string `json:"createAt" sql:"created_at"`
	UserID      int64  `json:"userId" sql:"user_id"`
}
