package controllers

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func VillageState(w http.ResponseWriter, r *http.Request) {
	village, _, ok := GetVillageState(w, r)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(village)
}
func GetVillageBuildings(w http.ResponseWriter, r *http.Request) {
	village, repo, ok := GetVillage(w, r)
	if !ok {
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
	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	var req models.MoveBuildingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.X < 0 || req.X > 49 || req.Y < 0 || req.Y > 49 {
		http.Error(w, "Invalid position", http.StatusBadRequest)
		return
	}

	size_x, size_y, err := repo.GetBuildingSize(village.ID, req.BuildingInstanceId)
	if err != nil {
		http.Error(w, "Building not found", http.StatusNotFound)
		return
	}

	canPlace, err := repositories.CanPlace(context.Background(), repo.DB, village.ID, size_x, size_y, req.X, req.Y, req.BuildingInstanceId)
	if err != nil {
		http.Error(w, "Error checking placement", http.StatusInternalServerError)
		return
	}
	if !canPlace {
		http.Error(w, "Position blocked by another building", http.StatusBadRequest)
		return
	}

	err = repo.MoveBuilding(village.ID, req.BuildingInstanceId, req.X, req.Y)
	if err != nil {
		http.Error(w, "Failed to move building", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Building moved successfully",
	})
}

func AddBuilding(w http.ResponseWriter, r *http.Request) {
	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	var req models.AddBuildingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = repo.AddBuildingTX(village.ID, req.BuildingID, req.X, req.Y)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func GetVillageTroops(w http.ResponseWriter, r *http.Request) {

	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	troops, err := repo.GetVillageTroops(village.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(troops)
}

func UpgradeBuilding(w http.ResponseWriter, r *http.Request) {

	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	var req models.UpgradeBuildingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = repo.BuildingUpgradeTX(village.ID, req.BuildingInstanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetShopBuildings(w http.ResponseWriter, r *http.Request) {
	village, repo, ok := GetVillage(w, r)
	if !ok {
		return
	}

	buildings, err := repo.GetShopBuildings(village.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}
