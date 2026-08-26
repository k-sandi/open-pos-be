package categories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Category struct {
	ID          string
	Name        string
	Description string
	IsActive    bool
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
	Create(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
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

func (r *repository) GetByID(ctx context.Context, id string) (*Category, error) {
	var c Category
	err := r.db.QueryRow(ctx, `
		SELECT id, name, description, is_active
		FROM categories
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&c.ID, &c.Name, &c.Description, &c.IsActive)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) List(ctx context.Context) ([]*Category, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, description, is_active
		FROM categories
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.IsActive); err == nil {
			categories = append(categories, &c)
		}
	}
	return categories, nil
}

func (r *repository) Create(ctx context.Context, c *Category) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO categories (name, description, is_active)
		VALUES ($1, $2, $3) RETURNING id
	`, c.Name, c.Description, c.IsActive).Scan(&c.ID)
}

func (r *repository) Update(ctx context.Context, c *Category) error {
	_, err := r.db.Exec(ctx, `
		UPDATE categories 
		SET name=$1, description=$2, is_active=$3, updated_at=NOW()
		WHERE id=$4 AND deleted_at IS NULL
	`, c.Name, c.Description, c.IsActive, c.ID)
	return err
}

func (r *repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE categories SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE categories SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
