package utils

import (
	"fmt"
	"reflect"
	"sport-booking-backend/models"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns formatted error response
func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

// HandleValidationErrors converts validation errors to API error format
func HandleValidationErrors(c *fiber.Ctx, err error) error {
	var validationErrors []models.ErrorDetail

	if validationErr, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErr {
			errorDetail := models.ErrorDetail{
				Code:    err.Tag(),
				Field:   err.Field(),
				Message: getValidationMessage(err),
			}
			validationErrors = append(validationErrors, errorDetail)
		}
	}

	response := models.ErrorResponse(
		"Validation failed",
		models.ValidationErrors{Errors: validationErrors},
	)

	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// HandleError returns appropriate error response based on error type
func HandleError(c *fiber.Ctx, err error, defaultMessage string) error {
	if err == nil {
		return nil
	}

	// Handle validation errors
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		return HandleValidationErrors(c, validationErr)
	}

	// Handle fiber errors
	if fiberErr, ok := err.(*fiber.Error); ok {
		response := models.ErrorResponse(fiberErr.Message, nil)
		return c.Status(fiberErr.Code).JSON(response)
	}

	// Handle generic errors
	response := models.ErrorResponse(defaultMessage, err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

// SuccessResponseWithPagination returns success response with pagination
func SuccessResponseWithPagination(c *fiber.Ctx, message string, data interface{}, pagination models.PaginationMeta) error {
	response := models.SuccessResponse(message, data, pagination)
	return c.JSON(response)
}

// SuccessResponse returns simple success response
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	response := models.SuccessResponse(message, data, nil)
	return c.JSON(response)
}
