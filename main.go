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
				Name:     "shppath",
				Aliases:  []string{"s"},
				Usage:    "",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "xlsmeta",
				Aliases:  []string{"x"},
				Usage:    "",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "directory",
				Aliases:  []string{"d"},
				Usage:    "path to input directory containing shapefiles",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "mode",
				Aliases:  []string{"m"},
				Usage:    "P/U. P prepares a config excel templates, U uploads data",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
