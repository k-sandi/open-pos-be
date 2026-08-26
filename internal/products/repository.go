package products

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Product represents the database entity
type Product struct {
	ID          string
	CategoryID  string
	SKU         string
	Name        string
	Description string
	Price       int64
	ImageURL    string
	IsActive    bool
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
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

func (r *repository) GetByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	err := r.db.QueryRow(ctx, `
		SELECT id, category_id, sku, name, description, price, image_url, is_active
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&p.ID, &p.CategoryID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.IsActive)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context) ([]*Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, category_id, sku, name, description, price, image_url, is_active
		FROM products
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.SKU, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.IsActive); err == nil {
			products = append(products, &p)
		}
	}
	return products, nil
}

func (r *repository) Create(ctx context.Context, p *Product) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO products (category_id, sku, name, description, price, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`, p.CategoryID, p.SKU, p.Name, p.Description, p.Price, p.ImageURL, p.IsActive).Scan(&p.ID)
}

func (r *repository) Update(ctx context.Context, p *Product) error {
	_, err := r.db.Exec(ctx, `
		UPDATE products 
		SET category_id=$1, sku=$2, name=$3, description=$4, price=$5, image_url=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL
	`, p.CategoryID, p.SKU, p.Name, p.Description, p.Price, p.ImageURL, p.IsActive, p.ID)
	return err
}

func (r *repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE products SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE products SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
