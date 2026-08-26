package auth_test

import (
	"testing"
	"time"

	"open-pos-be/internal/auth"
)

func TestGenerateToken(t *testing.T) {
	token, err := auth.GenerateToken("user123", "admin", 15*time.Minute)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected token string, got empty")
	}
}

func TestValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		token, err := auth.GenerateToken("user123", "admin", 15*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			t.Fatalf("expected token to be valid, got error: %v", err)
		}
		if claims.UserID != "user123" {
			t.Errorf("expected userID 'user123', got '%s'", claims.UserID)
		}
		if claims.Role != "admin" {
			t.Errorf("expected role 'admin', got '%s'", claims.Role)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token, err := auth.GenerateToken("user123", "admin", -1*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}

		claims, err := auth.ValidateToken(token)
		if err == nil {
			t.Fatalf("expected error for expired token, got claims: %+v", claims)
		}
	})

	t.Run("invalid token string", func(t *testing.T) {
		_, err := auth.ValidateToken("invalid.token.string")
		if err == nil {
			t.Fatal("expected error for invalid token string, got nil")
		}
	})
}
