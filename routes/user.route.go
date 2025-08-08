package routes

import (
	"sport-booking-backend/controllers"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUserRoutes(version fiber.Router, db *gorm.DB) {
	bookingController := controllers.NewBookingController(db)
	venueController := controllers.NewVenueController(db)

	user := version.Group("/user", utils.AuthMiddleware)
	
	venue := user.Group("/venues")
	venue.Get("/", venueController.GetListVenue)
	venue.Get("/:id", venueController.GetDetailVenue)
	venue.Post("/availability", venueController.CheckAvailability)

	booking := user.Group("/booking")
	booking.Post("/", bookingController.CreateBooking)
	booking.Post("/cancel/:id", bookingController.CancelBooking)
}