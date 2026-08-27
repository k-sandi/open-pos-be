package taxes

import (
	"context"
)

type CreateTaxDTO struct {
	Name     string
	Rate     float64
	IsActive bool
}

type UpdateTaxDTO struct {
	Name     string
	Rate     float64
	IsActive bool
}

type Service interface {
	GetTax(ctx context.Context, id string) (*Tax, error)
	ListTaxes(ctx context.Context) ([]*Tax, error)
	CreateTax(ctx context.Context, dto CreateTaxDTO) (*Tax, error)
	UpdateTax(ctx context.Context, id string, dto UpdateTaxDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	DeleteTax(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetTax(ctx context.Context, id string) (*Tax, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListTaxes(ctx context.Context) ([]*Tax, error) {
	return s.repo.List(ctx)
}

func (s *service) CreateTax(ctx context.Context, dto CreateTaxDTO) (*Tax, error) {
	tax := &Tax{
		Name:     dto.Name,
		Rate:     dto.Rate,
		IsActive: dto.IsActive,
	}
	if err := s.repo.Create(ctx, tax); err != nil {
		return nil, err
	}
	return tax, nil
}

func (s *service) UpdateTax(ctx context.Context, id string, dto UpdateTaxDTO) error {
	tax := &Tax{
		ID:       id,
		Name:     dto.Name,
		Rate:     dto.Rate,
		IsActive: dto.IsActive,
	}
	return s.repo.Update(ctx, tax)
}

func (s *service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *service) DeleteTax(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
