package core

import (
	"testing"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
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

	err := Upload(cfg)
	assert.Nil(t, err)
}
