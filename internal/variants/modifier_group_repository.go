package variants

import (
	"context"
)

type ModifierGroup struct {
	ID         string `json:"id"`
	ProductID  string `json:"product_id"`
	Name       string `json:"name"`
	MinChoices int    `json:"min_choices"`
	MaxChoices int    `json:"max_choices"`
	IsActive   bool   `json:"is_active"`
}

type ModifierGroupRepository interface {
	ListByProduct(ctx context.Context, productID string) ([]*ModifierGroup, error)
	GetByID(ctx context.Context, id string) (*ModifierGroup, error)
	Create(ctx context.Context, mg *ModifierGroup) error
	Update(ctx context.Context, mg *ModifierGroup) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	SoftDelete(ctx context.Context, id string) error
}


type modifierGroupRepository struct {
	db DBTX
}

func NewModifierGroupRepository(db DBTX) ModifierGroupRepository {
	return &modifierGroupRepository{db: db}
}

func (r *modifierGroupRepository) ListByProduct(ctx context.Context, productID string) ([]*ModifierGroup, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, product_id, name, min_choices, max_choices, is_active
		FROM modifier_groups
		WHERE product_id = $1 AND deleted_at IS NULL
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*ModifierGroup
	for rows.Next() {
		var mg ModifierGroup
		if err := rows.Scan(&mg.ID, &mg.ProductID, &mg.Name, &mg.MinChoices, &mg.MaxChoices, &mg.IsActive); err == nil {
			groups = append(groups, &mg)
		}
	}
	return groups, nil
}

func (r *modifierGroupRepository) GetByID(ctx context.Context, id string) (*ModifierGroup, error) {
	var mg ModifierGroup
	err := r.db.QueryRow(ctx, `
		SELECT id, product_id, name, min_choices, max_choices, is_active
		FROM modifier_groups
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&mg.ID, &mg.ProductID, &mg.Name, &mg.MinChoices, &mg.MaxChoices, &mg.IsActive)
	if err != nil {
		return nil, err
	}
	return &mg, nil
}

func (r *modifierGroupRepository) Create(ctx context.Context, mg *ModifierGroup) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO modifier_groups (product_id, name, min_choices, max_choices, is_active)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, mg.ProductID, mg.Name, mg.MinChoices, mg.MaxChoices, mg.IsActive).Scan(&mg.ID)
}

func (r *modifierGroupRepository) Update(ctx context.Context, mg *ModifierGroup) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifier_groups 
		SET name=$1, min_choices=$2, max_choices=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5 AND deleted_at IS NULL
	`, mg.Name, mg.MinChoices, mg.MaxChoices, mg.IsActive, mg.ID)
	return err
}

func (r *modifierGroupRepository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifier_groups SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *modifierGroupRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE modifier_groups SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
