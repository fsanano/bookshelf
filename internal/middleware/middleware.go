package middleware

import (
	"bookshelf/internal/models"
	"bookshelf/internal/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/signup" {
		// Signup is public
		return c.Next()
	}

	key := c.Get("Key")
	sign := c.Get("Sign")

	if key == "" || sign == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Bad credentials",
		})
	}

	models.UsersStore.Sync.RLock()
	user, ok := models.UsersStore.Data[key]
	models.UsersStore.Sync.RUnlock()

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Bad credentials",
		})
	}

	method := c.Method()
	url := c.Path()
	body := string(c.Body())
	userSecret := user.Secret

	signStr := method + url + string(body) + userSecret

	expectedSign := utils.MD5Sum(signStr)
	log.Println("sign", expectedSign, expectedSign == sign)

	if expectedSign != sign {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			Data:    nil,
			IsOk:    false,
			Message: "Bad credentials",
		})
	}

	// Store user key in locals for downstream handlers
	c.Locals("userKey", user.Key)
	return c.Next()
}
