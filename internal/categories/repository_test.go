package categories_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/categories"
)

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := categories.NewRepository(mock)

	mock.ExpectQuery("^SELECT (.+) FROM categories").
		WithArgs("123").
		WillReturnRows(mock.NewRows([]string{"id", "name", "description", "is_active"}).
			AddRow("123", "Beverages", "Drinks and more", true))

	cat, err := repo.GetByID(context.Background(), "123")
	if err != nil {
		t.Errorf("error was not expected while getting category: %s", err)
	}

	if cat.Name != "Beverages" {
		t.Errorf("expected name %s, got %s", "Beverages", cat.Name)
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

	repo := categories.NewRepository(mock)

	c := &categories.Category{
		Name:        "Snacks",
		Description: "Light meals",
		IsActive:    true,
	}

	mock.ExpectQuery("^INSERT INTO categories").
		WithArgs(c.Name, c.Description, c.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("456"))

	err = repo.Create(context.Background(), c)
	if err != nil {
		t.Errorf("error was not expected while creating category: %s", err)
	}

	if c.ID != "456" {
		t.Errorf("expected ID %s, got %s", "456", c.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
