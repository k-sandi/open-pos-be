package categories

import (
	"context"
)

type CreateCategoryDTO struct {
	Name        string
	Description string
}

type UpdateCategoryDTO struct {
	Name        string
	Description string
}

type Service interface {
	GetCategory(ctx context.Context, id string) (*Category, error)
	ListCategories(ctx context.Context) ([]*Category, error)
	CreateCategory(ctx context.Context, dto CreateCategoryDTO) (*Category, error)
	UpdateCategory(ctx context.Context, id string, dto UpdateCategoryDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	DeleteCategory(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCategory(ctx context.Context, id string) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListCategories(ctx context.Context) ([]*Category, error) {
	return s.repo.List(ctx)
}

func (s *service) CreateCategory(ctx context.Context, dto CreateCategoryDTO) (*Category, error) {
	category := &Category{
		Name:        dto.Name,
		Description: dto.Description,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) UpdateCategory(ctx context.Context, id string, dto UpdateCategoryDTO) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	category.Name = dto.Name
	category.Description = dto.Description

	return s.repo.Update(ctx, category)
}

func (s *service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *service) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
