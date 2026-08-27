package taxes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TaxResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Rate     float64 `json:"rate"`
	IsActive bool    `json:"is_active"`
}

type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type CreateTaxRequest struct {
	Name     string  `json:"name"`
	Rate     float64 `json:"rate"`
	IsActive bool    `json:"is_active"`
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Routes returns a Chi router mounted with tax routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListTaxes)
	r.Get("/{id}", h.GetTax)
	r.Post("/", h.CreateTax)
	r.Put("/{id}", h.UpdateTax)
	r.Patch("/{id}/status", h.UpdateStatusHandler)
	r.Delete("/{id}", h.DeleteTax)
	return r
}

// @Summary List Taxes
// @Tags taxes
// @Produce json
// @Security ApiKeyAuth
// @Router /taxes [get]
func (h *Handler) ListTaxes(w http.ResponseWriter, r *http.Request) {
	taxes, err := h.svc.ListTaxes(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var res []TaxResponse
	for _, t := range taxes {
		res = append(res, mapTaxToResponse(t))
	}
	if res == nil {
		res = []TaxResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @Summary Get Tax
// @Tags taxes
// @Produce json
// @Param id path string true "Tax ID"
// @Security ApiKeyAuth
// @Router /taxes/{id} [get]
func (h *Handler) GetTax(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tax, err := h.svc.GetTax(r.Context(), id)
	if err != nil {
		http.Error(w, "tax not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapTaxToResponse(tax))
}

// @Summary Create Tax
// @Tags taxes
// @Accept json
// @Produce json
// @Param request body CreateTaxRequest true "Tax Data"
// @Security ApiKeyAuth
// @Router /taxes [post]
func (h *Handler) CreateTax(w http.ResponseWriter, r *http.Request) {
	var req CreateTaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := CreateTaxDTO{
		Name:     req.Name,
		Rate:     req.Rate,
		IsActive: req.IsActive,
	}

	_, err := h.svc.CreateTax(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update Tax
// @Tags taxes
// @Accept json
// @Produce json
// @Param id path string true "Tax ID"
// @Param request body CreateTaxRequest true "Update Data"
// @Security ApiKeyAuth
// @Router /taxes/{id} [put]
func (h *Handler) UpdateTax(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateTaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := UpdateTaxDTO{
		Name:     req.Name,
		Rate:     req.Rate,
		IsActive: req.IsActive,
	}

	if err := h.svc.UpdateTax(r.Context(), id, dto); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Tax Status
// @Tags taxes
// @Accept json
// @Produce json
// @Param id path string true "Tax ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /taxes/{id}/status [patch]
func (h *Handler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
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

// @Summary Delete Tax
// @Tags taxes
// @Security ApiKeyAuth
// @Param id path string true "Tax ID"
// @Router /taxes/{id} [delete]
func (h *Handler) DeleteTax(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteTax(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapTaxToResponse(t *Tax) TaxResponse {
	return TaxResponse{
		ID:       t.ID,
		Name:     t.Name,
		Rate:     t.Rate,
		IsActive: t.IsActive,
	}
}
