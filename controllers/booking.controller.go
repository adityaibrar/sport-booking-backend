package controllers

import (
	"fmt"
	"sport-booking-backend/models"

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
	userID := c.Locals("user_id").(float64)

	var request models.BookingRequest
	if err := c.BodyParser(&request); err != nil {
		fmt.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var venue models.Venue
	if err := bc.DB.Find(&venue, request.VenueID).Error; err != nil {
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

	if err := bc.DB.Create(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create booing",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Booking created successfully",
		"data_booking": booking,
	})
}
