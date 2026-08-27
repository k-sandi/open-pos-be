package customers

import (
	"context"
	"testing"
)

type mockRepo struct {
	customer *Customer
	err      error
	pointsUpdated int
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*Customer, error) {
	return m.customer, m.err
}
func (m *mockRepo) GetByPhone(ctx context.Context, phone string) (*Customer, error) {
	return m.customer, m.err
}
func (m *mockRepo) Search(ctx context.Context, query string) ([]*Customer, error) {
	return nil, nil
}
func (m *mockRepo) List(ctx context.Context) ([]*Customer, error) {
	return nil, nil
}
func (m *mockRepo) Create(ctx context.Context, c *Customer) error {
	return nil
}
func (m *mockRepo) Update(ctx context.Context, c *Customer) error {
	return nil
}
func (m *mockRepo) UpdatePoints(ctx context.Context, id string, pointDelta int) error {
	m.pointsUpdated = pointDelta
	return nil
}
func (m *mockRepo) SoftDelete(ctx context.Context, id string) error {
	return nil
}

func TestService_UpdatePoints(t *testing.T) {
	repo := &mockRepo{
		customer: &Customer{ID: "1", LoyaltyPoints: 100},
	}
	svc := NewService(repo)

	// Test valid deduction
	err := svc.UpdatePoints(context.Background(), "1", -50)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if repo.pointsUpdated != -50 {
		t.Errorf("expected -50 points updated, got %d", repo.pointsUpdated)
	}

	// Test invalid deduction (drops below 0)
	err = svc.UpdatePoints(context.Background(), "1", -150)
	if err == nil || err.Error() != "insufficient loyalty points" {
		t.Errorf("expected 'insufficient loyalty points' error, got %v", err)
	}
}
