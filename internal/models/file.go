package models

type File struct {
	ID        int64  `json:"id" sql:"id"`
	Name      string `json:"name" sql:"name"`
	Location  string `json:"location" sql:"location"`
	CreatedAt string `json:"createdAt" sql:"created_at"`
}
