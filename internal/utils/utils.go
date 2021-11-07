package utils

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetId() string {
	return uuid.New().String()
}

func Check(c *fiber.Ctx, err error) bool {
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return false
	}
	return true
}

func Validate(c *fiber.Ctx, s interface{}) bool {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid data sent in the request (unsupported, required fields missing etc)",
		})
		return false
	}
	return true
}
