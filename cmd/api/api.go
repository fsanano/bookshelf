package api

import (
	"bookshelf/internal/handlers"
	"bookshelf/internal/middleware"
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

	// Authenticated endpoints
	app.Use(middleware.AuthMiddleware)
	app.Get("/myself", handlers.MyselfHandler)
	app.Post("/books", handlers.AddBookHandler)
	app.Get("/books/:title", handlers.GetBooksByTitleHandler)
	app.Get("/books", handlers.GetBooksHandler)
	app.Patch("/books/:id", handlers.UpdateBookHandler)
	app.Delete("/books/:id", handlers.DeleteBookHandler)
	app.Get("/cleanup", handlers.CleanupHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(app.Listen(":" + port))

	return nil
}
