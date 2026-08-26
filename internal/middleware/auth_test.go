package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"open-pos-be/internal/auth"
	"open-pos-be/internal/middleware"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("missing x-api-key returns 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		rr := httptest.NewRecorder()

		handlerCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		middleware.Auth(next).ServeHTTP(rr, req)

		if handlerCalled {
			t.Fatal("expected handler not to be called")
		}
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("x-api-key", "invalid.jwt.token")
		rr := httptest.NewRecorder()

		handlerCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		middleware.Auth(next).ServeHTTP(rr, req)

		if handlerCalled {
			t.Fatal("expected handler not to be called")
		}
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("expired token returns 401", func(t *testing.T) {
		token, err := auth.GenerateToken("user123", "admin", -1*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}

		req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("x-api-key", token)
		rr := httptest.NewRecorder()

		handlerCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlerCalled = true
		})

		middleware.Auth(next).ServeHTTP(rr, req)

		if handlerCalled {
			t.Fatal("expected handler not to be called")
		}
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", rr.Code)
		}
	})

	t.Run("valid token injects user id and role into context", func(t *testing.T) {
		token, err := auth.GenerateToken("user_abc", "cashier", 15*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error generating token: %v", err)
		}

		req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
		req.Header.Set("x-api-key", token)
		rr := httptest.NewRecorder()

		var recordedUserID string
		var recordedRole string
		var hasUserID, hasRole bool

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recordedUserID, hasUserID = middleware.GetUserID(r.Context())
			recordedRole, hasRole = middleware.GetRole(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		middleware.Auth(next).ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}
		if !hasUserID || recordedUserID != "user_abc" {
			t.Errorf("expected userID 'user_abc', got '%s' (found: %v)", recordedUserID, hasUserID)
		}
		if !hasRole || recordedRole != "cashier" {
			t.Errorf("expected role 'cashier', got '%s' (found: %v)", recordedRole, hasRole)
		}
	})
}
