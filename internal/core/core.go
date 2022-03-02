package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/structutil"
	dynamicstruct "github.com/Ompluscator/dynamic-struct"
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

	// dynamically allocate struct fields based on shp file columns
	baseDef := dynamicstruct.NewStruct().
		Build().
		New()
	var definition dynamicstruct.DynamicStruct
	for _, f := range fields {
		colName := strings.ToLower(f.String())
		definition = dynamicstruct.ExtendStruct(baseDef).
			AddField(strings.Title(colName), "", `db:"`+colName+`"`).
			Build()
	}

	templateType := reflect.TypeOf(definition)

	var insertTable = dq.TableDataSet{
		Name:       cfg.Dbtablename,
		Statements: map[string]string{},
		Fields:     definition,
	}

	// dataSlice := reflect.SliceOf(templateType)
	batchSize := 1000
	dataSlice := definition.NewSliceOfStructs()
	fmt.Println(dataSlice)
	for shpf.Next() {
		i, _ := shpf.Shape()

		// construct data struct from point
		newDynStruct := reflect.ValueOf(templateType)
		for j, f := range fields {
			val := shpf.ReadAttribute(i, j)
			structutil.SetField(&newDynStruct, strings.Title(strings.ToLower(f.String())), val)
		}
		// dataSlice = append(dataSlice, newDynStruct)
		if i%batchSize == 0 {
			err = st.DS.Insert(&insertTable).
				Records(&dataSlice).
				Batch(true).
				BatchSize(batchSize).
				Execute()
		}
	}
	return err
}
