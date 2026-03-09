package api

import (
	"encoding/json"
	"keep/internal/keep"
	"keep/internal/store"
	"net/http"
)

// ListKeeps handles GET /keeps
func ListKeeps(w http.ResponseWriter, r *http.Request) {
	keeps, err := store.List()
	if err != nil {
		http.Error(w, "Failed to load keeps", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keeps)
}

// CreateKeep handles POST /keeps
func CreateKeep(w http.ResponseWriter, r *http.Request) {
	var req struct {
		KeepContent string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newKeep, err := store.Create(*keep.NewKeep(req.KeepContent))
	if err != nil {
		http.Error(w, "Failed to create keep", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newKeep)
}

// UpdateKeep handles PUT /keeps/{uuid}
func UpdateKeep(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	if uuidString == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}
	var req struct {
		KeepContent string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := store.Update(keep.KeepUUID(uuidString), req.KeepContent); err != nil {
		http.Error(w, "Keep not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteKeep handles DELETE /keeps/{uuid}
func DeleteKeep(w http.ResponseWriter, r *http.Request) {
	uuidString := r.PathValue("uuid")
	if uuidString == "" {
		http.Error(w, "Missing UUID", http.StatusBadRequest)
		return
	}
	if err := store.Delete(keep.KeepUUID(uuidString)); err != nil {
		http.Error(w, "Keep not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
