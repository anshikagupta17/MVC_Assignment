package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/auth"
	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUser

	err := json.NewDecoder(r.Body).Decode(&req)
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

	hash_pass, err := auth.HashPass(req.Password)

	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	repo := repositories.UserRepository{
		DB: db.Conn,
	}

	user_id, err := repo.RegisterTX(req.Username, hash_pass)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.Register_Response{
		Message: "Registered successfully",
		UserId:  user_id,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
