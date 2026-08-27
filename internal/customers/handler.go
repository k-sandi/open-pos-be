package customers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CustomerResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Email         *string `json:"email,omitempty"`
	LoyaltyPoints int     `json:"loyalty_points"`
	IsActive      bool    `json:"is_active"`
}

type CreateCustomerRequest struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone"`
	Email         *string `json:"email,omitempty"`
	LoyaltyPoints int     `json:"loyalty_points"`
}

type UpdateCustomerRequest struct {
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    *string `json:"email,omitempty"`
	IsActive bool    `json:"is_active"`
}

type UpdatePointsRequest struct {
	Points int `json:"points"`
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Routes returns a Chi router mounted with customer routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.ListCustomers)
	r.Post("/", h.CreateCustomer)
	r.Get("/{id}", h.GetCustomer)
	r.Put("/{id}", h.UpdateCustomer)
	r.Delete("/{id}", h.DeleteCustomer)
	r.Patch("/{id}/points", h.UpdatePointsHandler)
	return r
}

// @Summary List Customers
// @Tags customers
// @Produce json
// @Param query query string false "Search query by name or phone"
// @Security ApiKeyAuth
// @Router /customers [get]
func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	var customers []*Customer
	var err error

	if query != "" {
		customers, err = h.svc.SearchCustomers(r.Context(), query)
	} else {
		customers, err = h.svc.ListCustomers(r.Context())
	}

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var res []CustomerResponse
	for _, c := range customers {
		res = append(res, mapCustomerToResponse(c))
	}
	if res == nil {
		res = []CustomerResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @Summary Create Customer
// @Tags customers
// @Accept json
// @Produce json
// @Param request body CreateCustomerRequest true "Customer Data"
// @Security ApiKeyAuth
// @Router /customers [post]
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := CreateCustomerDTO{
		Name:          req.Name,
		Phone:         req.Phone,
		Email:         req.Email,
		LoyaltyPoints: req.LoyaltyPoints,
	}

	_, err := h.svc.CreateCustomer(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get Customer
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Security ApiKeyAuth
// @Router /customers/{id} [get]
func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customer, err := h.svc.GetCustomer(r.Context(), id)
	if err != nil {
		http.Error(w, "customer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapCustomerToResponse(customer))
}

// @Summary Update Customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body UpdateCustomerRequest true "Update Data"
// @Security ApiKeyAuth
// @Router /customers/{id} [put]
func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := UpdateCustomerDTO{
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		IsActive: req.IsActive,
	}

	if err := h.svc.UpdateCustomer(r.Context(), id, dto); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update Customer Points
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body UpdatePointsRequest true "Points Delta"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /customers/{id}/points [patch]
func (h *Handler) UpdatePointsHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdatePointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdatePoints(r.Context(), id, req.Points); err != nil {
		if err.Error() == "insufficient loyalty points" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete Customer
// @Tags customers
// @Security ApiKeyAuth
// @Router /customers/{id} [delete]
func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteCustomer(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapCustomerToResponse(c *Customer) CustomerResponse {
	return CustomerResponse{
		ID:            c.ID,
		Name:          c.Name,
		Phone:         c.Phone,
		Email:         c.Email,
		LoyaltyPoints: c.LoyaltyPoints,
		IsActive:      c.IsActive,
	}
}
