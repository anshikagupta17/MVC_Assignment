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

	err = repo.UpgradeTroops(village.ID, req.TroopID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
