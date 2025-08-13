package controllers

import (
	"sport-booking-backend/models"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ReviewController struct {
	DB *gorm.DB
}

func NewReviewController(db *gorm.DB) *ReviewController {
	return &ReviewController{DB: db}
}

func (rc *ReviewController) CreateReview(c *fiber.Ctx) error {
	var request models.ReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.HandleError(c, err, "Invalid request body")
	}

	if err := utils.ValidateStruct(&request); err != nil {
		return utils.HandleValidationErrors(c, err)
	}

	if err := rc.DB.Preload("User").Preload("Venue").First(&models.Booking{}, "user_id = ? AND venue_id = ?", request.UserID, request.VenueID).Error; err != nil {
		return utils.HandleError(c, err, "Failed to create review before booking")
	}

	review := models.Review{
		UserID:  request.UserID,
		VenueID: request.VenueID,
		Rating:  request.Rating,
		Comment: request.Comment,
	}

	if err := rc.DB.Create(&review).Error; err != nil {
		return utils.HandleError(c, err, "Failed to create review")
	}

	if err := rc.DB.Preload("User").Preload("Venue").First(&review, "user_id = ? AND venue_id = ?", request.UserID, request.VenueID).Error; err != nil {
		return utils.HandleError(c, err, "Failed to retrieve created review")
	}

	response := review.ToResponse()
	return utils.SuccessResponse(c, "Created review successfully", response)
}

func (rc *ReviewController) GetReview(c *fiber.Ctx) error {
	var reviewFilter models.ReviewSearchRequest

	if err := c.QueryParser(&reviewFilter); err != nil {
		return utils.HandleError(c, err, "Invalid query parameters")
	}

	if reviewFilter.Page <= 0 {
		reviewFilter.Page = 1
	}
	if reviewFilter.Limit <= 0 {
		reviewFilter.Limit = 10
	}
	if reviewFilter.Limit > 100 {
		reviewFilter.Limit = 100
	}
	if reviewFilter.SortBy == "" {
		reviewFilter.SortBy = "created_at"
	}
	if reviewFilter.SortOrder == "" {
		reviewFilter.SortOrder = "desc"
	}

	query := rc.DB.Model(&models.Review{}).Preload("User").Preload("Venue")

	if reviewFilter.VenueCategory != "" {
		query = query.Joins("JOIN venues ON venues.id = reviews.venue_id").Where("venues.category = ?", reviewFilter.VenueCategory)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return utils.HandleError(c, err, "Failed to count reviews")
	}

	offset := (reviewFilter.Page - 1) * reviewFilter.Limit
	orderClause := reviewFilter.SortBy + " " + reviewFilter.SortOrder

	var reviews []models.Review
	if err := query.Order(orderClause).Offset(offset).Limit(reviewFilter.Limit).Find(&reviews).Error; err != nil {
		return utils.HandleError(c, err, "Failed to retrieve reviews")
	}

	if err := rc.DB.Preload("User").Preload("Venue").Find(&reviews).Error; err != nil {
		return utils.HandleError(c, err, "Failed to retrieve created review")
	}

	var reviewResponse []models.ReviewResponse
	for _, review := range reviews {
		reviewResponse = append(reviewResponse, review.ToResponse())
	}

	pagination := models.NewPaginationMeta(reviewFilter.Page, reviewFilter.Limit, total)
	return utils.SuccessResponseWithPagination(c, "Retrieved reviews successfully", reviewResponse, pagination)
}

func (rc *ReviewController) DeleteReview(c *fiber.Ctx) error {
	reviewId := c.Params("id")
	userId := c.Params("user_id")

	var review models.Review
	if err := rc.DB.First(&review, reviewId, userId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, err, "Review not found")
		}
		return utils.HandleError(c, err, "Failed to retrieve review")
	}

	if err := rc.DB.Delete(&review).Error; err != nil {
		return utils.HandleError(c, err, "Failed to delete review")
	}
	return utils.SuccessResponse(c, "Review deleted successfully", nil)
}

func (rc *ReviewController) HistoryReviewByUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	var reviews []models.Review
	if err := rc.DB.Where("user_id = ?", userId).Find(&reviews).Error; err != nil {
		return utils.HandleError(c, err, "Failed to retrieve user reviews")
	}

	var reviewResponses []models.ReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, review.ToResponse())
	}

	return utils.SuccessResponse(c, "Retrieved user reviews successfully", reviewResponses)
}
