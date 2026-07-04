package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func UpgradeTroops(w http.ResponseWriter, r *http.Request) {

	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	var req models.UpgradeTroopRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
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

	village, repo, ok := GetVillage(w, r)
	if !ok {
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

	err = repo.TrainTroopsTX(village.ID, req.TroopID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Troops trained successfully",
	})
}
