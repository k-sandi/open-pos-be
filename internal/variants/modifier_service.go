package variants

import (
	"context"
	"errors"
)

type ModifierDTO struct {
	ModifierGroupID string `json:"modifier_group_id"`
	Name            string `json:"name"`
	AdditionalPrice int64  `json:"additional_price"`
}

type ModifierService interface {
	ListModifiers(ctx context.Context, groupID string) ([]*Modifier, error)
	GetModifier(ctx context.Context, id string) (*Modifier, error)
	CreateModifier(ctx context.Context, dto ModifierDTO) (string, error)
	UpdateModifier(ctx context.Context, id string, dto ModifierDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	DeleteModifier(ctx context.Context, id string) error
}

type modifierService struct {
	repo ModifierRepository
}

func NewModifierService(repo ModifierRepository) ModifierService {
	return &modifierService{repo: repo}
}

func (s *modifierService) ListModifiers(ctx context.Context, groupID string) ([]*Modifier, error) {
	return s.repo.ListByGroup(ctx, groupID)
}

func (s *modifierService) GetModifier(ctx context.Context, id string) (*Modifier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *modifierService) CreateModifier(ctx context.Context, dto ModifierDTO) (string, error) {
	if dto.Name == "" || dto.ModifierGroupID == "" {
		return "", errors.New("name and modifier_group_id are required")
	}
	m := &Modifier{
		ModifierGroupID: dto.ModifierGroupID,
		Name:            dto.Name,
		AdditionalPrice: dto.AdditionalPrice,
		IsActive:        true,
	}
	err := s.repo.Create(ctx, m)
	if err != nil {
		return "", err
	}
	return m.ID, nil
}

func (s *modifierService) UpdateModifier(ctx context.Context, id string, dto ModifierDTO) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	m.ModifierGroupID = dto.ModifierGroupID
	m.Name = dto.Name
	m.AdditionalPrice = dto.AdditionalPrice

	return s.repo.Update(ctx, m)
}

func (s *modifierService) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *modifierService) DeleteModifier(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
