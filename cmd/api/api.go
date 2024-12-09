package api

import (
	"bookshelf/internal/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	cli "github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "api",
	Usage:  "API for providing job's candidase",
	Action: run,
}

func run(c *cli.Context) error {
	app := fiber.New()

	// Public endpoint
	app.Post("/signup", handlers.SignupHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))

	return nil
}
