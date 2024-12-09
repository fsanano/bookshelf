package api

import (
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	cli "github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "api",
	Usage:  "API for providing job's candidase",
	Action: run,
}

func run(c *cli.Context) error {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	debug.SetPanicOnFault(true) // will cause panic instead program fault in order to keep application alive

	app.Use(pprof.New())

	return nil
}
