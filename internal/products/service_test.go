package products_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/products"
)

type mockRepository struct {
	products.Repository
	mockGetByID func(ctx context.Context, id string) (*products.Product, error)
	mockCreate  func(ctx context.Context, p *products.Product) error
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*products.Product, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepository) Create(ctx context.Context, p *products.Product) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, p)
	}
	return errors.New("not implemented")
}

func TestService_GetProduct(t *testing.T) {
	repo := &mockRepository{
		mockGetByID: func(ctx context.Context, id string) (*products.Product, error) {
			if id == "1" {
				return &products.Product{ID: "1", Name: "Coffee"}, nil
			}
			return nil, errors.New("product not found")
		},
	}

	svc := products.NewService(repo)

	p, err := svc.GetProduct(context.Background(), "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Name != "Coffee" {
		t.Errorf("expected name Coffee, got %s", p.Name)
	}
}

func TestService_CreateProduct(t *testing.T) {
	repo := &mockRepository{
		mockCreate: func(ctx context.Context, p *products.Product) error {
			p.ID = "new-id"
			return nil
		},
	}

	svc := products.NewService(repo)
	
	dto := products.CreateProductDTO{
		CategoryID:  "cat1",
		SKU:         "SKU1",
		Name:        "Coffee",
		Description: "Desc",
		Price:       100,
		ImageURL:    "img.png",
	}

	p, err := svc.CreateProduct(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", p.ID)
	}
	
	if !p.IsActive {
		t.Errorf("expected IsActive to be true")
	}
}
