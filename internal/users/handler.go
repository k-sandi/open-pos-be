package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// StatusRequest represents the payload for updating a user's active status.
type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

// Routes returns a Chi router mounted with user routes.
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Patch("/{id}/status", UpdateStatusHandler)
	return r
}

// UpdateStatusHandler handles PATCH /api/v1/users/:id/status
// @Summary Update User Status
// @Description Updates the active status of a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /users/{id}/status [patch]
func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: DB update logic goes here. Mocking success for now.
	w.WriteHeader(http.StatusOK)
}
