package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	customMiddleware "open-pos-be/internal/middleware"
)

// StatusRequest represents the payload for updating a user's active status.
type StatusRequest struct {
	IsActive bool `json:"is_active"`
}

type UserResponse struct {
	ID         string `json:"id"`
	EmployeeID string `json:"employee_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RoleID     string `json:"role_id"`
	RoleName   string `json:"role_name"`
	IsActive   bool   `json:"is_active"`
}

type CreateUserRequest struct {
	EmployeeID string `json:"employee_id"`
	PIN        string `json:"pin"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RoleID     string `json:"role_id"`
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// Routes returns a Chi router mounted with user routes.
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/me", h.GetMe)
	r.Get("/", h.ListUsers)
	r.Post("/", h.CreateUser)
	r.Put("/{id}", h.UpdateUser)
	r.Delete("/{id}", h.DeleteUser)
	r.Patch("/{id}/status", h.UpdateStatusHandler)
	return r
}

// @Summary Get Current User Profile
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Router /users/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := customMiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapUserToResponse(user))
}

// @Summary List Users
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Router /users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var res []UserResponse
	for _, u := range users {
		res = append(res, mapUserToResponse(u))
	}
	if res == nil {
		res = []UserResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// @Summary Create User
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User Data"
// @Security ApiKeyAuth
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := CreateUserDTO{
		EmployeeID: req.EmployeeID,
		PIN:        req.PIN,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		RoleID:     req.RoleID,
	}

	_, err := h.svc.CreateUser(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update User
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body CreateUserRequest true "Update Data"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := UpdateUserDTO{
		EmployeeID: req.EmployeeID,
		PIN:        req.PIN,
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		RoleID:     req.RoleID,
	}

	if err := h.svc.UpdateUser(r.Context(), id, dto); err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Update User Status
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body StatusRequest true "Status Update Request"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Security ApiKeyAuth
// @Router /users/{id}/status [patch]
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

// @Summary Delete User
// @Tags users
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func mapUserToResponse(u *User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		EmployeeID: u.EmployeeID,
		Name:       u.Name,
		Email:      u.Email,
		Phone:      u.Phone,
		RoleID:     u.RoleID,
		RoleName:   u.RoleName,
		IsActive:   u.IsActive,
	}
}
