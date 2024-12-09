package main

import (
	"fmt"
	"os"

	"github.com/fsanano/bookshelf/cmd/api"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "bookshelf",
		Commands: []*cli.Command{
			&api.Cmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Print(err.Error())
	}
}
