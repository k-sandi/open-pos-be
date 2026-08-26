package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	
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
	IsActive   bool   `json:"is_active"`
}

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
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
func (h *Handler) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: DB update logic goes here. Mocking success for now.
	w.WriteHeader(http.StatusOK)
}

// @Summary Get Current User Profile
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Router /users/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
    // Extract claims from context
	userID, ok := customMiddleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user UserResponse
	err := h.db.QueryRow(r.Context(), `
		SELECT id, employee_id, name, email, phone, role_id, is_active
		FROM users WHERE id = $1 AND deleted_at IS NULL
	`, userID).Scan(&user.ID, &user.EmployeeID, &user.Name, &user.Email, &user.Phone, &user.RoleID, &user.IsActive)
	
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Summary List Users
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Router /users [get]
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT id, employee_id, name, email, phone, role_id, is_active
		FROM users WHERE deleted_at IS NULL
	`)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []UserResponse{}
	for rows.Next() {
		var u UserResponse
		if err := rows.Scan(&u.ID, &u.EmployeeID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.IsActive); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

type CreateUserRequest struct {
	EmployeeID string `json:"employee_id"`
	PIN        string `json:"pin"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	RoleID     string `json:"role_id"`
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
	
	pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)

	_, err := h.db.Exec(r.Context(), `
		INSERT INTO users (employee_id, pin_hash, name, email, phone, role_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, req.EmployeeID, string(pinHash), req.Name, req.Email, req.Phone, req.RoleID)
	
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

    pinHash, _ := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
    
	_, err := h.db.Exec(r.Context(), `
		UPDATE users 
		SET employee_id=$1, pin_hash=$2, name=$3, email=$4, phone=$5, role_id=$6, updated_at=NOW()
		WHERE id=$7 AND deleted_at IS NULL
	`, req.EmployeeID, string(pinHash), req.Name, req.Email, req.Phone, req.RoleID, id)
	
	if err != nil {
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
	_, err := h.db.Exec(r.Context(), `
		UPDATE users SET deleted_at=NOW() WHERE id=$1
	`, id)
	if err != nil {
		http.Error(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
