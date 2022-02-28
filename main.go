package main

import (
	"log"
	"os"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/core"
	"github.com/urfave/cli/v2"
)

// main entry point into app containing an args parser wrapper
func main() {
	var dbuser, dbpass, dbname, dbhost, dbport, file string
	app := &cli.App{
		Name:   "nsi-loader",
		Usage:  "upload nsi shapefile to postgis database",
		Action: core.Upload,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dbuser",
				Aliases:     []string{"u"},
				Usage:       "database username",
				Destination: &dbuser,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "dbpass",
				Aliases:     []string{"p"},
				Usage:       "database password",
				Destination: &dbpass,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "dbname",
				Aliases:     []string{"n"},
				Usage:       "database name",
				Destination: &dbname,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "dbhost",
				Aliases:     []string{"t"},
				Usage:       "database hostname",
				Destination: &dbhost,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "dbport",
				Aliases:     []string{"o"},
				Usage:       "database access port",
				Destination: &dbport,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "filepath",
				Aliases:     []string{"d"},
				Usage:       "path to input shapefile",
				Destination: &file,
				Required:    true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
