package main

import (
	"log"
	"sport-booking-backend/config"
	"sport-booking-backend/middleware"
	"sport-booking-backend/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create Fiber app with configuration
	app := fiber.New(fiber.Config{
		AppName:      "Sport Booking API",
		ServerHeader: "Sport Booking Backend",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle Fiber errors
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Log error for debugging
			if cfg.Debug {
				log.Printf("Error: %v", err)
			}

			// Return JSON error response
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": "Internal Server Error",
				"error":   err.Error(),
			})
		},
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Database connection with proper configuration
	gormConfig := &gorm.Config{}
	if cfg.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Auto-migrate database schema
	if err := config.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "ok",
			"service":     "Sport Booking API",
			"version":     "1.0.0",
			"environment": cfg.Environment,
		})
	})

	// API documentation endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Sport Booking API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
			"health":  "/health",
			"endpoints": fiber.Map{
				"auth":  "/api/v1/auth",
				"user":  "/api/v1/user",
				"admin": "/api/v1/admin",
			},
		})
	})

	// Setup API routes
	routes.SetupRoutes(app, db)

	// Setup error handler for 404
	middleware.SetupErrorHandler(app)

	// Start server
	address := cfg.Host + ":" + cfg.Port
	log.Printf("🚀 Server starting on %s in %s mode", address, cfg.Environment)

	if err := app.Listen(address); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
