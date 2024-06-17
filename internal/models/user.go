package models

type User struct {
	ID          int64  `json:"id" sql:"id"`
	Email       string `json:"email" sql:"email"`
	Password    string `json:"password" sql:"password"`
	Name        string `json:"name" sql:"name"`
	Address     string `json:"address" sql:"address"`
	PhoneNumber string `json:"phoneNumber" sql:"phone_number"`
	CreatedAt   string `json:"createdAt" sql:"created_at"`
	RolerID     int64  `json:"rolerId" sql:"roler_id"`
}
