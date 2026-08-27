package customers

import (
	"context"
	"errors"
)

type CreateCustomerDTO struct {
	Name          string
	Phone         string
	Email         *string
	LoyaltyPoints int
}

type UpdateCustomerDTO struct {
	Name     string
	Phone    string
	Email    *string
	IsActive bool
}

type Service interface {
	GetCustomer(ctx context.Context, id string) (*Customer, error)
	SearchCustomers(ctx context.Context, query string) ([]*Customer, error)
	ListCustomers(ctx context.Context) ([]*Customer, error)
	CreateCustomer(ctx context.Context, dto CreateCustomerDTO) (*Customer, error)
	UpdateCustomer(ctx context.Context, id string, dto UpdateCustomerDTO) error
	UpdatePoints(ctx context.Context, id string, pointDelta int) error
	DeleteCustomer(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) SearchCustomers(ctx context.Context, query string) ([]*Customer, error) {
	if query == "" {
		return s.repo.List(ctx)
	}
	return s.repo.Search(ctx, query)
}

func (s *service) ListCustomers(ctx context.Context) ([]*Customer, error) {
	return s.repo.List(ctx)
}

func (s *service) CreateCustomer(ctx context.Context, dto CreateCustomerDTO) (*Customer, error) {
	customer := &Customer{
		Name:          dto.Name,
		Phone:         dto.Phone,
		Email:         dto.Email,
		LoyaltyPoints: dto.LoyaltyPoints,
		IsActive:      true,
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *service) UpdateCustomer(ctx context.Context, id string, dto UpdateCustomerDTO) error {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	customer.Name = dto.Name
	customer.Phone = dto.Phone
	customer.Email = dto.Email
	customer.IsActive = dto.IsActive

	return s.repo.Update(ctx, customer)
}

func (s *service) UpdatePoints(ctx context.Context, id string, pointDelta int) error {
	if pointDelta == 0 {
		return nil
	}

	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if customer.LoyaltyPoints+pointDelta < 0 {
		return errors.New("insufficient loyalty points")
	}

	return s.repo.UpdatePoints(ctx, id, pointDelta)
}

func (s *service) DeleteCustomer(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
