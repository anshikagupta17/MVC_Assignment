package models

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}
