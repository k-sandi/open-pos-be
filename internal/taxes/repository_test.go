package taxes_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/taxes"
)

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := taxes.NewRepository(mock)

	mock.ExpectQuery("^SELECT (.+) FROM taxes").
		WithArgs("123").
		WillReturnRows(mock.NewRows([]string{"id", "name", "rate", "is_active"}).
			AddRow("123", "VAT", 11.0, true))

	tax, err := repo.GetByID(context.Background(), "123")
	if err != nil {
		t.Errorf("error was not expected while getting tax: %s", err)
	}

	if tax.Name != "VAT" {
		t.Errorf("expected name %s, got %s", "VAT", tax.Name)
	}
	if tax.Rate != 11.0 {
		t.Errorf("expected rate %f, got %f", 11.0, tax.Rate)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := taxes.NewRepository(mock)

	tax := &taxes.Tax{
		Name:     "Service Tax",
		Rate:     5.0,
		IsActive: true,
	}

	mock.ExpectQuery("^INSERT INTO taxes").
		WithArgs(tax.Name, tax.Rate, tax.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("456"))

	err = repo.Create(context.Background(), tax)
	if err != nil {
		t.Errorf("error was not expected while creating tax: %s", err)
	}

	if tax.ID != "456" {
		t.Errorf("expected ID %s, got %s", "456", tax.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
