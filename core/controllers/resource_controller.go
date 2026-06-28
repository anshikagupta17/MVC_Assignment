package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func ResourceCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user_id, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(user_id)
	if err != nil {
		http.Error(w, "Village not found", http.StatusBadRequest)
		return
	}
	collected_resources, err := repo.CollectVillageResources(village.ID)
	if err != nil {
		http.Error(w, "Collected resources not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(collected_resources)

}
