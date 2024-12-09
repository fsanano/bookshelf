package handlers

import (
	"bookshelf/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func MyselfHandler(c *fiber.Ctx) error {
	userKey := c.Locals("userKey").(string)

	models.UsersStore.Sync.RLock()
	user, ok := models.UsersStore.Data[userKey]
	models.UsersStore.Sync.RUnlock()
	log.Println("user", user, ok)

	if !ok {
		// This theoretically shouldn't happen because we have auth middleware,
		// but let's handle just in case.
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Bad credentials",
		})
	}

	return c.JSON(models.Response{
		Data:    user,
		IsOk:    true,
		Message: "ok",
	})
}
