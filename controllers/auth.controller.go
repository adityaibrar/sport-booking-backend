package controllers

import (
	"fmt"
	"sport-booking-backend/models"
	"sport-booking-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var input models.RegisterRequest
	fmt.Printf("Registering user: %+v\n", input)
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var count int64
	if err := ac.DB.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing user",
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hashing password",
		})
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Address:  input.Address,
	}

	if err := ac.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	userResponse := fiber.Map{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"phone":   user.Phone,
		"address": user.Address,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "User successfully created",
		"user_data": userResponse,
	})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var request models.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username and password is required",
		})
	}

	var user models.User
	if err := ac.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email",
		})
	}

	if !utils.ChechPasswordHash(request.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Wrong password",
		})
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	reponse := fiber.Map{
		"status":  true,
		"message": "Login successful",
		"data_user": fiber.Map{
			"token": token,
			"user":  user,
		},
	}

	return c.Status(fiber.StatusOK).JSON(reponse)

}
