package controllers

import (
	"encoding/json"
	"net/http"
)

func ResourceCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	village, repo, ok := GetVillage(w, r)
	if !ok {
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
