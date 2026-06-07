package models

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register_Response struct {
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login_Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
