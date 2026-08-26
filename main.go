package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "open-pos-be/docs"

	"open-pos-be/internal/router"
)

// @title Open POS API
// @version 1.0
// @description Backend API for the Open POS application.
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, relying on environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL environment variable is not set")
		os.Exit(1)
	}

	// Initialize PostgreSQL connection pool
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Verify database connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("Unable to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to PostgreSQL database successfully!")

	// Setup router with all routes and middlewares
	r := router.Setup(dbPool, logger)

	logger.Info("Starting server...",
		slog.String("port", port),
		slog.String("swagger_ui", "http://localhost:"+port+"/swagger/index.html"),
		slog.String("api_base", "http://localhost:"+port+"/api/v1"),
	)
	
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
