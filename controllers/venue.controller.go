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
