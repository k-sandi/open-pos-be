package orders

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Order struct {
	ID            string
	OrderNumber   string
	CashierID     string
	TaxID         *string
	Subtotal      int64
	TaxAmount     int64
	TotalAmount   int64
	PaymentMethod string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time

	Items []OrderItem
}

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	UnitPrice int64
	Subtotal  int64

	Modifiers []OrderItemModifier
}

type OrderItemModifier struct {
	ID          string
	OrderItemID string
	ModifierID  string
	Price       int64
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

type Repository interface {
	GetProductPrice(ctx context.Context, productID string) (int64, error)
	GetModifierPrice(ctx context.Context, modifierID string) (int64, error)
	GetTaxRate(ctx context.Context, taxID string) (float64, error)
	CreateOrderTx(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	List(ctx context.Context) ([]*Order, error)
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetProductPrice(ctx context.Context, productID string) (int64, error) {
	var price int64
	err := r.db.QueryRow(ctx, "SELECT price FROM products WHERE id = $1 AND is_active = true AND deleted_at IS NULL", productID).Scan(&price)
	return price, err
}

func (r *repository) GetModifierPrice(ctx context.Context, modifierID string) (int64, error) {
	var price int64
	err := r.db.QueryRow(ctx, "SELECT additional_price FROM modifiers WHERE id = $1 AND is_active = true AND deleted_at IS NULL", modifierID).Scan(&price)
	return price, err
}

func (r *repository) GetTaxRate(ctx context.Context, taxID string) (float64, error) {
	var rate float64
	err := r.db.QueryRow(ctx, "SELECT rate FROM taxes WHERE id = $1 AND is_active = true AND deleted_at IS NULL", taxID).Scan(&rate)
	return rate, err
}

func (r *repository) CreateOrderTx(ctx context.Context, order *Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Create order
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (order_number, cashier_id, tax_id, subtotal, tax_amount, total_amount, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, order.OrderNumber, order.CashierID, order.TaxID, order.Subtotal, order.TaxAmount, order.TotalAmount, order.PaymentMethod, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	for i := range order.Items {
		item := &order.Items[i]
		item.OrderID = order.ID

		err = tx.QueryRow(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price, subtotal)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal).Scan(&item.ID)
		if err != nil {
			return err
		}

		for j := range item.Modifiers {
			modifier := &item.Modifiers[j]
			modifier.OrderItemID = item.ID

			err = tx.QueryRow(ctx, `
				INSERT INTO order_item_modifiers (order_item_id, modifier_id, price)
				VALUES ($1, $2, $3)
				RETURNING id
			`, modifier.OrderItemID, modifier.ModifierID, modifier.Price).Scan(&modifier.ID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) GetByID(ctx context.Context, id string) (*Order, error) {
	var order Order
	err := r.db.QueryRow(ctx, `
		SELECT id, order_number, cashier_id, tax_id, subtotal, tax_amount, total_amount, payment_method, status, created_at, updated_at
		FROM orders
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&order.ID, &order.OrderNumber, &order.CashierID, &order.TaxID, &order.Subtotal, &order.TaxAmount, &order.TotalAmount, &order.PaymentMethod, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, order_id, product_id, quantity, unit_price, subtotal
		FROM order_items
		WHERE order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) > 0 {
		var itemIDs []string
		itemMap := make(map[string]*OrderItem)
		for i := range items {
			itemIDs = append(itemIDs, items[i].ID)
			itemMap[items[i].ID] = &items[i]
		}

		modRows, err := r.db.Query(ctx, `
			SELECT id, order_item_id, modifier_id, price
			FROM order_item_modifiers
			WHERE order_item_id = ANY($1)
		`, itemIDs)
		if err == nil {
			defer modRows.Close()
			for modRows.Next() {
				var mod OrderItemModifier
				if err := modRows.Scan(&mod.ID, &mod.OrderItemID, &mod.ModifierID, &mod.Price); err != nil {
					return nil, err
				}
				if item, ok := itemMap[mod.OrderItemID]; ok {
					item.Modifiers = append(item.Modifiers, mod)
				}
			}
			if err := modRows.Err(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	order.Items = items

	return &order, nil
}

func (r *repository) List(ctx context.Context) ([]*Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, order_number, cashier_id, tax_id, subtotal, tax_amount, total_amount, payment_method, status, created_at, updated_at
		FROM orders
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		order := &Order{}
		if err := rows.Scan(&order.ID, &order.OrderNumber, &order.CashierID, &order.TaxID, &order.Subtotal, &order.TaxAmount, &order.TotalAmount, &order.PaymentMethod, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
