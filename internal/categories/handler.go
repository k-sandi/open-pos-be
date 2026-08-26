package categories

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Routes returns a Chi router mounted with category routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListCategories)
	r.Get("/{id}", h.GetCategory)
	r.Post("/", h.CreateCategory)
	r.Put("/{id}", h.UpdateCategory)
	r.Patch("/{id}/status", h.UpdateStatusHandler)
	r.Delete("/{id}", h.DeleteCategory)
	return r
}

// @Summary List Categories
// @Tags categories
// @Produce json
// @Security ApiKeyAuth
// @Router /categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategories(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var res []CategoryResponse
	for _, c := range categories {
		res = append(res, mapCategoryToResponse(c))
	}
	if res == nil {
		res = []CategoryResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @Summary Get Category
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Security ApiKeyAuth
// @Router /categories/{id} [get]
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	category, err := h.svc.GetCategory(r.Context(), id)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapCategoryToResponse(category))
}

// @Summary Create Category
// @Tags categories
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Category Data"
// @Security ApiKeyAuth
// @Router /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := CreateCategoryDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	_, err := h.svc.CreateCategory(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update Category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body CreateCategoryRequest true "Update Data"
// @Security ApiKeyAuth
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := UpdateCategoryDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.svc.UpdateCategory(r.Context(), id, dto); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Category Status
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /categories/{id}/status [patch]
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

// @Summary Delete Category
// @Tags categories
// @Security ApiKeyAuth
// @Param id path string true "Category ID"
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteCategory(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapCategoryToResponse(c *Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		IsActive:    c.IsActive,
	}
}
