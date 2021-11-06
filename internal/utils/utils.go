package utils

import (
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
