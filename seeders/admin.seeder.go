package seeders

import (
	"sport-booking-backend/models"
	"sport-booking-backend/utils"

	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) error {
	admin := models.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: "password",
		Role:     models.UserRoleAdmin,
		IsActive: true,
		Phone:    "1234567890",
		Address:  "123 Admin Street, Admin City, Admin Country",
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = hashedPassword

	// Create the admin user
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}
