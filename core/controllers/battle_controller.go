package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func FindOpponent(w http.ResponseWriter, r *http.Request) {

	village, repo, ok := GetVillage(w, r)
	if !ok {
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

	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	var req models.AttackRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
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
