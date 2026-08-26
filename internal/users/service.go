package users

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserDTO struct {
	EmployeeID string 
	PIN        string 
	Name       string 
	Email      string 
	Phone      string 
	RoleID     string 
}

type UpdateUserDTO struct {
	EmployeeID string 
	PIN        string 
	Name       string 
	Email      string 
	Phone      string 
	RoleID     string 
}

type Service interface {
	GetProfile(ctx context.Context, userID string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, dto CreateUserDTO) (*User, error)
	UpdateUser(ctx context.Context, id string, dto UpdateUserDTO) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	DeleteUser(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *service) ListUsers(ctx context.Context) ([]*User, error) {
	return s.repo.List(ctx)
}

func (s *service) CreateUser(ctx context.Context, dto CreateUserDTO) (*User, error) {
	pinHash, err := bcrypt.GenerateFromPassword([]byte(dto.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash pin")
	}

	user := &User{
		EmployeeID: dto.EmployeeID,
		PINHash:    string(pinHash),
		Name:       dto.Name,
		Email:      dto.Email,
		Phone:      dto.Phone,
		RoleID:     dto.RoleID,
		IsActive:   true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) UpdateUser(ctx context.Context, id string, dto UpdateUserDTO) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.EmployeeID = dto.EmployeeID
	user.Name = dto.Name
	user.Email = dto.Email
	user.Phone = dto.Phone
	user.RoleID = dto.RoleID

	if dto.PIN != "" {
		pinHash, err := bcrypt.GenerateFromPassword([]byte(dto.PIN), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash pin")
		}
		user.PINHash = string(pinHash)
	}

	return s.repo.Update(ctx, user)
}

func (s *service) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	return s.repo.UpdateStatus(ctx, id, isActive)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
