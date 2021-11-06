package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/renbou/dontstress/internal/dao"
	"github.com/renbou/dontstress/internal/models"
	"strings"
)

func auth(c *fiber.Ctx) error {
	admin := models.Admin{Id: strings.TrimSpace(c.Get("Authorization"))}
	if admin.Id != "" {
		_, err := dao.AdminDao().Get(&admin)
		if err == nil {
			return c.Next()
		}
	}
	return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized request")
}

func New() fiber.Handler {
	return auth
}
