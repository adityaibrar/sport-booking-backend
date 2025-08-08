package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Authentication
	SetupAuthRoutes(v1, db)
	// Admin and Venue
	SetupAdminRoutes(v1, db)
	// User and Booking
	SetupUserRoutes(v1, db)
}
