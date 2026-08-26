package variants_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/variants"
)

type mockModifierRepository struct {
	variants.ModifierRepository
	mockGetByID func(ctx context.Context, id string) (*variants.Modifier, error)
	mockCreate  func(ctx context.Context, m *variants.Modifier) error
}

func (m *mockModifierRepository) GetByID(ctx context.Context, id string) (*variants.Modifier, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockModifierRepository) Create(ctx context.Context, mod *variants.Modifier) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, mod)
	}
	return errors.New("not implemented")
}

func TestModifierService_GetModifier(t *testing.T) {
	repo := &mockModifierRepository{
		mockGetByID: func(ctx context.Context, id string) (*variants.Modifier, error) {
			if id == "123" {
				return &variants.Modifier{ID: "123", Name: "Vanilla"}, nil
			}
			return nil, errors.New("not found")
		},
	}

	svc := variants.NewModifierService(repo)

	m, err := svc.GetModifier(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if m.Name != "Vanilla" {
		t.Errorf("expected name Vanilla, got %s", m.Name)
	}
}

func TestModifierService_CreateModifier(t *testing.T) {
	repo := &mockModifierRepository{
		mockCreate: func(ctx context.Context, m *variants.Modifier) error {
			m.ID = "new-id"
			return nil
		},
	}

	svc := variants.NewModifierService(repo)

	dto := variants.ModifierDTO{
		ModifierGroupID: "group1",
		Name:            "Vanilla",
		AdditionalPrice: 5000,
	}

	id, err := svc.CreateModifier(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", id)
	}
}
