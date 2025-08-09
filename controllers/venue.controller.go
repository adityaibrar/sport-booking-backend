package controllers

import (
	"sport-booking-backend/models"
	"sport-booking-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// VenueController handles all venue-related HTTP requests
type VenueController struct {
	DB *gorm.DB
}

// NewVenueController creates a new venue controller instance
func NewVenueController(db *gorm.DB) *VenueController {
	return &VenueController{DB: db}
}

// CreateVenue creates a new venue (admin only)
func (vc *VenueController) CreateVenue(c *fiber.Ctx) error {
	var request models.VenueRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	venue := models.Venue{
		Name:         request.Name,
		Category:     request.Category,
		PricePerHour: request.PricePerHour,
		Description:  request.Description,
		Location:     request.Location,
		OpenTime:     request.OpenTime,
		CloseTime:    request.CloseTime,
		Capacity:     request.Capacity,
		IsActive:     true,
	}

	if err := vc.DB.Create(&venue).Error; err != nil {
		return utils.HandleError(c, err, "Failed to create venue")
	}

	response := venue.ToResponse()
	return utils.SuccessResponse(c, "Venue created successfully", response)
}

// GetListVenue retrieves all venues with optional filters
func (vc *VenueController) GetListVenue(c *fiber.Ctx) error {
	var searchParams models.VenueSearchRequest

	// Parse query parameters
	if err := c.QueryParser(&searchParams); err != nil {
		return utils.HandleError(c, err, "Invalid query parameters")
	}

	// Set defaults
	if searchParams.Page <= 0 {
		searchParams.Page = 1
	}
	if searchParams.Limit <= 0 {
		searchParams.Limit = 10
	}
	if searchParams.Limit > 100 {
		searchParams.Limit = 100
	}
	if searchParams.SortBy == "" {
		searchParams.SortBy = "created_at"
	}
	if searchParams.SortOrder == "" {
		searchParams.SortOrder = "desc"
	}

	// Build query
	query := vc.DB.Model(&models.Venue{})

	// Apply filters
	if searchParams.Name != "" {
		query = query.Where("name LIKE ?", "%"+searchParams.Name+"%")
	}

	if searchParams.Category != "" {
		query = query.Where("category = ?", searchParams.Category)
	}

	if searchParams.Location != "" {
		query = query.Where("location LIKE ?", "%"+searchParams.Location+"%")
	}

	if searchParams.MinPrice > 0 {
		query = query.Where("price_per_hour >= ?", searchParams.MinPrice)
	}

	if searchParams.MaxPrice > 0 {
		query = query.Where("price_per_hour <= ?", searchParams.MaxPrice)
	}

	if searchParams.IsActive != nil {
		query = query.Where("is_active = ?", *searchParams.IsActive)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.HandleError(c, err, "Failed to count venues")
	}

	// Apply pagination and sorting
	offset := (searchParams.Page - 1) * searchParams.Limit
	orderClause := searchParams.SortBy + " " + searchParams.SortOrder

	var venues []models.Venue
	if err := query.Order(orderClause).Offset(offset).Limit(searchParams.Limit).Find(&venues).Error; err != nil {
		return utils.HandleError(c, err, "Failed to fetch venues")
	}

	// Convert to response format
	var venueResponses []models.VenueResponse
	for _, venue := range venues {
		venueResponses = append(venueResponses, venue.ToResponse())
	}

	// Create pagination metadata
	pagination := models.NewPaginationMeta(searchParams.Page, searchParams.Limit, total)

	return utils.SuccessResponseWithPagination(c, "Venues retrieved successfully", venueResponses, pagination)
}

// GetDetailVenue retrieves a specific venue by ID
func (vc *VenueController) GetDetailVenue(c *fiber.Ctx) error {
	venueID := c.Params("id")

	var venue models.Venue
	if err := vc.DB.First(&venue, venueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Venue not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch venue")
	}

	response := venue.ToResponse()
	return utils.SuccessResponse(c, "Venue retrieved successfully", response)
}

// UpdateVenue updates an existing venue (admin only)
func (vc *VenueController) UpdateVenue(c *fiber.Ctx) error {
	venueID := c.Params("id")

	var request models.VenueRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	var venue models.Venue
	if err := vc.DB.First(&venue, venueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Venue not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch venue")
	}

	// Update venue fields
	venue.Name = request.Name
	venue.Category = request.Category
	venue.PricePerHour = request.PricePerHour
	venue.Description = request.Description
	venue.Location = request.Location
	venue.OpenTime = request.OpenTime
	venue.CloseTime = request.CloseTime
	venue.Capacity = request.Capacity

	if err := vc.DB.Save(&venue).Error; err != nil {
		return utils.HandleError(c, err, "Failed to update venue")
	}

	response := venue.ToResponse()
	return utils.SuccessResponse(c, "Venue updated successfully", response)
}

// DeleteVenue deletes a venue (admin only)
func (vc *VenueController) DeleteVenue(c *fiber.Ctx) error {
	venueID := c.Params("id")

	var venue models.Venue
	if err := vc.DB.First(&venue, venueID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Venue not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch venue")
	}

	// Check if venue has active bookings
	var activeBookings int64
	if err := vc.DB.Model(&models.Booking{}).
		Where("venue_id = ? AND status IN ?", venueID, []string{"pending", "confirmed"}).
		Count(&activeBookings).Error; err != nil {
		return utils.HandleError(c, err, "Failed to check venue bookings")
	}

	if activeBookings > 0 {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Cannot delete venue with active bookings"), "")
	}

	// Soft delete the venue
	if err := vc.DB.Delete(&venue).Error; err != nil {
		return utils.HandleError(c, err, "Failed to delete venue")
	}

	return utils.SuccessResponse(c, "Venue deleted successfully", nil)
}

// CheckAvailability checks if a venue is available for booking
func (vc *VenueController) CheckAvailability(c *fiber.Ctx) error {
	var request models.BookingAvailabilityRequest

	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	// Parse request date and time
	requestDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD format"), "")
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Invalid time format. Use HH:MM format"), "")
	}

	// Combine date and time to create full datetime
	requestedStartDateTime := time.Date(
		requestDate.Year(), requestDate.Month(), requestDate.Day(),
		startTime.Hour(), startTime.Minute(), 0, 0, time.Local,
	)
	requestedEndDateTime := requestedStartDateTime.Add(time.Duration(request.Duration) * time.Hour)

	// Get the specific venue
	var venue models.Venue
	if err := vc.DB.Where("id = ? AND is_active = ?", request.VenueID, true).First(&venue).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Venue not found or inactive"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch venue")
	}

	// Check if venue is open during requested time
	if !venue.IsOpenAt(request.StartTime) {
		return utils.SuccessResponse(c, "Availability check completed", models.VenueAvailabilityResponse{
			VenueID:     venue.ID,
			Date:        request.Date,
			IsAvailable: false,
		})
	}

	// Check for time overlap with existing bookings
	var count int64
	err = vc.DB.Model(&models.Booking{}).
		Where("venue_id = ? AND status IN ?", venue.ID, []string{"pending", "confirmed"}).
		Where("start_time < ? AND DATE_ADD(start_time, INTERVAL duration HOUR) > ?",
			requestedEndDateTime, requestedStartDateTime).
		Count(&count).Error

	if err != nil {
		return utils.HandleError(c, err, "Failed to check availability")
	}

	isAvailable := count == 0

	availabilityResponse := models.VenueAvailabilityResponse{
		VenueID:     venue.ID,
		Date:        request.Date,
		IsAvailable: isAvailable,
	}

	// Generate available slots if requested
	if isAvailable {
		availabilityResponse.AvailableSlots = venue.GetOperatingHours()
	}

	return utils.SuccessResponse(c, "Availability check completed", availabilityResponse)
}
