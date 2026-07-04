package controllers

import (
	"net/http"

	db "github.com/anshikagupta17/MVC_Assignment/core/database"
	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/anshikagupta17/MVC_Assignment/core/repositories"
)

func getUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	return userID, true
}

func GetVillage(w http.ResponseWriter, r *http.Request) (models.Village, repositories.VillageRepository, bool) {
	userID, ok := getUserID(w, r)
	if !ok {
		return models.Village{}, repositories.VillageRepository{}, false
	}

	repo := repositories.VillageRepository{DB: db.Conn}

	village, err := repo.GetVillage(userID)
	if err != nil {
		http.Error(w, "village not found", http.StatusNotFound)
		return models.Village{}, repositories.VillageRepository{}, false
	}

	return village, repo, true
}

func GetVillageState(w http.ResponseWriter, r *http.Request) (models.VillageResponse, repositories.VillageRepository, bool) {
	userID, ok := getUserID(w, r)
	if !ok {
		return models.VillageResponse{}, repositories.VillageRepository{}, false
	}

	repo := repositories.VillageRepository{DB: db.Conn}

	village, err := repo.VillageState(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.VillageResponse{}, repositories.VillageRepository{}, false
	}

	return village, repo, true
}
