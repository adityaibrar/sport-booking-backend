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

func (bc *BookingController) UpdateBookingStatus(c *fiber.Ctx) error {
	idVenue := c.Params("id")
	var booking models.Booking
	if err := bc.DB.Where("id = ? AND status = ?", idVenue, "pending").First(&booking).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found or already confirmed",
		})
	}

	booking.Status = "confirmed"
	if err := bc.DB.Save(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to confirm booking",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking confirmed successfully",
	})
}

func (bc *BookingController) GetUserBookings(c *fiber.Ctx) error {
	// Get search parameter from query
	searchName := c.Query("name", "")

	var bookings []models.Booking
	query := bc.DB.Preload("User").Preload("Venue")

	// Filter by user name if provided
	if searchName != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").
			Where("users.name LIKE ?", "%"+searchName+"%")
	}

	// Order by creation date (newest first)
	if err := query.Order("bookings.created_at DESC").Find(&bookings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch bookings",
			"debug": err.Error(),
		})
	}

	// Transform to response format
	var bookingResponses []models.BookingResponse
	for _, booking := range bookings {
		response := models.BookingResponse{
			ID:          booking.ID,
			StartTime:   booking.StartTime,
			Duration:    booking.Duration,
			TotalPrice:  booking.TotalPrice,
			Status:      booking.Status,
			PaymentQRIS: booking.PaymentQRIS,
			CreatedAt:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// User data
		response.User.ID = booking.User.ID
		response.User.Name = booking.User.Name
		response.User.Email = booking.User.Email

		// Venue data
		response.Venue.ID = booking.Venue.ID
		response.Venue.Name = booking.Venue.Name
		response.Venue.Category = booking.Venue.Category
		response.Venue.PricePerHour = booking.Venue.PricePerHour
		response.Venue.Description = booking.Venue.Description

		bookingResponses = append(bookingResponses, response)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Bookings retrieved successfully",
		"data_booking": bookingResponses,
		"total":        len(bookingResponses),
		"filter": fiber.Map{
			"name": searchName,
		},
	})
}

func (bc *BookingController) CancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking
	if err := bc.DB.Where("id = ? AND status = ?", id, "pending").First(&booking).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Booking not found or cannot be cancelled",
		})
	}
	booking.Status = "cancelled"
	if err := bc.DB.Save(&booking).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel booking",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking cancelled successfully",
	})
}

func (bc *BookingController) HistoryBookingUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	var bookings []models.Booking
	// preload untuk mengambil data user dan venue
	if err := bc.DB.Preload("User").Preload("Venue").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&bookings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch booking history",
			"debug": err.Error(),
		})
	}

	var bookingResponses []models.BookingResponse
	for _, booking := range bookings {
		response := models.BookingResponse{
			ID:          booking.ID,
			StartTime:   booking.StartTime,
			Duration:    booking.Duration,
			TotalPrice:  booking.TotalPrice,
			Status:      booking.Status,
			PaymentQRIS: booking.PaymentQRIS,
			CreatedAt:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		response.User.ID = booking.User.ID
		response.User.Name = booking.User.Name
		response.User.Email = booking.User.Email
		response.User.Phone = fmt.Sprintf("%d", booking.User.Phone)

		response.Venue.ID = booking.Venue.ID
		response.Venue.Name = booking.Venue.Name
		response.Venue.Category = booking.Venue.Category
		response.Venue.PricePerHour = booking.Venue.PricePerHour
		response.Venue.Description = booking.Venue.Description

		bookingResponses = append(bookingResponses, response)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Booking history retrieved successfully",
		"data_booking": bookingResponses,
		"total":        len(bookingResponses),
	})
}
