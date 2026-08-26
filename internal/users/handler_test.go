package users_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-pos-be/internal/users"
)

type mockService struct{}

func (m *mockService) GetProfile(ctx context.Context, userID string) (*users.User, error) { return nil, nil }
func (m *mockService) ListUsers(ctx context.Context) ([]*users.User, error) { return nil, nil }
func (m *mockService) CreateUser(ctx context.Context, dto users.CreateUserDTO) (*users.User, error) { return nil, nil }
func (m *mockService) UpdateUser(ctx context.Context, id string, dto users.UpdateUserDTO) error { return nil }
func (m *mockService) UpdateStatus(ctx context.Context, id string, isActive bool) error { return nil }
func (m *mockService) DeleteUser(ctx context.Context, id string) error { return nil }


func TestUpdateStatusHandler(t *testing.T) {
	h := users.NewHandler(&mockService{})
	req := httptest.NewRequest("PATCH", "/api/v1/users/123/status", strings.NewReader(`{"is_active": false}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateStatusHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestUpdateStatusHandler_InvalidBody(t *testing.T) {
	h := users.NewHandler(&mockService{})
	req := httptest.NewRequest("PATCH", "/api/v1/users/123/status", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	h.UpdateStatusHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRoutes_PatchUserStatus(t *testing.T) {
	h := users.NewHandler(&mockService{})
	r := h.Routes()

	req := httptest.NewRequest("PATCH", "/123/status", strings.NewReader(`{"is_active": true}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
