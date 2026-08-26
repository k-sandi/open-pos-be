package products

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type ProductResponse struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	ImageURL    string `json:"image_url"`
	IsActive    bool   `json:"is_active"`
}

type CreateProductRequest struct {
	CategoryID  string `json:"category_id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	ImageURL    string `json:"image_url"`
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListProducts)
	r.Get("/{id}", h.GetProduct)
	r.Post("/", h.CreateProduct)
	r.Put("/{id}", h.UpdateProduct)
	r.Delete("/{id}", h.DeleteProduct)
	r.Patch("/{id}/status", h.UpdateStatusHandler)
	return r
}

// @Summary List Products
// @Tags products
// @Produce json
// @Security ApiKeyAuth
// @Router /products [get]
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var res []ProductResponse
	for _, p := range products {
		res = append(res, mapProductToResponse(p))
	}
	if res == nil {
		res = []ProductResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @Summary Get Product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Security ApiKeyAuth
// @Router /products/{id} [get]
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapProductToResponse(product))
}

// @Summary Create Product
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Product Data"
// @Security ApiKeyAuth
// @Router /products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := CreateProductDTO{
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
	}

	_, err := h.svc.CreateProduct(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update Product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body CreateProductRequest true "Update Data"
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := UpdateProductDTO{
		CategoryID:  req.CategoryID,
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
	}

	if err := h.svc.UpdateProduct(r.Context(), id, dto); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Product Status
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /products/{id}/status [patch]
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

// @Summary Delete Product
// @Tags products
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapProductToResponse(p *Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
	}
}
