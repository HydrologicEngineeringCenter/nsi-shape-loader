package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	dynamicstruct "github.com/Ompluscator/dynamic-struct"
	"github.com/jonas-p/go-shp"
	"github.com/stretchr/testify/assert"
)

func TestCore(t *testing.T) {

	cfg := config.Config{
		Dbuser:      "admin",
		Dbpass:      "notPassword",
		Dbname:      "gis",
		Dbtablename: "nsi",
		Dbhost:      "host.docker.internal",
		Dbport:      "25432",
		FilePath:    "/workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp",
	}

	_, err := shp.Open(cfg.FilePath)
	assert.Nil(t, err, "opening shapefile failed")

	fmt.Println("Reading shapefile from: " + cfg.FilePath)
	shpf, err := shp.Open(cfg.FilePath)
	defer shpf.Close()

	fields := shpf.Fields()

	// dynamically allocate struct fields based on shp file columns
	var definition dynamicstruct.DynamicStruct
	for _, f := range fields {
		colName := strings.ToLower(f.String())
		definition = dynamicstruct.ExtendStruct(definition).
			AddField(strings.Title(colName), "", `db:"`+colName+`"`).
			Build()
	}

	fmt.Println(definition)

}
