package controllers

import (
	"sport-booking-backend/models"

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
