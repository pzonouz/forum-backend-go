package models

type User struct {
	ID               int64  `json:"id" sql:"id"`
	Email            string `json:"email" sql:"email"`
	Password         string `json:"password" sql:"password"`
	NickName         string `json:"nickName" sql:"nickname"`
	Address          string `json:"address" sql:"address"`
	PhoneNumber      string `json:"phoneNumber" sql:"phone_number"`
	Role             string `json:"role" sql:"role"`
	CreatedAt        string `json:"createdAt" sql:"created_at"`
	Token            string `json:"token" sql:"token"`
	IsForgetPassword bool   `json:"isForgetPassword" sql:"is_forget_password"`
}
