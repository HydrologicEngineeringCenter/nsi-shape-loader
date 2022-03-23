package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/util"
	"github.com/google/uuid"
	suuid "github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

// Config is for the general app
type Config struct {
	Mode types.Mode
	PathConfig
	StoreConfig
	AccessConfig model.Access
}

type PathConfig struct {
	ShpPath string
	XlsPath string
}

// StoreConfig holds only params required for database connection
type StoreConfig struct {
	Dbuser string
	Dbpass string
	Dbname string
	Dbhost string
	Dbport string
}

func (c *StoreConfig) Rdbmsconfig() dq.RdbmsConfig {
	return dq.RdbmsConfig{
		Dbuser:   c.Dbuser,
		Dbpass:   c.Dbpass,
		Dbhost:   c.Dbhost,
		Dbport:   c.Dbport,
		Dbname:   c.Dbname,
		DbDriver: "postgres",
		DbStore:  "pgx",
	}
}

// NewConfig generates new config from cli args context
func NewConfig(c *cli.Context) (Config, error) {

	mode := types.ModeReverse[c.String("mode")]
	// validate for valid mode
	if mode != types.Access && mode != types.Prep && mode != types.Upload {
		return Config{}, errors.New("invalid mode, respecify --mode")
	}

	var storeCfg StoreConfig
	var pathCfg PathConfig
	var accessCfg model.Access

	// validate nonempty db access info
	if mode == types.Prep || mode == types.Upload {
		// validate sql db connection is not empty
		sqlConn := c.String("sqlConn")
		if sqlConn == "" {
			return Config{}, errors.New("invalid sql connection string, --sqlConn should not be empty")
		}
		var user, pass, name, host, port string
		user = strings.ReplaceAll(regexp.MustCompile(`user=.+?\s`).FindString(sqlConn), " ", "")
		pass = strings.ReplaceAll(regexp.MustCompile(`password=.+?\s`).FindString(sqlConn), " ", "")
		name = strings.ReplaceAll(regexp.MustCompile(`dbname=.+?\s`).FindString(sqlConn), " ", "")
		host = strings.ReplaceAll(regexp.MustCompile(`host=.+?\s`).FindString(sqlConn), " ", "")
		port = strings.ReplaceAll(regexp.MustCompile(`port=.+?\s`).FindString(sqlConn), " ", "")

		if util.StrContains([]string{sqlConn, user, pass, name, host, port}, "") {
			return Config{}, errors.New("invalid sql connection string, respecify --sqlConn")
		}

		storeCfg = StoreConfig{
			Dbuser: user,
			Dbpass: pass,
			Dbname: name,
			Dbhost: host,
			Dbport: port,
		}
		pathCfg = PathConfig{
			ShpPath: c.Path("shpPath"),
			XlsPath: c.Path("xlsPath"),
		}

		if pathCfg.ShpPath == "" {
			return Config{}, errors.New("invalid path to shp file, --shpPath should not be empty")
		}
		if mode == types.Upload && pathCfg.XlsPath == "" {
			return Config{}, errors.New("invalid path to xls file, --xlsPath should not be empty")
		}
	}

	if mode == types.Access {
		role := types.Role(c.String("role"))
		if !util.StrContains([]string{string(types.Admin), string(types.Owner), string(types.User)}, string(role)) {
			return Config{}, errors.New("invalid role, --role accepts only admin, owner, or user")
		}
		group := c.String("group")
		if group == "" {
			return Config{}, errors.New("invalid group, --group must not be empty")
		}
		sDatasetId := c.String("datasetId")
		if sDatasetId == "" {
			return Config{}, errors.New("invalid datasetId, --datasetId must not be empty")
		}
		datasetId, err := suuid.FromString(sDatasetId)
		gDatasetId, err := uuid.FromBytes(datasetId.Bytes())
		if err != nil {
			return Config{}, err
		}
		accessCfg = model.Access{
			DatasetId: gDatasetId,
			Group:     group,
			Role:      role,
		}
	}

	return Config{
		Mode:         mode,
		PathConfig:   pathCfg,
		StoreConfig:  storeCfg,
		AccessConfig: accessCfg,
	}, nil
}
