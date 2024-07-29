package models

type FileRequest struct {
	ID        int64  `json:"id" sql:"id"`
	Request   string `json:"request" sql:"request"`
	UserID    int64  `json:"userId" sql:"user_id"`
	Done      bool   `json:"done" sql:"done"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
}
