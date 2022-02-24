package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "nsi-loader",
		Usage: "upload nsi shapefile to postgis database",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
