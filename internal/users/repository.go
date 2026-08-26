package users

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// User represents the database entity
type User struct {
	ID         string
	EmployeeID string
	PINHash    string
	Name       string
	Email      string
	Phone      string
	RoleID     string
	RoleName   string
	IsActive   bool
}

type Repository interface {
	GetByEmployeeID(ctx context.Context, employeeID string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	UpdateStatus(ctx context.Context, id string, isActive bool) error
	SoftDelete(ctx context.Context, id string) error
}

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) GetByEmployeeID(ctx context.Context, employeeID string) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.employee_id, u.pin_hash, u.name, u.email, u.phone, u.role_id, r.name as role_name, u.is_active
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.employee_id = $1 AND u.deleted_at IS NULL
	`, employeeID).Scan(&u.ID, &u.EmployeeID, &u.PINHash, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RoleName, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.employee_id, u.pin_hash, u.name, u.email, u.phone, u.role_id, r.name as role_name, u.is_active
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`, id).Scan(&u.ID, &u.EmployeeID, &u.PINHash, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RoleName, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) List(ctx context.Context) ([]*User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.employee_id, u.pin_hash, u.name, u.email, u.phone, u.role_id, r.name as role_name, u.is_active
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.EmployeeID, &u.PINHash, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RoleName, &u.IsActive); err == nil {
			users = append(users, &u)
		}
	}
	return users, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO users (employee_id, pin_hash, name, email, phone, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`, u.EmployeeID, u.PINHash, u.Name, u.Email, u.Phone, u.RoleID, u.IsActive).Scan(&u.ID)
}

func (r *repository) Update(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users 
		SET employee_id=$1, pin_hash=$2, name=$3, email=$4, phone=$5, role_id=$6, is_active=$7, updated_at=NOW()
		WHERE id=$8 AND deleted_at IS NULL
	`, u.EmployeeID, u.PINHash, u.Name, u.Email, u.Phone, u.RoleID, u.IsActive, u.ID)
	return err
}

func (r *repository) UpdateStatus(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET is_active=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL
	`, isActive, id)
	return err
}

func (r *repository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET deleted_at=NOW() WHERE id=$1
	`, id)
	return err
}
