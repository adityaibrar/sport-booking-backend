package controllers

import (
	"fmt"
	"sport-booking-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BookingController struct {
	DB *gorm.DB
}

func NewBookingController(db *gorm.DB) *BookingController {
	return &BookingController{DB: db}
}

func (bc *BookingController) CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var request models.BookingRequest
	if err := c.BodyParser(&request); err != nil {
		fmt.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var venue models.Venue
	if err := bc.DB.First(&venue, request.VenueID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Venue not found",
		})
	}

	resultChan := make(chan error)
	go func() {
		var count int64
		err := bc.DB.Model(&models.Booking{}).
			Where(
				"venue_id = ? AND status IN (?, ?)", request.VenueID, "pending", "confirmed",
			).
			Where("start_time <= ? AND DATE_ADD(start_time, INTERVAL duration HOUR) >= ?", request.StartTime, request.StartTime).
			Count(&count).Error

		if err != nil {
			resultChan <- err
			return
		}
		if count > 0 {
			resultChan <- fiber.NewError(fiber.StatusConflict, "Venue not available at selected time")
			return
		}
		resultChan <- nil
	}()

	if err := <-resultChan; err != nil {
		return err
	}

	totalPrice := float64(request.Duration) * venue.PricePerHour

	booking := models.Booking{
		UserID:     userID,
		VenueID:    request.VenueID,
		StartTime:  request.StartTime,
		Duration:   request.Duration,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	// Check for overlapping bookings for the same venue
	var overlappingCount int64
	// Convert StartTime to time.Time for comparison using format "2006-01-02 15:04:05.999"
	startTime, err := time.Parse("2006-01-02 15:04:05.999", request.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start time format. Expected format: YYYY-MM-DD HH:mm:ss.SSS",
		})
	}

	endTime := startTime.Add(time.Duration(request.Duration) * time.Hour)
	endTimeStr := endTime.Format("2006-01-02 15:04:05.999")
	if err := bc.DB.Model(&models.Booking{}).
		Where("venue_id = ? AND status IN (?, ?)", request.VenueID, "pending", "confirmed").
		Where("? < DATE_ADD(start_time, INTERVAL duration HOUR) AND start_time < ?", request.StartTime, endTimeStr).
		Count(&overlappingCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check booking conflicts",
		})
	}
	if overlappingCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Venue already booked for the selected time slot",
		})
	}

	if err := bc.DB.Create(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booking",
		})
	}

	// Ambil data user dan venue untuk response
	var user models.User
	if err := bc.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user data",
		})
	}

	// Data venue sudah diambil sebelumnya (variable 'venue')
	response := fiber.Map{
		"message": "Booking created successfully",
		"data_booking": fiber.Map{
			"user":  user,
			"venue": venue,
			"booking": fiber.Map{
				"id":          booking.ID,
				"start_time":  booking.StartTime,
				"duration":    booking.Duration,
				"total_price": booking.TotalPrice,
				"status":      booking.Status,
			},
		},
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}
