package variants

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)


type ModifierGroupHandler struct {
	svc ModifierGroupService
}

func NewModifierGroupHandler(svc ModifierGroupService) *ModifierGroupHandler {
	return &ModifierGroupHandler{svc: svc}
}

// Routes returns a Chi router mounted with modifier group routes.
func (h *ModifierGroupHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/products/{product_id}/modifier-groups", h.ListByProduct)
	r.Post("/modifier-groups", h.CreateGroup)
	r.Put("/modifier-groups/{id}", h.UpdateGroup)
	r.Patch("/modifier-groups/{id}/status", h.UpdateStatus)
	r.Delete("/modifier-groups/{id}", h.DeleteGroup)
	return r
}

// @Summary List Modifier Groups by Product
// @Tags modifier-groups
// @Produce json
// @Param product_id path string true "Product ID"
// @Security ApiKeyAuth
// @Router /products/{product_id}/modifier-groups [get]
func (h *ModifierGroupHandler) ListByProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")
	
	groups, err := h.svc.ListByProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if groups == nil {
		groups = []*ModifierGroup{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// @Summary Create Modifier Group
// @Tags modifier-groups
// @Accept json
// @Produce json
// @Param request body ModifierGroupDTO true "Modifier Group Data"
// @Security ApiKeyAuth
// @Router /modifier-groups [post]
func (h *ModifierGroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req ModifierGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mg, err := h.svc.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mg)
}

// @Summary Update Modifier Group
// @Tags modifier-groups
// @Accept json
// @Produce json
// @Param id path string true "Modifier Group ID"
// @Param request body ModifierGroupDTO true "Update Data"
// @Security ApiKeyAuth
// @Router /modifier-groups/{id} [put]
func (h *ModifierGroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ModifierGroupDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.Update(r.Context(), id, req); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Modifier Group Status
// @Tags modifier-groups
// @Accept json
// @Produce json
// @Param id path string true "Modifier Group ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /modifier-groups/{id}/status [patch]
func (h *ModifierGroupHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
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

// @Summary Delete Modifier Group
// @Tags modifier-groups
// @Security ApiKeyAuth
// @Router /modifier-groups/{id} [delete]
func (h *ModifierGroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
