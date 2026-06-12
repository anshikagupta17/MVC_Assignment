package controllers

import (
	"net/http"

	"encoding/json"

	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

type MoveBuildingRequest struct {
	BuildingInstanceId int64 `json:"building_instance_id"`
	X                  int   `json:"x"`
	Y                  int   `json:"y"`
}

func GetVillage(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(userId)

	if err != nil {
		http.Error(w, "Village not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(village)

}

func GetVillageBuildings(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}
	village, err := repo.GetVillage(userId)
	if err != nil {
		http.Error(w, "Village not found", http.StatusNotFound)
		return
	}
	building, err := repo.VillageBuildings(village.ID)
	if err != nil {
		http.Error(w, "Buildings not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(building)
}

func MoveBuilding(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	repo := repositories.VillageRepository{
		DB: db.Conn,
	}

	village, err := repo.GetVillage(userId)
	if err != nil {
		http.Error(w, "Village not found", http.StatusNotFound)
		return
	}

	var req MoveBuildingRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.X < 0 || req.X > 49 || req.Y < 0 || req.Y > 49 {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}

	check, err := repo.CanPlaceBuilding(village.ID, req.BuildingInstanceId, req.X, req.Y)
	if err != nil {
		http.Error(w, "Error checking placement", http.StatusInternalServerError)
		return
	}

	if !check {
		http.Error(w, "Position blocked by another building", http.StatusBadRequest)
		return
	}

	err = repo.MoveBuilding(village.ID, req.BuildingInstanceId, req.X, req.Y)
	if err != nil {
		http.Error(w, "Failed to move building", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Building moved successfully",
	})
}
