package handlers

import (
	"bookshelf/internal/models"

	"github.com/gofiber/fiber/v2"
)

func SignupHandler(c *fiber.Ctx) error {
	var u models.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Invalid request body",
		})
	}

	if u.Key == "" || u.Secret == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Key and Secret required",
		})
	}

	models.UsersStore.Sync.Lock()
	defer models.UsersStore.Sync.Unlock()

	if _, exists := models.UsersStore.Data[u.Key]; exists {
		return c.Status(fiber.StatusConflict).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "User already exists",
		})
	}

	models.UsersStore.Data[u.Key] = u

	return c.JSON(models.Response{
		Data:    u,
		IsOk:    true,
		Message: "ok",
	})
}