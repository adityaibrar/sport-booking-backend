package controllers

import (
	"fmt"
	"sport-booking-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type VenueController struct {
	DB *gorm.DB
}

func NewVenueController(db *gorm.DB) *VenueController {
	return &VenueController{DB: db}
}

func (vc *VenueController) CreateVenue(c *fiber.Ctx) error {
	var request models.VenueRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	venue := models.Venue{
		Name:         request.Name,
		Category:     request.Category,
		PricePerHour: request.PricePerHour,
		Description:  request.Description,
	}

	if err := vc.DB.Create(&venue).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create venue",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Successfully create venue",
		"data_venue": venue,
	})
}

func (vc *VenueController) UpdateVenue(c *fiber.Ctx) error {
	id := c.Params("id")
	var venue models.Venue

	if err := vc.DB.First(&venue, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Venue not found",
		})
	}

	var updateVenue models.VenueRequest

	if err := c.BodyParser(&updateVenue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	venue.Name = updateVenue.Name
	venue.Category = updateVenue.Category
	venue.PricePerHour = updateVenue.PricePerHour
	venue.Description = updateVenue.Description

	if err := vc.DB.Save(&venue).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update a venue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Successfully update a venue",
		"data_venue": venue,
	})
}

func (vc *VenueController) GetListVenue(c *fiber.Ctx) error {
	var venues []models.Venue

	if err := vc.DB.Find(&venues).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get venues",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Successfuly get list venue",
		"data_venues": venues,
	})
}

func (vc *VenueController) GetDetailVenue(c *fiber.Ctx) error {
	id := c.Params("id")
	var venue models.Venue

	if err := vc.DB.Find(&venue, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Venue not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Succesfully get venue",
		"data_venue": venue,
	})
}

func (vc *VenueController) DeleteVenue(c *fiber.Ctx) error {
	id := c.Params("id")
	var venue models.Venue

	if err := vc.DB.First(&venue, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Venue not found",
		})
	}

	if err := vc.DB.Delete(&venue).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed delete venue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Venue successfuly deleted",
	})
}

func (vc *VenueController) CheckAvailability(c *fiber.Ctx) error {
	var request models.BookingSearch

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Parse requested start and end times (just time, not full datetime)
	requestedStartTime := request.StartTime // e.g., "19:00"

	// Parse the time to calculate end time
	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid time format. Use HH:MM format",
		})
	}
	endTime := startTime.Add(time.Duration(request.Duration) * time.Hour)
	requestedEndTime := endTime.Format("15:04") // e.g., "21:00"

	// Get all venues based on category filter
	var venues []models.Venue
	query := vc.DB
	if request.Category != "" {
		query = query.Where("category = ?", request.Category)
	}
	if err := query.Find(&venues).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get venues",
		})
	}

	// Check availability for each venue
	availableVenues := make([]models.Venue, 0)

	for _, venue := range venues {
		var count int64

		// Check for time overlap with existing bookings
		// Two time ranges (A and B) overlap if: A.start < B.end AND B.start < A.end
		// For our case: requested.start < existing.end AND existing.start < requested.end
		err := vc.DB.Model(&models.Booking{}).
			Where("venue_id = ? AND status IN (?, ?)", venue.ID, "pending", "confirmed").
			Where(`
				? < ADDTIME(start_time, SEC_TO_TIME(duration * 3600)) AND 
				start_time < ?
			`, requestedStartTime, requestedEndTime).
			Count(&count).Error

		if err != nil {
			fmt.Printf("Error checking availability for venue %d: %v\n", venue.ID, err)
			continue
		}

		fmt.Printf("Venue %d (%s): Found %d conflicting bookings\n", venue.ID, venue.Name, count)

		// If no overlapping bookings found, venue is available
		if count == 0 {
			availableVenues = append(availableVenues, venue)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Available venues retrieved successfully",
		"data":    availableVenues,
		"total":   len(availableVenues),
	})
}

func (vc *VenueController) GetAllVenues(c *fiber.Ctx) error {
	// Get category from query parameter (optional)
	category := c.Query("category")

	// Get all venues with their bookings (preload bookings)
	var venues []models.Venue
	query := vc.DB.Preload("Bookings", func(db *gorm.DB) *gorm.DB {
		return db.Where("status IN (?)", []string{"pending", "confirmed", "completed"}).
			Order("start_time ASC")
	})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&venues).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get venues",
		})
	}

	type BookingDetail struct {
		ID          uint    `json:"id"`
		UserID      uint    `json:"user_id"`
		StartTime   string  `json:"start_time"`
		Duration    int     `json:"duration"`
		EndTime     string  `json:"end_time"`
		TotalPrice  float64 `json:"total_price"`
		Status      string  `json:"status"`
		PaymentQRIS string  `json:"payment_qris,omitempty"`
	}

	type VenueWithStatus struct {
		ID              uint            `json:"id"`
		Name            string          `json:"name"`
		Category        string          `json:"category"`
		PricePerHour    float64         `json:"price_per_hour"`
		Description     string          `json:"description"`
		TotalBookings   int             `json:"total_bookings"`
		ActiveBookings  int             `json:"active_bookings"`
		BookingDetails  []BookingDetail `json:"booking_details,omitempty"`
		IsCurrentlyBusy bool            `json:"is_currently_busy"`
		CreatedAt       interface{}     `json:"created_at"`
		UpdatedAt       interface{}     `json:"updated_at"`
	}

	venuesWithStatus := make([]VenueWithStatus, 0, len(venues))
	currentTime := time.Now().Format("15:04")

	for _, venue := range venues {
		activeBookings := 0
		isCurrentlyBusy := false
		bookingDetails := make([]BookingDetail, 0)

		for _, booking := range venue.Bookings {
			if booking.Status == "pending" || booking.Status == "confirmed" {
				activeBookings++
			}
			if booking.Status == "confirmed" || booking.Status == "pending" {
				startTime, err := time.Parse("15:04", booking.StartTime)
				if err == nil {
					endTime := startTime.Add(time.Duration(booking.Duration) * time.Hour)
					currentTimeParsed, err := time.Parse("15:04", currentTime)
					if err == nil {
						if (currentTimeParsed.After(startTime) || currentTimeParsed.Equal(startTime)) &&
							currentTimeParsed.Before(endTime) {
							isCurrentlyBusy = true
						}
					}
				}
			}
			endTime := ""
			if startTime, err := time.Parse("15:04", booking.StartTime); err == nil {
				endTime = startTime.Add(time.Duration(booking.Duration) * time.Hour).Format("15:04")
			}
			bookingDetails = append(bookingDetails, BookingDetail{
				ID:          booking.ID,
				UserID:      booking.UserID,
				StartTime:   booking.StartTime,
				Duration:    booking.Duration,
				EndTime:     endTime,
				TotalPrice:  booking.TotalPrice,
				Status:      booking.Status,
				PaymentQRIS: booking.PaymentQRIS,
			})
		}

		venueWithStatus := VenueWithStatus{
			ID:              venue.ID,
			Name:            venue.Name,
			Category:        venue.Category,
			PricePerHour:    venue.PricePerHour,
			Description:     venue.Description,
			TotalBookings:   len(venue.Bookings),
			ActiveBookings:  activeBookings,
			BookingDetails:  bookingDetails,
			IsCurrentlyBusy: isCurrentlyBusy,
			CreatedAt:       venue.CreatedAt,
			UpdatedAt:       venue.UpdatedAt,
		}
		venuesWithStatus = append(venuesWithStatus, venueWithStatus)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Venues retrieved successfully",
		"data":    venuesWithStatus,
		"total":   len(venuesWithStatus),
		"filter": map[string]interface{}{
			"category": category,
		},
	})
}

func (vc *VenueController) GetAvailableVenues(c *fiber.Ctx) error {
	var request models.BookingSearch

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Parse requested start and end times (just time, not full datetime)
	requestedStartTime := request.StartTime // e.g., "19:00"

	// Parse the time to calculate end time
	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid time format. Use HH:MM format",
		})
	}
	endTime := startTime.Add(time.Duration(request.Duration) * time.Hour)
	requestedEndTime := endTime.Format("15:04") // e.g., "21:00"

	var venues []models.Venue
	if err := vc.DB.Find(&venues).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get venues",
		})
	}

	availableVenues := make([]models.Venue, 0)

	for _, venue := range venues {
		var count int64

		err := vc.DB.Model(&models.Booking{}).
			Where("venue_id = ? AND status IN (?, ?)", venue.ID, "pending", "confirmed").
			Where(`
				? < ADDTIME(start_time, SEC_TO_TIME(duration * 3600)) AND 
				start_time < ?
			`, requestedStartTime, requestedEndTime).
			Count(&count).Error

		if err != nil {
			fmt.Printf("Error checking availability for venue %d: %v\n", venue.ID, err)
			continue
		}

		if count == 0 {
			availableVenues = append(availableVenues, venue)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Available venues retrieved successfully",
		"data":    availableVenues,
		"total":   len(availableVenues),
	})
}
