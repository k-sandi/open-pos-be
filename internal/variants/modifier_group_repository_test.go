package variants_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/variants"
)

func TestModifierGroupRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := variants.NewModifierGroupRepository(mock)

	mock.ExpectQuery("^SELECT (.+) FROM modifier_groups").
		WithArgs("123").
		WillReturnRows(mock.NewRows([]string{"id", "product_id", "name", "min_choices", "max_choices", "is_active"}).
			AddRow("123", "prod-1", "Ice Level", 0, 1, true))

	mg, err := repo.GetByID(context.Background(), "123")
	if err != nil {
		t.Errorf("error was not expected while getting modifier group: %s", err)
	}

	if mg.Name != "Ice Level" {
		t.Errorf("expected name %s, got %s", "Ice Level", mg.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModifierGroupRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := variants.NewModifierGroupRepository(mock)

	mg := &variants.ModifierGroup{
		ProductID:  "prod-1",
		Name:       "Sugar Level",
		MinChoices: 1,
		MaxChoices: 1,
		IsActive:   true,
	}

	mock.ExpectQuery("^INSERT INTO modifier_groups").
		WithArgs(mg.ProductID, mg.Name, mg.MinChoices, mg.MaxChoices, mg.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("456"))

	err = repo.Create(context.Background(), mg)
	if err != nil {
		t.Errorf("error was not expected while creating modifier group: %s", err)
	}

	if mg.ID != "456" {
		t.Errorf("expected ID %s, got %s", "456", mg.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
