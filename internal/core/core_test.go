package core

import (
	"fmt"
	"os"
	"testing"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/config"
	"github.com/stretchr/testify/assert"
)

// func TestCore(t *testing.T) {

// 	cfg := config.Config{
// 		Dbuser:      "admin",
// 		Dbpass:      "notPassword",
// 		Dbname:      "gis",
// 		Dbtablename: "nsi_v2022",
// 		Dbhost:      "host.docker.internal",
// 		Dbport:      "25432",
// 		FilePath:    "/workspaces/shape-sql-loader/test/nsi/NSI_V2_Archives/V2022/15001.shp",
// 	}

// 	output, err := Upload(cfg)
// 	assert.Nil(t, err)
// 	fmt.Println(output)
// }

func TestPrep(t *testing.T) {
	cfg := config.Config{
		Mode:        "P",
		StoreConfig: config.StoreConfig{},
	}
	path, err := os.Getwd()
	assert.Nil(t, err)
	dir := "/workspaces/shape-sql-loader"
	if path != dir {
		os.Chdir(dir)
	}
	path, err = os.Getwd()
	assert.Nil(t, err)
	fmt.Println("pwd: " + path)
	err = Prep(cfg)
	assert.Nil(t, err)
}
