package controllers

import (
	"sport-booking-backend/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

func (ac *AdminController) GetDashboardAnalytics(c *fiber.Ctx) error {
	var stats struct {
		TotalRevenue    float64 `json:"total_revenue"`
		TotalBookings   int64   `json:"total_bookings"`
		PendingBookings int64   `json:"pending_bookings"`
	}

	ac.DB.Model(&models.Booking{}).Where("status = ?", "completed").Select("SUM(total_price) as total_revenue").Scan(&stats)

	ac.DB.Model(&models.Booking{}).Count(&stats.TotalBookings)

	ac.DB.Model(&models.Booking{}).Where("status = ? ", "pending").Count(&stats.PendingBookings)

	return c.Status(fiber.StatusOK).JSON(stats)
}
