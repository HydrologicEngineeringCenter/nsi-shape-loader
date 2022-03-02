package core

import (
	"fmt"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/structutil"
	"github.com/jonas-p/go-shp"
	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

func Core(c *cli.Context) error {
	cfg, err := config.NewConfig(c)
	err = Upload(cfg)
	return err
}

func Upload(cfg config.Config) error {
	st, err := store.NewStore(cfg)
	if err != nil {
		return err
	}

	fmt.Println("Reading shapefile from: " + cfg.FilePath)
	shpf, err := shp.Open(cfg.FilePath)
	defer shpf.Close()

	fields := shpf.Fields()

	var insertTable = dq.TableDataSet{
		Name:       cfg.Dbtablename,
		Statements: map[string]string{},
		Fields:     model.Point{},
	}

	batchSize := 1000
	var records []model.Point
	for shpf.Next() {
		i, _ := shpf.Shape()

		// construct data struct from point
		var newPoint model.Point
		for j, f := range fields {
			val := shpf.ReadAttribute(i, j)
			fieldStr := strings.Title(strings.ToLower(f.String()))
			structutil.SetField(&newPoint, fieldStr, val)
		}
		records = append(records, newPoint)
		if i%batchSize == 0 {
			err = st.DS.Insert(&insertTable).
				Records(&records).
				Batch(true).
				BatchSize(batchSize).
				Execute()
		}
	}
	return err
}
