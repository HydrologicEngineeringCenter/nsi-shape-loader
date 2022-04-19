package main

import (
	"log"
	"os"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/core"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/urfave/cli/v2"
)

// main entry point into app containing an args parser wrapper
func main() {

	app := &cli.App{
		Name:    config.APP_NAME,
		Version: config.APP_VERSION,
		Usage:   "Upload ESRI shapefiles to PostGIS database",
		Commands: []*cli.Command{
			{
				Name:    "prepare",
				Aliases: []string{"p"},
				Usage:   "Prepare a excel config template",
				Action: func(c *cli.Context) error {
					err := core.Core(c, types.Prep)
					return err
				},
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:     "shpPath",
						Aliases:  []string{"s"},
						Usage:    "Path to shp file",
						Required: true,
					},
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "upload shp file to PostGIS",
				Action: func(c *cli.Context) error {
					err := core.Core(c, types.Upload)
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "sqlConn",
						Aliases:  []string{"s"},
						Usage:    "PostGIS connection string",
						Required: true,
					},
					&cli.PathFlag{
						Name:     "xlsPath",
						Aliases:  []string{"x"},
						Usage:    "Path to metadata xlsx file",
						Required: true,
					},
					&cli.PathFlag{
						Name:     "shpPath",
						Aliases:  []string{"p"},
						Usage:    "Path to shp file",
						Required: true,
					},
				},
			},
			{
				Name:    "adduser",
				Aliases: []string{"a"},
				Usage:   "add user and their role to group",
				Action: func(c *cli.Context) error {
					err := core.Core(c, types.Access)
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "sqlConn",
						Aliases:  []string{"s"},
						Usage:    "PostGIS connection string",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "user",
						Aliases:  []string{"u"},
						Usage:    "user id",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "group",
						Aliases:  []string{"g"},
						Usage:    "group name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "role",
						Aliases:  []string{"r"},
						Usage:    "admin / owner / user",
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
