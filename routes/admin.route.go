package routes

import (
	"sport-booking-backend/controllers"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAdminRoutes(version fiber.Router, db *gorm.DB) {
	adminController := controllers.NewAdminController(db)
	venueController := controllers.NewVenueController(db)

	admin := version.Group("/admin", utils.AuthMiddleware, utils.CheckRole)
	dashboard := admin.Group("/dashboard")
	dashboard.Get("/", adminController.GetDashboardAnalytics)

	venue := admin.Group("/venue")
	venue.Post("/", venueController.CreateVenue)
	venue.Put("/:id", venueController.UpdateVenue)
	venue.Get("/", venueController.GetListVenue)
	venue.Get("/:id", venueController.GetDetailVenue)
	venue.Delete("/:id", venueController.DeleteVenue)
}