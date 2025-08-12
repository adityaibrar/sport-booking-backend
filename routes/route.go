package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, client *snap.Client) {

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Authentication
	SetupAuthRoutes(v1, db)
	// Admin and Venue
	SetupAdminRoutes(v1, db, client)
	// User and Booking
	SetupUserRoutes(v1, db, client)
}
