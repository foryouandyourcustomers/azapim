package main

import (
	"log"
	"os"

	"github.com/foryouandyourcustomers/azapim/internal/cli"
	ucli "github.com/urfave/cli/v2"
)

func main() {
	app := &ucli.App{
		Flags:    cli.GlobalFlags,
		Before:   cli.BeforeFunction,
		Commands: cli.Collection,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
