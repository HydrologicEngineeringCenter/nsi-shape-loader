package main

import (
	"log"
	"os"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/core"
	"github.com/urfave/cli/v2"
)

// main entry point into app containing an args parser wrapper
func main() {
	app := &cli.App{
		Name:   "nsi-loader",
		Usage:  "upload nsi shapefile to postgis database",
		Action: core.Core,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "sql",
				Aliases:  []string{"s"},
				Usage:    "sql connection params",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "schematable",
				Aliases:  []string{"e"},
				Usage:    "table name in format ie. schema.table",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "directory",
				Aliases:  []string{"d"},
				Usage:    "path to input directory containing shapefiles",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
