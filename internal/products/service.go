package products

import "context"

type CreateProductDTO struct {
	CategoryID  string
	SKU         string
	Name        string
	Description string
	Price       int64
	ImageURL    string
}

type UpdateProductDTO struct {
	CategoryID  string
	SKU         string
	Name        string
	Description string
	Price       int64
	ImageURL    string
}

type Service interface {
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context) ([]*Product, error)
	CreateProduct(ctx context.Context, dto CreateProductDTO) (*Product, error)
	UpdateProduct(ctx context.Context, id string, dto UpdateProductDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	DeleteProduct(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListProducts(ctx context.Context) ([]*Product, error) {
	return s.repo.List(ctx)
}

func (s *service) CreateProduct(ctx context.Context, dto CreateProductDTO) (*Product, error) {
	p := &Product{
		CategoryID:  dto.CategoryID,
		SKU:         dto.SKU,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		ImageURL:    dto.ImageURL,
		IsActive:    true,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) UpdateProduct(ctx context.Context, id string, dto UpdateProductDTO) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	p.CategoryID = dto.CategoryID
	p.SKU = dto.SKU
	p.Name = dto.Name
	p.Description = dto.Description
	p.Price = dto.Price
	p.ImageURL = dto.ImageURL

	return s.repo.Update(ctx, p)
}

func (s *service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *service) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
