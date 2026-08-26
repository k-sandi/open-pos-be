package taxes_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/taxes"
)

type mockRepository struct {
	taxes.Repository
	mockGetByID func(ctx context.Context, id string) (*taxes.Tax, error)
	mockCreate  func(ctx context.Context, t *taxes.Tax) error
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*taxes.Tax, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepository) Create(ctx context.Context, t *taxes.Tax) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, t)
	}
	return errors.New("not implemented")
}

func TestService_GetTax(t *testing.T) {
	repo := &mockRepository{
		mockGetByID: func(ctx context.Context, id string) (*taxes.Tax, error) {
			if id == "123" {
				return &taxes.Tax{ID: "123", Name: "VAT", Rate: 11.0}, nil
			}
			return nil, errors.New("tax not found")
		},
	}

	svc := taxes.NewService(repo)

	tax, err := svc.GetTax(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tax.Name != "VAT" {
		t.Errorf("expected name VAT, got %s", tax.Name)
	}
}

func TestService_CreateTax(t *testing.T) {
	repo := &mockRepository{
		mockCreate: func(ctx context.Context, tx *taxes.Tax) error {
			tx.ID = "new-id"
			return nil
		},
	}

	svc := taxes.NewService(repo)
	
	dto := taxes.CreateTaxDTO{
		Name:     "Service Tax",
		Rate:     5.0,
		IsActive: true,
	}

	tax, err := svc.CreateTax(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if tax.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", tax.ID)
	}
	if tax.Rate != 5.0 {
		t.Errorf("expected Rate 5.0, got %f", tax.Rate)
	}
}
