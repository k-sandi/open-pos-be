package variants

import (
	"context"
	"errors"
)

type ModifierGroupDTO struct {
	ProductID  string `json:"product_id"`
	Name       string `json:"name"`
	MinChoices int    `json:"min_choices"`
	MaxChoices int    `json:"max_choices"`
	IsActive   bool   `json:"is_active"`
}

type ModifierGroupService interface {
	ListByProduct(ctx context.Context, productID string) ([]*ModifierGroup, error)
	GetByID(ctx context.Context, id string) (*ModifierGroup, error)
	Create(ctx context.Context, dto ModifierGroupDTO) (*ModifierGroup, error)
	Update(ctx context.Context, id string, dto ModifierGroupDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	Delete(ctx context.Context, id string) error
}

type modifierGroupService struct {
	repo ModifierGroupRepository
}

func NewModifierGroupService(repo ModifierGroupRepository) ModifierGroupService {
	return &modifierGroupService{repo: repo}
}

func (s *modifierGroupService) ListByProduct(ctx context.Context, productID string) ([]*ModifierGroup, error) {
	return s.repo.ListByProduct(ctx, productID)
}

func (s *modifierGroupService) GetByID(ctx context.Context, id string) (*ModifierGroup, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *modifierGroupService) Create(ctx context.Context, dto ModifierGroupDTO) (*ModifierGroup, error) {
	mg := &ModifierGroup{
		ProductID:  dto.ProductID,
		Name:       dto.Name,
		MinChoices: dto.MinChoices,
		MaxChoices: dto.MaxChoices,
		IsActive:   dto.IsActive,
	}

	err := s.repo.Create(ctx, mg)
	if err != nil {
		return nil, err
	}

	return mg, nil
}

func (s *modifierGroupService) Update(ctx context.Context, id string, dto ModifierGroupDTO) error {
	mg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if mg == nil {
		return errors.New("modifier group not found")
	}

	mg.Name = dto.Name
	mg.MinChoices = dto.MinChoices
	mg.MaxChoices = dto.MaxChoices
	mg.IsActive = dto.IsActive

	return s.repo.Update(ctx, mg)
}

func (s *modifierGroupService) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *modifierGroupService) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
