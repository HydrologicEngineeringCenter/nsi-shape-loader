package main

import (
	"log"
	"os"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/core"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/global"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/urfave/cli/v2"
)

// main entry point into app containing an args parser wrapper
func main() {

	app := &cli.App{
		Name:    global.APP_NAME,
		Version: global.APP_VERSION,
		Usage:   "Upload ESRI shapefiles to PostGIS database",
		Commands: []*cli.Command{
			{
				Name:  "prepare",
				Usage: "Prepare an excel config template",
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
				Name:  "mod",
				Usage: "Options to modify data",
				Subcommands: []*cli.Command{
					{
						Name:  "inventory",
						Usage: "Upload shp file to PostGIS",
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
						Name:  "user",
						Usage: "Add user and their role to group",
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
					{
						Name:  "elevation",
						Usage: "Add elevation data to inventory table",
						Action: func(c *cli.Context) error {
							err := core.Core(c, types.Elevation)
							return err
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "dataset",
								Aliases:  []string{"d"},
								Usage:    "Dataset name",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "version",
								Aliases:  []string{"v"},
								Usage:    "Dataset version",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "quality",
								Aliases:  []string{"q"},
								Usage:    "Dataset quality",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "sqlConn",
								Aliases:  []string{"s"},
								Usage:    "PostGIS connection string",
								Required: true,
							},
						},
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
