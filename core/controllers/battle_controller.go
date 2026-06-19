package controllers

import (
	"encoding/json"
	"net/http"

	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func FindOpponent(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

	opponent, err := repo.FindOpponent(village.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buildings, err := repo.GetOpponentBuildings(opponent.VillageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.OpponentResponse{
		VillageID:     opponent.VillageID,
		TownhallLevel: opponent.TownhallLevel,
		Trophies:      opponent.Trophies,
		Gold:          opponent.Gold,
		Elixir:        opponent.Elixir,
		Buildings:     buildings,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AttackVillage(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.AttackRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
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

	result, err := repo.AttackVillage(village.ID, req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
