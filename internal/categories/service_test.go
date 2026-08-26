package categories_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/categories"
)

type mockRepository struct {
	categories.Repository
	mockGetByID func(ctx context.Context, id string) (*categories.Category, error)
	mockCreate  func(ctx context.Context, c *categories.Category) error
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*categories.Category, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepository) Create(ctx context.Context, c *categories.Category) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, c)
	}
	return errors.New("not implemented")
}

func TestService_GetCategory(t *testing.T) {
	repo := &mockRepository{
		mockGetByID: func(ctx context.Context, id string) (*categories.Category, error) {
			if id == "123" {
				return &categories.Category{ID: "123", Name: "Beverages"}, nil
			}
			return nil, errors.New("category not found")
		},
	}

	svc := categories.NewService(repo)

	c, err := svc.GetCategory(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if c.Name != "Beverages" {
		t.Errorf("expected name Beverages, got %s", c.Name)
	}
}

func TestService_CreateCategory(t *testing.T) {
	repo := &mockRepository{
		mockCreate: func(ctx context.Context, c *categories.Category) error {
			c.ID = "new-id"
			return nil
		},
	}

	svc := categories.NewService(repo)
	
	dto := categories.CreateCategoryDTO{
		Name:        "Snacks",
		Description: "Light meals",
	}

	c, err := svc.CreateCategory(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", c.ID)
	}
}
