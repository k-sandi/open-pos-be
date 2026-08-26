package users_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"open-pos-be/internal/users"
)

func TestUpdateStatusHandler(t *testing.T) {
	req := httptest.NewRequest("PATCH", "/api/v1/users/123/status", strings.NewReader(`{"is_active": false}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	users.UpdateStatusHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestUpdateStatusHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("PATCH", "/api/v1/users/123/status", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	users.UpdateStatusHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRoutes_PatchUserStatus(t *testing.T) {
	r := users.Routes()

	req := httptest.NewRequest("PATCH", "/123/status", strings.NewReader(`{"is_active": true}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
