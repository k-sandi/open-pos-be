package variants_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/variants"
)

type mockModifierGroupRepository struct {
	variants.ModifierGroupRepository
	mockGetByID func(ctx context.Context, id string) (*variants.ModifierGroup, error)
	mockCreate  func(ctx context.Context, mg *variants.ModifierGroup) error
}

func (m *mockModifierGroupRepository) GetByID(ctx context.Context, id string) (*variants.ModifierGroup, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockModifierGroupRepository) Create(ctx context.Context, mg *variants.ModifierGroup) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, mg)
	}
	return errors.New("not implemented")
}

func TestModifierGroupService_GetByID(t *testing.T) {
	repo := &mockModifierGroupRepository{
		mockGetByID: func(ctx context.Context, id string) (*variants.ModifierGroup, error) {
			if id == "123" {
				return &variants.ModifierGroup{ID: "123", Name: "Ice Level"}, nil
			}
			return nil, errors.New("modifier group not found")
		},
	}

	svc := variants.NewModifierGroupService(repo)

	mg, err := svc.GetByID(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if mg.Name != "Ice Level" {
		t.Errorf("expected name Ice Level, got %s", mg.Name)
	}
}

func TestModifierGroupService_Create(t *testing.T) {
	repo := &mockModifierGroupRepository{
		mockCreate: func(ctx context.Context, mg *variants.ModifierGroup) error {
			mg.ID = "new-id"
			return nil
		},
	}

	svc := variants.NewModifierGroupService(repo)
	
	dto := variants.ModifierGroupDTO{
		ProductID:  "prod-1",
		Name:       "Sugar Level",
		MinChoices: 1,
		MaxChoices: 1,
		IsActive:   true,
	}

	mg, err := svc.Create(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mg.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", mg.ID)
	}
}
