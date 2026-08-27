package orders

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"open-pos-be/internal/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Auth)
	r.Post("/", h.CreateOrder)
	return r
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Creates a new order, calculates tax and subtotal, and saves order items and modifiers
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body CreateOrderRequest true "Order Request"
// @Success      201  {object}  Order
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /orders [post]
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	cashierID, ok := middleware.GetUserID(r.Context())
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	order, err := h.service.CreateOrder(r.Context(), cashierID, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
