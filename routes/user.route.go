package routes

import (
	"sport-booking-backend/controllers"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(version fiber.Router, db *gorm.DB) {
	bookingController := controllers.NewBookingController(db)
	venueController := controllers.NewVenueController(db)
	reviewController := controllers.NewReviewController(db)

	// User routes group with authentication middleware
	user := version.Group("/user", utils.AuthMiddleware)

	// Venue routes for users
	venues := user.Group("/venues")
	venues.Get("/", venueController.GetListVenue)                   // GET /api/v1/user/venues
	venues.Get("/:id", venueController.GetDetailVenue)              // GET /api/v1/user/venues/:id
	venues.Post("/availability", venueController.CheckAvailability) // POST /api/v1/user/venues/availability

	// Booking routes for users
	bookings := user.Group("/bookings")
	bookings.Post("/", bookingController.CreateBooking)                // POST /api/v1/user/bookings
	bookings.Get("/history/:id", bookingController.HistoryBookingUser) // GET /api/v1/user/bookings/history/:id
	bookings.Post("/:id/cancel", bookingController.CancelBooking)      // POST /api/v1/user/bookings/:id/cancel

	// Review routes for users
	reviews := user.Group("/reviews")
	reviews.Post("/", reviewController.CreateReview) // POST /api/v1/user/reviews
	// reviews.Get("/history/:id", reviewController.HistoryReviewUser)    // GET /api/v1/user/reviews/history/:id
	reviews.Get("/", reviewController.GetReview)          // GET /api/v1/user/reviews
	reviews.Get("/:id", reviewController.GetReview)       // GET /api/v1/user/reviews/:id
	reviews.Delete("/:id", reviewController.DeleteReview) // DELETE /api/v1/user/reviews/:id

}
