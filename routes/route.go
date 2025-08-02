package routes

import (
	"sport-booking-backend/controllers"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	venueController := controllers.NewVenueController(db)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Authentication
	v1.Post("/register", authController.Register)
	v1.Post("/login", authController.Login)

	admin := v1.Group("/admin", utils.AuthMiddleware, utils.CheckRole)
	// Venue
	venue := admin.Group("/venue")
	venue.Post("/", venueController.CreateVenue)
	venue.Put("/:id", venueController.UpdateVenue)
	venue.Get("/", venueController.GetListVenue)
	venue.Get("/:id", venueController.GetDetailVenue)
	venue.Delete("/:id", venueController.DeleteVenue)
}
