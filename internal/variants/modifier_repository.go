package variants

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is an interface for pgx.Pool and pgx.Tx
type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

// Modifier represents the database entity
type Modifier struct {
	ID              string
	ModifierGroupID string
	Name            string
	AdditionalPrice int64
	IsActive        bool
}

type ModifierRepository interface {
	ListByGroup(ctx context.Context, groupID string) ([]*Modifier, error)
	GetByID(ctx context.Context, id string) (*Modifier, error)
	Create(ctx context.Context, m *Modifier) error
	Update(ctx context.Context, m *Modifier) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	SoftDelete(ctx context.Context, id string) error
}

type modifierRepository struct {
	db DBTX
}

func NewModifierRepository(db DBTX) ModifierRepository {
	return &modifierRepository{db: db}
}

func (r *modifierRepository) ListByGroup(ctx context.Context, groupID string) ([]*Modifier, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, modifier_group_id, name, additional_price, is_active
		FROM modifiers
		WHERE modifier_group_id = $1 AND deleted_at IS NULL
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modifiers []*Modifier
	for rows.Next() {
		var m Modifier
		if err := rows.Scan(&m.ID, &m.ModifierGroupID, &m.Name, &m.AdditionalPrice, &m.IsActive); err == nil {
			modifiers = append(modifiers, &m)
		}
	}
	return modifiers, nil
}

func (r *modifierRepository) GetByID(ctx context.Context, id string) (*Modifier, error) {
	var m Modifier
	err := r.db.QueryRow(ctx, `
		SELECT id, modifier_group_id, name, additional_price, is_active
		FROM modifiers
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&m.ID, &m.ModifierGroupID, &m.Name, &m.AdditionalPrice, &m.IsActive)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *modifierRepository) Create(ctx context.Context, m *Modifier) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO modifiers (modifier_group_id, name, additional_price, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, m.ModifierGroupID, m.Name, m.AdditionalPrice, m.IsActive).Scan(&m.ID)
}

func (r *modifierRepository) Update(ctx context.Context, m *Modifier) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifiers 
		SET modifier_group_id=$1, name=$2, additional_price=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5 AND deleted_at IS NULL
	`, m.ModifierGroupID, m.Name, m.AdditionalPrice, m.IsActive, m.ID)
	return err
}

func (r *modifierRepository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifiers SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *modifierRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifiers SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
