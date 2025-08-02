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
	bookingController := controllers.NewBookingController(db)
	adminController := controllers.NewAdminController(db)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Authentication
	v1.Post("/register", authController.Register)
	v1.Post("/login", authController.Login)

	admin := v1.Group("/admin", utils.AuthMiddleware, utils.CheckRole)
	dashboard := admin.Group("/dashboard")
	dashboard.Get("/", adminController.GetDashboardAnalytics)
	// Venue
	venue := admin.Group("/venue")
	venue.Post("/", venueController.CreateVenue)
	venue.Put("/:id", venueController.UpdateVenue)
	venue.Get("/", venueController.GetListVenue)
	venue.Get("/:id", venueController.GetDetailVenue)
	venue.Delete("/:id", venueController.DeleteVenue)

	user := v1.Group("/user", utils.AuthMiddleware)

	user.Get("/venues", venueController.GetListVenue)
	user.Get("/venues/:id", venueController.GetDetailVenue)
	// Booking
	booking := user.Group("/booking")
	booking.Post("/", bookingController.CreateBooking)
}
