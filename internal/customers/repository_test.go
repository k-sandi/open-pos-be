package customers

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	repo := NewRepository(mock)
	
	email := "test@test.com"
	rows := mock.NewRows([]string{"id", "name", "phone", "email", "loyalty_points", "is_active"}).
		AddRow("1", "John", "12345", &email, 10, true)

	mock.ExpectQuery("^SELECT id, name, phone, email, loyalty_points, is_active FROM customers WHERE id = \\$1 AND deleted_at IS NULL$").
		WithArgs("1").
		WillReturnRows(rows)

	customer, err := repo.GetByID(context.Background(), "1")
	if err != nil {
		t.Errorf("error was not expected while getting customer: %s", err)
	}
	
	if customer.ID != "1" || customer.Name != "John" {
		t.Errorf("expected customer John, got %v", customer)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_UpdatePoints(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %s", err)
	}
	defer mock.Close(context.Background())

	repo := NewRepository(mock)
	
	mock.ExpectExec("^UPDATE customers SET loyalty_points = loyalty_points \\+ \\$1, updated_at = NOW\\(\\) WHERE id=\\$2 AND deleted_at IS NULL$").
		WithArgs(50, "1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdatePoints(context.Background(), "1", 50)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
