package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func UpgradeTroops(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UpgradeTroopRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(userID)
	if err != nil {
		http.Error(w, "village not found", http.StatusNotFound)
		return
	}

	err = repo.UpgradeTroopsTX(village.ID, req.TroopID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func TrainTroops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.TrainTroopsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(userID)
	if err != nil {
		http.Error(w, "Village not found", http.StatusNotFound)
		return
	}

	err = repo.TrainTroops(village.ID, req.TroopID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Troops trained successfully",
	})
}
