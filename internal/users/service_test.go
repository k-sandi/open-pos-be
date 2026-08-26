package users_test

import (
	"context"
	"errors"
	"testing"

	"open-pos-be/internal/users"
)

type mockRepository struct {
	users.Repository
	mockGetByID func(ctx context.Context, id string) (*users.User, error)
	mockCreate  func(ctx context.Context, u *users.User) error
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	if m.mockGetByID != nil {
		return m.mockGetByID(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepository) Create(ctx context.Context, u *users.User) error {
	if m.mockCreate != nil {
		return m.mockCreate(ctx, u)
	}
	return errors.New("not implemented")
}

func TestService_GetProfile(t *testing.T) {
	repo := &mockRepository{
		mockGetByID: func(ctx context.Context, id string) (*users.User, error) {
			if id == "123" {
				return &users.User{ID: "123", Name: "John Doe"}, nil
			}
			return nil, errors.New("user not found")
		},
	}

	svc := users.NewService(repo)

	u, err := svc.GetProfile(context.Background(), "123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", u.Name)
	}
}

func TestService_CreateUser(t *testing.T) {
	repo := &mockRepository{
		mockCreate: func(ctx context.Context, u *users.User) error {
			u.ID = "new-id"
			return nil
		},
	}

	svc := users.NewService(repo)
	
	dto := users.CreateUserDTO{
		EmployeeID: "EMP99",
		PIN:        "123456",
		Name:       "Test User",
		Email:      "test@example.com",
		Phone:      "000000",
		RoleID:     "role1",
	}

	u, err := svc.CreateUser(context.Background(), dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if u.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", u.ID)
	}
	
	if len(u.PINHash) == 0 {
		t.Errorf("expected PIN to be hashed")
	}
}
