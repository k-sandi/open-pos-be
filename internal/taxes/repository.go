package taxes

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Tax struct {
	ID       string
	Name     string
	Rate     float64
	IsActive bool
}

type Repository interface {
	GetByID(ctx context.Context, id string) (*Tax, error)
	List(ctx context.Context) ([]*Tax, error)
	Create(ctx context.Context, t *Tax) error
	Update(ctx context.Context, t *Tax) error
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

func (r *repository) GetByID(ctx context.Context, id string) (*Tax, error) {
	var t Tax
	err := r.db.QueryRow(ctx, `
		SELECT id, name, rate, is_active
		FROM taxes
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&t.ID, &t.Name, &t.Rate, &t.IsActive)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *repository) List(ctx context.Context) ([]*Tax, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, rate, is_active
		FROM taxes
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var taxes []*Tax
	for rows.Next() {
		var t Tax
		if err := rows.Scan(&t.ID, &t.Name, &t.Rate, &t.IsActive); err == nil {
			taxes = append(taxes, &t)
		}
	}
	return taxes, nil
}

func (r *repository) Create(ctx context.Context, t *Tax) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO taxes (name, rate, is_active)
		VALUES ($1, $2, $3) RETURNING id
	`, t.Name, t.Rate, t.IsActive).Scan(&t.ID)
}

func (r *repository) Update(ctx context.Context, t *Tax) error {
	_, err := r.db.Exec(ctx, `
		UPDATE taxes 
		SET name=$1, rate=$2, is_active=$3, updated_at=NOW()
		WHERE id=$4 AND deleted_at IS NULL
	`, t.Name, t.Rate, t.IsActive, t.ID)
	return err
}

func (r *repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE taxes SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE taxes SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
