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
	r.Post("/logout", h.Logout)
	r.Post("/refresh", h.Refresh)
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

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// @Summary Refresh Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh Request"
// @Success 200 {object} LoginResponse
// @Router /auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := ValidateToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, _ := GenerateToken(claims.UserID, claims.Role, 15*time.Minute)
	// Optionally generate a new refresh token, or keep the old one
	newRefreshToken, _ := GenerateToken(claims.UserID, claims.Role, 7*24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}

// @Summary Logout User
// @Tags auth
// @Produce json
// @Success 200 {string} string "OK"
// @Router /auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// For stateless JWT, logout is primarily handled by the client clearing the token.
	// We just return a success response.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "logged out successfully"}`))
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
