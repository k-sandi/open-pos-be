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

	// 5. Seed Modifiers
	log.Println("Seeding variants and modifiers...")
	var groupID string
	err = dbPool.QueryRow(ctx, `
		INSERT INTO modifier_groups (product_id, name, min_choices, max_choices, is_active)
		VALUES (
			(SELECT id FROM products WHERE sku = 'BEV-001'), 
			$1, 1, 1, true
		)
		ON CONFLICT DO NOTHING
		RETURNING id
	`, "Ice Level").Scan(&groupID)

	if err != nil && err.Error() != "no rows in result set" {
		log.Fatalf("Failed to seed modifier group: %v", err)
	}

	if groupID != "" {
		_, err = dbPool.Exec(ctx, `
			INSERT INTO modifiers (modifier_group_id, name, additional_price, is_active)
			VALUES 
				($1, 'Normal Ice', 0, true),
				($1, 'Less Ice', 0, true),
				($1, 'No Ice', 0, true)
		`, groupID)
		if err != nil {
			log.Fatalf("Failed to seed modifiers: %v", err)
		}
	}

	// 6. Seed Taxes
	log.Println("Seeding taxes...")
	_, err = dbPool.Exec(ctx, `
		INSERT INTO taxes (name, rate, is_active)
		VALUES ('PB1 10%', 10.00, true)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		log.Fatalf("Failed to seed taxes: %v", err)
	}
	
	// 7. Seed Customers
	log.Println("Seeding customers...")
	_, err = dbPool.Exec(ctx, `
		INSERT INTO customers (name, phone, email, loyalty_points, is_active)
		VALUES ('John Doe', '08123456789', 'john.doe@example.com', 100, true)
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		log.Fatalf("Failed to seed customers: %v", err)
	}

	log.Println("Database seeded successfully! Default Admin - Employee ID: admin01, PIN: 123456")
}
