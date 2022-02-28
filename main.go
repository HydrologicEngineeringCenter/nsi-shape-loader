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
		Action: core.Upload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dbuser",
				Aliases:  []string{"u"},
				Usage:    "database username",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbpass",
				Aliases:  []string{"p"},
				Usage:    "database password",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbhost",
				Aliases:  []string{"t"},
				Usage:    "database hostname",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbport",
				Aliases:  []string{"o"},
				Usage:    "database access port",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbname",
				Aliases:  []string{"n"},
				Usage:    "database name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dbtname",
				Aliases:  []string{"a"},
				Usage:    "table name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filepath",
				Aliases:  []string{"d"},
				Usage:    "path to input shapefile",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
