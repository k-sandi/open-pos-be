package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"log/slog"
	
	"open-pos-be/internal/auth"
	"open-pos-be/internal/categories"
	customMiddleware "open-pos-be/internal/middleware"
	"open-pos-be/internal/products"
	"open-pos-be/internal/users"
	"open-pos-be/internal/variants"
)

// Setup configures and returns the main chi.Router with all mounted endpoints and middlewares.
func Setup(dbPool *pgxpool.Pool, logger *slog.Logger) chi.Router {
	// Initialize Layers
	userRepo := users.NewRepository(dbPool)
	userService := users.NewService(userRepo)
	userHandler := users.NewHandler(userService)

	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	categoryRepo := categories.NewRepository(dbPool)
	categoryService := categories.NewService(categoryRepo)
	categoryHandler := categories.NewHandler(categoryService)

	productRepo := products.NewRepository(dbPool)
	productService := products.NewService(productRepo)
	productHandler := products.NewHandler(productService)

	modifierGroupRepo := variants.NewModifierGroupRepository(dbPool)
	modifierGroupService := variants.NewModifierGroupService(modifierGroupRepo)
	modifierGroupHandler := variants.NewModifierGroupHandler(modifierGroupService)

	modifierRepo := variants.NewModifierRepository(dbPool)
	modifierService := variants.NewModifierService(modifierRepo)
	modifierHandler := variants.NewModifierHandler(modifierService)

	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(customMiddleware.StructuredLogger(logger))
	r.Use(middleware.Recoverer)

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"), // Using relative path is better than hardcoded localhost
		))

		r.Mount("/api/v1/auth", authHandler.Routes())
	})

	// Protected Routes (requires x-api-key or Bearer token)
	r.Group(func(r chi.Router) {
		r.Use(customMiddleware.Auth)

		// Mount Users API Routes
		r.Mount("/api/v1/users", userHandler.Routes())
		
		// Mount Categories API Routes
		r.Mount("/api/v1/categories", categoryHandler.Routes())

		// Mount Products API Routes
		r.Mount("/api/v1/products", productHandler.Routes())

		// Mount Variants API Routes
		r.Mount("/api/v1/modifier-groups", modifierGroupHandler.Routes())
		r.Mount("/api/v1/modifiers", modifierHandler.Routes())
	})

	return r
}
