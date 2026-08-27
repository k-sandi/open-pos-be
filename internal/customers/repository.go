package customers

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Customer represents the database entity
type Customer struct {
	ID            string
	Name          string
	Phone         string
	Email         *string
	LoyaltyPoints int
	IsActive      bool
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Customer, error)
	GetByPhone(ctx context.Context, phone string) (*Customer, error)
	Search(ctx context.Context, query string) ([]*Customer, error)
	List(ctx context.Context) ([]*Customer, error)
	Create(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	UpdatePoints(ctx context.Context, id string, pointDelta int) error
	SoftDelete(ctx context.Context, id string) error
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*Customer, error) {
	var c Customer
	err := r.db.QueryRow(ctx, `
		SELECT id, name, phone, email, loyalty_points, is_active
		FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.LoyaltyPoints, &c.IsActive)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) GetByPhone(ctx context.Context, phone string) (*Customer, error) {
	var c Customer
	err := r.db.QueryRow(ctx, `
		SELECT id, name, phone, email, loyalty_points, is_active
		FROM customers
		WHERE phone = $1 AND deleted_at IS NULL
	`, phone).Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.LoyaltyPoints, &c.IsActive)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) Search(ctx context.Context, query string) ([]*Customer, error) {
	searchPattern := "%" + query + "%"
	rows, err := r.db.Query(ctx, `
		SELECT id, name, phone, email, loyalty_points, is_active
		FROM customers
		WHERE (name ILIKE $1 OR phone ILIKE $1) AND deleted_at IS NULL
	`, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.LoyaltyPoints, &c.IsActive); err == nil {
			customers = append(customers, &c)
		}
	}
	return customers, nil
}

func (r *repository) List(ctx context.Context) ([]*Customer, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, phone, email, loyalty_points, is_active
		FROM customers
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Email, &c.LoyaltyPoints, &c.IsActive); err == nil {
			customers = append(customers, &c)
		}
	}
	return customers, nil
}

func (r *repository) Create(ctx context.Context, c *Customer) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO customers (name, phone, email, loyalty_points, is_active)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, c.Name, c.Phone, c.Email, c.LoyaltyPoints, c.IsActive).Scan(&c.ID)
}

func (r *repository) Update(ctx context.Context, c *Customer) error {
	_, err := r.db.Exec(ctx, `
		UPDATE customers 
		SET name=$1, phone=$2, email=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5 AND deleted_at IS NULL
	`, c.Name, c.Phone, c.Email, c.IsActive, c.ID)
	return err
}

func (r *repository) UpdatePoints(ctx context.Context, id string, pointDelta int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE customers 
		SET loyalty_points = loyalty_points + $1, updated_at = NOW() 
		WHERE id=$2 AND deleted_at IS NULL
	`, pointDelta, id)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE customers SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
