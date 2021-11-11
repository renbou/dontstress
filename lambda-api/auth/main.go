package auth

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
	"strings"
)

func auth(c *fiber.Ctx) error {
	admin := models.Admin{Id: strings.TrimSpace(c.Get("Authorization"))}
	validate := validator.New()
	if err := validate.Struct(admin); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized request")
	}

	_, err := dao.AdminDao().Get(&admin)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized request")

	}

	return c.Next()
}

func New() fiber.Handler {
	return auth
}
