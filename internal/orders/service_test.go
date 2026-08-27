package orders

import (
	"context"
	"testing"
)

type mockRepo struct {
	getProductPrice  func(ctx context.Context, productID string) (int64, error)
	getModifierPrice func(ctx context.Context, modifierID string) (int64, error)
	getTaxRate       func(ctx context.Context, taxID string) (float64, error)
	createOrderTx    func(ctx context.Context, order *Order) error
	getByID          func(ctx context.Context, id string) (*Order, error)
	list             func(ctx context.Context) ([]*Order, error)
}

func (m *mockRepo) GetProductPrice(ctx context.Context, productID string) (int64, error) {
	return m.getProductPrice(ctx, productID)
}

func (m *mockRepo) GetModifierPrice(ctx context.Context, modifierID string) (int64, error) {
	return m.getModifierPrice(ctx, modifierID)
}

func (m *mockRepo) GetTaxRate(ctx context.Context, taxID string) (float64, error) {
	return m.getTaxRate(ctx, taxID)
}

func (m *mockRepo) CreateOrderTx(ctx context.Context, order *Order) error {
	return m.createOrderTx(ctx, order)
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*Order, error) {
	if m.getByID != nil {
		return m.getByID(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) List(ctx context.Context) ([]*Order, error) {
	if m.list != nil {
		return m.list(ctx)
	}
	return nil, nil
}

func TestService_CreateOrder(t *testing.T) {
	taxID := "tax-1"
	req := CreateOrderRequest{
		TaxID:         &taxID,
		PaymentMethod: "CASH",
		Items: []CreateOrderItemReq{
			{
				ProductID: "prod-1",
				Quantity:  2,
				Modifiers: []string{"mod-1"},
			},
		},
	}

	repo := &mockRepo{
		getProductPrice: func(ctx context.Context, productID string) (int64, error) {
			if productID == "prod-1" {
				return 20000, nil
			}
			return 0, nil
		},
		getModifierPrice: func(ctx context.Context, modifierID string) (int64, error) {
			if modifierID == "mod-1" {
				return 5000, nil
			}
			return 0, nil
		},
		getTaxRate: func(ctx context.Context, taxID string) (float64, error) {
			if taxID == "tax-1" {
				return 10.0, nil
			}
			return 0, nil
		},
		createOrderTx: func(ctx context.Context, order *Order) error {
			order.ID = "order-123"
			return nil
		},
	}

	svc := NewService(repo)

	order, err := svc.CreateOrder(context.Background(), "cashier-1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order == nil {
		t.Fatalf("expected order to be created, got nil")
	}

	// Subtotal: 2 * (20000 + 5000) = 50000
	if order.Subtotal != 50000 {
		t.Errorf("expected subtotal 50000, got %v", order.Subtotal)
	}

	// Tax: 10% of 50000 = 5000
	if order.TaxAmount != 5000 {
		t.Errorf("expected tax amount 5000, got %v", order.TaxAmount)
	}

	// Total: 55000
	if order.TotalAmount != 55000 {
		t.Errorf("expected total amount 55000, got %v", order.TotalAmount)
	}

	if order.CashierID != "cashier-1" {
		t.Errorf("expected cashier ID 'cashier-1', got %v", order.CashierID)
	}

	if order.ID != "order-123" {
		t.Errorf("expected order ID 'order-123', got %v", order.ID)
	}
}

func TestService_GetOrder(t *testing.T) {
	repo := &mockRepo{
		getByID: func(ctx context.Context, id string) (*Order, error) {
			if id == "order-123" {
				return &Order{ID: "order-123"}, nil
			}
			return nil, nil
		},
	}
	svc := NewService(repo)

	order, err := svc.GetOrder(context.Background(), "order-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order == nil || order.ID != "order-123" {
		t.Errorf("expected order ID 'order-123', got %v", order)
	}
}

func TestService_ListOrders(t *testing.T) {
	repo := &mockRepo{
		list: func(ctx context.Context) ([]*Order, error) {
			return []*Order{{ID: "order-123"}}, nil
		},
	}
	svc := NewService(repo)

	orders, err := svc.ListOrders(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(orders) != 1 || orders[0].ID != "order-123" {
		t.Errorf("expected 1 order with ID 'order-123', got %v", orders)
	}
}
