package variants_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"open-pos-be/internal/variants"
)

func TestModifierRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := variants.NewModifierRepository(mock)

	mock.ExpectQuery("^SELECT id, modifier_group_id, name, additional_price, is_active FROM modifiers").
		WithArgs("123").
		WillReturnRows(mock.NewRows([]string{"id", "modifier_group_id", "name", "additional_price", "is_active"}).
			AddRow("123", "group1", "Vanilla", int64(5000), true))

	m, err := repo.GetByID(context.Background(), "123")
	if err != nil {
		t.Errorf("error was not expected while getting modifier: %s", err)
	}

	if m.Name != "Vanilla" {
		t.Errorf("expected name %s, got %s", "Vanilla", m.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestModifierRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := variants.NewModifierRepository(mock)

	m := &variants.Modifier{
		ModifierGroupID: "group1",
		Name:            "Chocolate",
		AdditionalPrice: 6000,
		IsActive:        true,
	}

	mock.ExpectQuery("^INSERT INTO modifiers").
		WithArgs(m.ModifierGroupID, m.Name, m.AdditionalPrice, m.IsActive).
		WillReturnRows(mock.NewRows([]string{"id"}).AddRow("456"))

	err = repo.Create(context.Background(), m)
	if err != nil {
		t.Errorf("error was not expected while creating modifier: %s", err)
	}

	if m.ID != "456" {
		t.Errorf("expected ID %s, got %s", "456", m.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
