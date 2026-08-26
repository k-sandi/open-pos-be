package users_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/users"
)

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := users.NewRepository(mock)

	// Expectation
	mock.ExpectQuery("^SELECT (.+) FROM users u JOIN roles r").
		WithArgs("123").
		WillReturnRows(mock.NewRows([]string{"id", "employee_id", "pin_hash", "name", "email", "phone", "role_id", "role_name", "is_active"}).
			AddRow("123", "EMP01", "hash", "John Doe", "john@example.com", "123456789", "role1", "Admin", true))

	user, err := repo.GetByID(context.Background(), "123")
	if err != nil {
		t.Errorf("error was not expected while getting user: %s", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name %s, got %s", "John Doe", user.Name)
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

	repo := users.NewRepository(mock)

	u := &users.User{
		EmployeeID: "EMP02",
		PINHash:    "hash",
		Name:       "Jane Doe",
		Email:      "jane@example.com",
		Phone:      "987654321",
		RoleID:     "role2",
		IsActive:   true,
	}

	mock.ExpectQuery("^INSERT INTO users").
		WithArgs(u.EmployeeID, u.PINHash, u.Name, u.Email, u.Phone, u.RoleID, u.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("456"))

	err = repo.Create(context.Background(), u)
	if err != nil {
		t.Errorf("error was not expected while creating user: %s", err)
	}

	if u.ID != "456" {
		t.Errorf("expected ID %s, got %s", "456", u.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
