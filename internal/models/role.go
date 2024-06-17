package models

type Role struct {
	ID        int64  `json:"id" sql:"id"`
	Name      string `json:"name" sql:"name"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
}
