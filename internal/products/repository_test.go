package products_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/products"
)

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := products.NewRepository(mock)

	mock.ExpectQuery("^SELECT (.+) FROM products").
		WithArgs("1").
		WillReturnRows(mock.NewRows([]string{"id", "category_id", "sku", "name", "description", "price", "image_url", "is_active"}).
			AddRow("1", "cat1", "SKU1", "Coffee", "Desc", int64(100), "img.png", true))

	product, err := repo.GetByID(context.Background(), "1")
	if err != nil {
		t.Errorf("error was not expected while getting product: %s", err)
	}

	if product.Name != "Coffee" {
		t.Errorf("expected name %s, got %s", "Coffee", product.Name)
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

	repo := products.NewRepository(mock)

	p := &products.Product{
		CategoryID:  "cat1",
		SKU:         "SKU1",
		Name:        "Coffee",
		Description: "Desc",
		Price:       100,
		ImageURL:    "img.png",
		IsActive:    true,
	}

	mock.ExpectQuery("^INSERT INTO products").
		WithArgs(p.CategoryID, p.SKU, p.Name, p.Description, p.Price, p.ImageURL, p.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("1"))

	err = repo.Create(context.Background(), p)
	if err != nil {
		t.Errorf("error was not expected while creating product: %s", err)
	}

	if p.ID != "1" {
		t.Errorf("expected ID %s, got %s", "1", p.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := products.NewRepository(mock)

	p := &products.Product{
		ID:          "1",
		CategoryID:  "cat1",
		SKU:         "SKU1",
		Name:        "Coffee",
		Description: "Desc",
		Price:       100,
		ImageURL:    "img.png",
		IsActive:    true,
	}

	mock.ExpectExec("^UPDATE products").
		WithArgs(p.CategoryID, p.SKU, p.Name, p.Description, p.Price, p.ImageURL, p.IsActive, p.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), p)
	if err != nil {
		t.Errorf("error was not expected while updating product: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
