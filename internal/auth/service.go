package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"open-pos-be/internal/jwt"
	"open-pos-be/internal/users"
)

type Service interface {
	Login(ctx context.Context, employeeID, pin string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}

type service struct {
	userRepo users.Repository
}

func NewService(userRepo users.Repository) Service {
	return &service{userRepo: userRepo}
}

func (s *service) Login(ctx context.Context, employeeID, pin string) (string, string, error) {
	u, err := s.userRepo.GetByEmployeeID(ctx, employeeID)
	if err != nil || !u.IsActive {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PINHash), []byte(pin)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, _ := jwt.GenerateToken(u.ID, u.RoleName, 15*time.Minute)
	refreshToken, _ := jwt.GenerateToken(u.ID, u.RoleName, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, _ := jwt.GenerateToken(claims.UserID, claims.Role, 15*time.Minute)
	newRefreshToken, _ := jwt.GenerateToken(claims.UserID, claims.Role, 7*24*time.Hour)

	return accessToken, newRefreshToken, nil
}
