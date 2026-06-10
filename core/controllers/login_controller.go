package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/auth"
	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUser

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

	repo := repositories.UserRepository{
		DB: db.Conn,
	}

	user, err := repo.GetUser(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = models.CheckPass(req.Password, user.PassWord)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := models.Login_Response{
		Message: "Login was successful",
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	username := r.Context().Value("username")

	response := map[string]interface{}{
		"user_id":  userID,
		"username": username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetVillage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(userID)
	if err != nil {
		http.Error(w, "Village not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(village)
}
