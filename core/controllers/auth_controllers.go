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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	user := core.User{
		UserName: req.Username,
		PassWord: hash_pass,
	}
}
