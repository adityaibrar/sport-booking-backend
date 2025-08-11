package routes

import (
	"sport-booking-backend/controllers"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupAdminRoutes configures all admin-related routes
func SetupAdminRoutes(version fiber.Router, db *gorm.DB) {
	adminController := controllers.NewAdminController(db)
	venueController := controllers.NewVenueController(db)
	bookingController := controllers.NewBookingController(db)

	// Admin routes group with authentication and role checking middleware
	admin := version.Group("/admin", utils.AuthMiddleware, utils.CheckRole)

	// Dashboard routes
	dashboard := admin.Group("/dashboard")
	dashboard.Get("/", adminController.GetDashboardAnalytics)      // GET /api/v1/admin/dashboard
	dashboard.Get("/stats", adminController.GetDashboardAnalytics) // GET /api/v1/admin/dashboard/stats

	// Venue management routes
	venues := admin.Group("/venues")
	venues.Post("/", venueController.CreateVenue)      // POST /api/v1/admin/venues
	venues.Get("/", venueController.GetListVenue)      // GET /api/v1/admin/venues
	venues.Get("/:id", venueController.GetDetailVenue) // GET /api/v1/admin/venues/:id
	venues.Put("/:id", venueController.UpdateVenue)    // PUT /api/v1/admin/venues/:id
	venues.Delete("/:id", venueController.DeleteVenue) // DELETE /api/v1/admin/venues/:id

	// Booking management routes
	bookings := admin.Group("/bookings")
	bookings.Get("/", bookingController.GetUserBookings)               // GET /api/v1/admin/bookings
	bookings.Put("/:id/status", bookingController.UpdateBookingStatus) // PUT /api/v1/admin/bookings/:id/status
	bookings.Get("/:id", bookingController.GetUserBookings)            // GET /api/v1/admin/bookings/:id (individual booking details)

	// User management routes (if needed in the future)
	// users := admin.Group("/users")
	// users.Get("/", userController.GetAllUsers)
	// users.Get("/:id", userController.GetUserDetails)
	// users.Put("/:id/status", userController.UpdateUserStatus)
}
