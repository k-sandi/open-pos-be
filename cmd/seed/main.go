package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbPool.Close()

	// 1. Seed Roles
	roles := []string{"Admin", "Manager", "Cashier"}
	log.Println("Seeding roles...")
	for _, role := range roles {
		_, err := dbPool.Exec(ctx, `
			INSERT INTO roles (name, is_active) 
			VALUES ($1, true) 
			ON CONFLICT (name) DO NOTHING
		`, role)
		if err != nil {
			log.Fatalf("Failed to seed role %s: %v", role, err)
		}
	}
	log.Println("Roles seeded successfully.")

	// Get Admin Role ID
	var adminRoleID string
	err = dbPool.QueryRow(ctx, "SELECT id FROM roles WHERE name = 'Admin'").Scan(&adminRoleID)
	if err != nil {
		log.Fatalf("Failed to get Admin role ID: %v", err)
	}

	// 2. Seed Default Admin User
	employeeID := "admin01"
	plainPIN := "123456"
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(plainPIN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash pin: %v", err)
	}

	log.Println("Seeding default admin user...")
	_, err = dbPool.Exec(ctx, `
		INSERT INTO users (employee_id, pin_hash, name, email, phone, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
		ON CONFLICT (employee_id) DO NOTHING
	`, employeeID, string(hashedPIN), "Super Admin", "admin@openpos.local", "+1234567890", adminRoleID)

	if err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// 3. Seed Categories
	log.Println("Seeding categories...")
	var categoryID string
	err = dbPool.QueryRow(ctx, `
		INSERT INTO categories (name, description, is_active)
		VALUES ($1, $2, true)
		ON CONFLICT (name) DO UPDATE SET is_active = true
		RETURNING id
	`, "Beverages", "All drinks").Scan(&categoryID)
	if err != nil {
		log.Fatalf("Failed to seed category: %v", err)
	}

	// 4. Seed Products
	log.Println("Seeding products...")
	_, err = dbPool.Exec(ctx, `
		INSERT INTO products (category_id, sku, name, description, price, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		ON CONFLICT (sku) DO NOTHING
	`, categoryID, "BEV-001", "Espresso", "Strong coffee", 20000)
	if err != nil {
		log.Fatalf("Failed to seed product: %v", err)
	}

	log.Println("Database seeded successfully! Default Admin - Employee ID: admin01, PIN: 123456")
}
