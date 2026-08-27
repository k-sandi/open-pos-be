package orders

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
)

func TestRepository_GetProductPrice(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewRepository(mock)

	productID := "prod-1"
	expectedPrice := int64(15000)

	mock.ExpectQuery("SELECT price FROM products WHERE id = \\$1").
		WithArgs(productID).
		WillReturnRows(pgxmock.NewRows([]string{"price"}).AddRow(expectedPrice))

	price, err := repo.GetProductPrice(context.Background(), productID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if price != expectedPrice {
		t.Errorf("expected price %v, got %v", expectedPrice, price)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_CreateOrderTx(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewRepository(mock)

	taxID := "tax-1"
	order := &Order{
		OrderNumber:   "ORD-123",
		CashierID:     "cashier-1",
		TaxID:         &taxID,
		Subtotal:      20000,
		TaxAmount:     2000,
		TotalAmount:   22000,
		PaymentMethod: "CASH",
		Status:        "PAID",
		Items: []OrderItem{
			{
				ProductID: "prod-1",
				Quantity:  1,
				UnitPrice: 15000,
				Subtotal:  20000,
				Modifiers: []OrderItemModifier{
					{
						ModifierID: "mod-1",
						Price:      5000,
					},
				},
			},
		},
	}

	mock.ExpectBegin()

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(order.OrderNumber, order.CashierID, order.TaxID, order.Subtotal, order.TaxAmount, order.TotalAmount, order.PaymentMethod, order.Status).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("order-1", time.Now(), time.Now()))

	mock.ExpectQuery("INSERT INTO order_items").
		WithArgs("order-1", "prod-1", 1, int64(15000), int64(20000)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("item-1"))

	mock.ExpectQuery("INSERT INTO order_item_modifiers").
		WithArgs("item-1", "mod-1", int64(5000)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("mod-item-1"))

	mock.ExpectCommit()

	err = repo.CreateOrderTx(context.Background(), order)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if order.ID != "order-1" {
		t.Errorf("expected order ID 'order-1', got %v", order.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewRepository(mock)

	orderID := "order-1"
	taxID := "tax-1"
	tm := time.Now()

	mock.ExpectQuery("SELECT id, order_number, cashier_id, tax_id, subtotal, tax_amount, total_amount, payment_method, status, created_at, updated_at FROM orders").
		WithArgs(orderID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "order_number", "cashier_id", "tax_id", "subtotal", "tax_amount", "total_amount", "payment_method", "status", "created_at", "updated_at"}).
			AddRow(orderID, "ORD-123", "cashier-1", &taxID, int64(20000), int64(2000), int64(22000), "CASH", "PAID", tm, tm))

	mock.ExpectQuery("SELECT id, order_id, product_id, quantity, unit_price, subtotal FROM order_items").
		WithArgs(orderID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "order_id", "product_id", "quantity", "unit_price", "subtotal"}).
			AddRow("item-1", orderID, "prod-1", 1, int64(15000), int64(20000)))

	mock.ExpectQuery("SELECT id, order_item_id, modifier_id, price FROM order_item_modifiers").
		WithArgs([]string{"item-1"}).
		WillReturnRows(pgxmock.NewRows([]string{"id", "order_item_id", "modifier_id", "price"}).
			AddRow("mod-item-1", "item-1", "mod-1", int64(5000)))

	order, err := repo.GetByID(context.Background(), orderID)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if order.ID != orderID {
		t.Errorf("expected order ID %v, got %v", orderID, order.ID)
	}
	if len(order.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(order.Items))
	}
	if len(order.Items[0].Modifiers) != 1 {
		t.Errorf("expected 1 modifier, got %d", len(order.Items[0].Modifiers))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	repo := NewRepository(mock)
	taxID := "tax-1"
	tm := time.Now()

	mock.ExpectQuery("SELECT id, order_number, cashier_id, tax_id, subtotal, tax_amount, total_amount, payment_method, status, created_at, updated_at FROM orders").
		WillReturnRows(pgxmock.NewRows([]string{"id", "order_number", "cashier_id", "tax_id", "subtotal", "tax_amount", "total_amount", "payment_method", "status", "created_at", "updated_at"}).
			AddRow("order-1", "ORD-123", "cashier-1", &taxID, int64(20000), int64(2000), int64(22000), "CASH", "PAID", tm, tm))

	orders, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(orders))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
