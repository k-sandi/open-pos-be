package variants_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"open-pos-be/internal/variants"
)

type mockModifierService struct {
	variants.ModifierService
	mockList   func(ctx context.Context, groupID string) ([]*variants.Modifier, error)
	mockCreate func(ctx context.Context, dto variants.ModifierDTO) (string, error)
}

func (m *mockModifierService) ListModifiers(ctx context.Context, groupID string) ([]*variants.Modifier, error) {
	if m.mockList != nil {
		return m.mockList(ctx, groupID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockModifierService) CreateModifier(ctx context.Context, dto variants.ModifierDTO) (string, error) {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, dto)
	}
	return "", errors.New("not implemented")
}

func TestModifierHandler_ListModifiers(t *testing.T) {
	svc := &mockModifierService{
		mockList: func(ctx context.Context, groupID string) ([]*variants.Modifier, error) {
			return []*variants.Modifier{
				{ID: "1", Name: "Vanilla", ModifierGroupID: groupID},
			}, nil
		},
	}

	handler := variants.NewModifierHandler(svc)
	r := chi.NewRouter()
	r.Mount("/modifiers", handler.Routes())

	req := httptest.NewRequest("GET", "/modifiers/group/group1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status OK, got %v", status)
	}

	if !bytes.Contains(rr.Body.Bytes(), []byte("Vanilla")) {
		t.Errorf("expected response to contain Vanilla")
	}
}

func TestModifierHandler_CreateModifier(t *testing.T) {
	svc := &mockModifierService{
		mockCreate: func(ctx context.Context, dto variants.ModifierDTO) (string, error) {
			return "new-id", nil
		},
	}

	handler := variants.NewModifierHandler(svc)
	r := chi.NewRouter()
	r.Mount("/modifiers", handler.Routes())

	body := []byte(`{"modifier_group_id": "group1", "name": "Vanilla", "additional_price": 5000}`)
	req := httptest.NewRequest("POST", "/modifiers", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status Created, got %v", status)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if res["id"] != "new-id" {
		t.Errorf("expected id new-id, got %v", res["id"])
	}
}
