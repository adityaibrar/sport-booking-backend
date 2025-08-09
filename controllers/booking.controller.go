package controllers

import (
	"errors"
	"sport-booking-backend/models"
	"sport-booking-backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BookingController handles all booking-related HTTP requests
type BookingController struct {
	DB *gorm.DB
}

// NewBookingController creates a new booking controller instance
func NewBookingController(db *gorm.DB) *BookingController {
	return &BookingController{DB: db}
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking for a venue with validation and conflict checking
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param booking body models.BookingRequest true "Booking request payload"
// @Success 201 {object} models.APIResponse{data=models.BookingResponse}
// @Failure 400 {object} models.APIResponse{error=models.ValidationErrors}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/user/bookings [post]
func (bc *BookingController) CreateBooking(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var request models.BookingRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	// Verify venue exists and is active
	var venue models.Venue
	if err := bc.DB.Where("id = ? AND is_active = ?", request.VenueID, true).First(&venue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Venue not found or inactive"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch venue")
	}

	// Validate booking time
	if request.StartTime.Before(time.Now()) {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Cannot book in the past"), "")
	}

	// Check if venue is open during booking time
	startTimeStr := request.StartTime.Format("15:04")
	endTime := request.StartTime.Add(time.Duration(request.Duration) * time.Hour)
	endTimeStr := endTime.Format("15:04")

	if !venue.IsOpenAt(startTimeStr) || !venue.IsOpenAt(endTimeStr) {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Venue is closed during requested time"), "")
	}

	// Check for booking conflicts
	var conflictCount int64
	err := bc.DB.Model(&models.Booking{}).
		Where("venue_id = ? AND status IN ?", request.VenueID, []string{"pending", "confirmed"}).
		Where("start_time < ? AND DATE_ADD(start_time, INTERVAL duration HOUR) > ?", endTime, request.StartTime).
		Count(&conflictCount).Error

	if err != nil {
		return utils.HandleError(c, err, "Failed to check booking conflicts")
	}

	if conflictCount > 0 {
		return utils.HandleError(c, fiber.NewError(fiber.StatusConflict, "Venue already booked for the selected time slot"), "")
	}

	// Calculate total price
	totalPrice := float64(request.Duration) * venue.PricePerHour

	// Create booking
	booking := models.Booking{
		UserID:     userID,
		VenueID:    request.VenueID,
		StartTime:  request.StartTime,
		Duration:   request.Duration,
		TotalPrice: totalPrice,
		Status:     models.BookingStatusPending,
		Notes:      request.Notes,
	}

	if err := bc.DB.Create(&booking).Error; err != nil {
		return utils.HandleError(c, err, "Failed to create booking")
	}

	// Load associations for response
	if err := bc.DB.Preload("User").Preload("Venue").First(&booking, booking.ID).Error; err != nil {
		return utils.HandleError(c, err, "Failed to load booking details")
	}

	response := booking.ToResponse()
	return utils.SuccessResponse(c, "Booking created successfully", response)
}

// UpdateBookingStatus godoc
// @Summary Update booking status
// @Description Update the status of a booking (admin only)
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Param status body models.BookingUpdateRequest true "Status update request"
// @Success 200 {object} models.APIResponse{data=models.BookingResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/admin/bookings/{id}/status [put]
func (bc *BookingController) UpdateBookingStatus(c *fiber.Ctx) error {
	bookingID := c.Params("id")

	var request models.BookingUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	var booking models.Booking
	if err := bc.DB.Preload("User").Preload("Venue").First(&booking, bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Booking not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch booking")
	}

	// Update booking
	booking.Status = request.Status
	if request.Notes != "" {
		booking.Notes = request.Notes
	}

	if err := bc.DB.Save(&booking).Error; err != nil {
		return utils.HandleError(c, err, "Failed to update booking")
	}

	response := booking.ToResponse()
	return utils.SuccessResponse(c, "Booking status updated successfully", response)
}

// GetUserBookings godoc
// @Summary Get bookings with search and pagination
// @Description Retrieve bookings with optional search filters and pagination
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_name query string false "Filter by user name"
// @Param venue_name query string false "Filter by venue name"
// @Param status query string false "Filter by status" Enums(pending, confirmed, completed, cancelled)
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param sort_by query string false "Sort field" Enums(created_at, start_time, total_price)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} models.APIResponse{data=[]models.BookingResponse,meta=models.PaginationMeta}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/admin/bookings [get]
func (bc *BookingController) GetUserBookings(c *fiber.Ctx) error {
	var searchParams models.BookingSearchRequest

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

	// Validate search parameters
	if err := utils.ValidateStruct(&searchParams); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	// Build query
	query := bc.DB.Model(&models.Booking{}).Preload("User").Preload("Venue")

	// Apply filters
	if searchParams.UserName != "" {
		query = query.Joins("JOIN users ON users.id = bookings.user_id").
			Where("users.name ILIKE ?", "%"+searchParams.UserName+"%")
	}

	if searchParams.VenueName != "" {
		query = query.Joins("JOIN venues ON venues.id = bookings.venue_id").
			Where("venues.name ILIKE ?", "%"+searchParams.VenueName+"%")
	}

	if searchParams.Status != "" {
		query = query.Where("bookings.status = ?", searchParams.Status)
	}

	if searchParams.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", searchParams.StartDate)
		if err != nil {
			return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format"), "")
		}
		query = query.Where("DATE(bookings.start_time) >= ?", startDate)
	}

	if searchParams.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", searchParams.EndDate)
		if err != nil {
			return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format"), "")
		}
		query = query.Where("DATE(bookings.start_time) <= ?", endDate)
	}

	if searchParams.MinPrice > 0 {
		query = query.Where("bookings.total_price >= ?", searchParams.MinPrice)
	}

	if searchParams.MaxPrice > 0 {
		query = query.Where("bookings.total_price <= ?", searchParams.MaxPrice)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.HandleError(c, err, "Failed to count bookings")
	}

	// Apply pagination and sorting
	offset := (searchParams.Page - 1) * searchParams.Limit
	orderClause := "bookings." + searchParams.SortBy + " " + searchParams.SortOrder

	var bookings []models.Booking
	if err := query.Order(orderClause).Offset(offset).Limit(searchParams.Limit).Find(&bookings).Error; err != nil {
		return utils.HandleError(c, err, "Failed to fetch bookings")
	}

	// Convert to response format
	var bookingResponses []models.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, booking.ToResponse())
	}

	// Create pagination metadata
	pagination := models.NewPaginationMeta(searchParams.Page, searchParams.Limit, total)

	return utils.SuccessResponseWithPagination(c, "Bookings retrieved successfully", bookingResponses, pagination)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking if it's in pending status or confirmed but not started yet
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking ID"
// @Success 200 {object} models.APIResponse{data=models.BookingResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/user/bookings/{id}/cancel [post]
func (bc *BookingController) CancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	userID := c.Locals("user_id").(uint)

	var booking models.Booking
	if err := bc.DB.Preload("User").Preload("Venue").First(&booking, bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "Booking not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch booking")
	}

	// Check if user owns this booking
	if booking.UserID != userID {
		return utils.HandleError(c, fiber.NewError(fiber.StatusForbidden, "You can only cancel your own bookings"), "")
	}

	// Check if booking can be cancelled
	if !booking.CanBeCancelled() {
		return utils.HandleError(c, fiber.NewError(fiber.StatusBadRequest, "Booking cannot be cancelled"), "")
	}

	// Update booking status
	booking.Status = models.BookingStatusCancelled
	if err := bc.DB.Save(&booking).Error; err != nil {
		return utils.HandleError(c, err, "Failed to cancel booking")
	}

	response := booking.ToResponse()
	return utils.SuccessResponse(c, "Booking cancelled successfully", response)
}

// HistoryBookingUser godoc
// @Summary Get user's booking history
// @Description Retrieve booking history for a specific user
// @Tags bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param status query string false "Filter by status"
// @Success 200 {object} models.APIResponse{data=[]models.BookingResponse,meta=models.PaginationMeta}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/user/bookings/history/{id} [get]
func (bc *BookingController) HistoryBookingUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status")

	if limit > 100 {
		limit = 100
	}

	// Verify user exists
	var user models.User
	if err := bc.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.HandleError(c, fiber.NewError(fiber.StatusNotFound, "User not found"), "")
		}
		return utils.HandleError(c, err, "Failed to fetch user")
	}

	// Build query
	query := bc.DB.Model(&models.Booking{}).Preload("User").Preload("Venue").
		Where("user_id = ?", userID)

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.HandleError(c, err, "Failed to count bookings")
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	var bookings []models.Booking
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bookings).Error; err != nil {
		return utils.HandleError(c, err, "Failed to fetch booking history")
	}

	// Convert to response format
	var bookingResponses []models.BookingResponse
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, booking.ToResponse())
	}

	// Create pagination metadata
	pagination := models.NewPaginationMeta(page, limit, total)

	return utils.SuccessResponseWithPagination(c, "Booking history retrieved successfully", bookingResponses, pagination)
}
