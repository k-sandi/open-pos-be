package variants

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// StatusRequest represents the payload for updating active status.
type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type ModifierHandler struct {
	svc ModifierService
}

func NewModifierHandler(svc ModifierService) *ModifierHandler {
	return &ModifierHandler{svc: svc}
}

// Routes returns a Chi router for modifier endpoints.
func (h *ModifierHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/group/{groupID}", h.ListModifiers)
	r.Post("/", h.CreateModifier)
	r.Put("/{id}", h.UpdateModifier)
	r.Patch("/{id}/status", h.UpdateStatus)
	r.Delete("/{id}", h.DeleteModifier)
	return r
}

// @Summary List Modifiers by Group
// @Tags modifiers
// @Produce json
// @Param groupID path string true "Modifier Group ID"
// @Security ApiKeyAuth
// @Router /modifiers/group/{groupID} [get]
func (h *ModifierHandler) ListModifiers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	modifiers, err := h.svc.ListModifiers(r.Context(), groupID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if modifiers == nil {
		modifiers = []*Modifier{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modifiers)
}

// @Summary Create Modifier
// @Tags modifiers
// @Accept json
// @Produce json
// @Param request body ModifierDTO true "Modifier Data"
// @Security ApiKeyAuth
// @Router /modifiers [post]
func (h *ModifierHandler) CreateModifier(w http.ResponseWriter, r *http.Request) {
	var req ModifierDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.svc.CreateModifier(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

// @Summary Update Modifier
// @Tags modifiers
// @Accept json
// @Produce json
// @Param id path string true "Modifier ID"
// @Param request body ModifierDTO true "Modifier Data"
// @Security ApiKeyAuth
// @Router /modifiers/{id} [put]
func (h *ModifierHandler) UpdateModifier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ModifierDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateModifier(r.Context(), id, req); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Modifier Status
// @Tags modifiers
// @Accept json
// @Produce json
// @Param id path string true "Modifier ID"
// @Param request body StatusRequest true "Status Update Request"
// @Security ApiKeyAuth
// @Router /modifiers/{id}/status [patch]
func (h *ModifierHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateStatus(r.Context(), id, req.IsActive); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete Modifier
// @Tags modifiers
// @Security ApiKeyAuth
// @Router /modifiers/{id} [delete]
func (h *ModifierHandler) DeleteModifier(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteModifier(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
