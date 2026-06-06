package controllers

import (
	"MVC_Assignment/core"
	"MVC_Assignment/core/models"
	"encoding/json"
	"net/http"
)

var req models.RegisterUser

func Register(w http.ResponseWriter, r *http.Request) {
	json.NewDecoder(r.Body).Decode(&req)
	hash_pass, err := password.HashPass(req.Password)
	if err != nil {
		return "", err
	}
	if req.Username == "" {

	}
	if req.Password == "" {

	}

	user := core.User{
		UserName: req.Username,
		PassWord: hash_pass,
	}
}
