package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", h.Login)
	return r
}

type LoginRequest struct {
	EmployeeID string `json:"employee_id"`
	PIN        string `json:"pin"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary Login User
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userID, roleName, pinHash string
	err := h.db.QueryRow(r.Context(), `
		SELECT u.id, r.name, u.pin_hash 
		FROM users u 
		JOIN roles r ON u.role_id = r.id 
		WHERE u.employee_id = $1 AND u.deleted_at IS NULL AND u.is_active = true
	`, req.EmployeeID).Scan(&userID, &roleName, &pinHash)
	
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pinHash), []byte(req.PIN)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, _ := GenerateToken(userID, roleName, 15*time.Minute)
	refreshToken, _ := GenerateToken(userID, roleName, 7*24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
