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
				Name:     "mode",
				Aliases:  []string{"m"},
				Usage:    "prep/upload/access. 'prep' prepares a config excel templates, 'upload' uploads data, 'access' changes access group and role",
				Required: true,
			},

			// xlsPath flag required for both Prep and Upload modes
			&cli.PathFlag{
				Name:    "xlsPath",
				Aliases: []string{"x"},
				Usage:   "",
			},
			// Upload
			&cli.PathFlag{
				Name:    "shpPath",
				Aliases: []string{"s"},
				Usage:   "",
			},

			// consider adding this flag for uploading multiple files
			// &cli.StringFlag{
			// 	Name:     "directory",
			// 	Aliases:  []string{"d"},
			// 	Usage:    "path to input directory containing shapefiles",
			// 	Required: false,
			// },

			// access mode
			&cli.PathFlag{
				Name:    "datasetId",
				Aliases: []string{"i"},
				Usage:   "",
			},
			&cli.PathFlag{
				Name:    "group",
				Aliases: []string{"i"},
				Usage:   "",
			},
			&cli.PathFlag{
				Name:    "role",
				Aliases: []string{"i"},
				Usage:   "",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
