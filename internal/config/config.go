package config

import (
	"errors"
	"fmt"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/util"
	"github.com/dlclark/regexp2"
	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

// GLOBAL VARS DUMP
const (
	APP_NAME              = "seahorse"
	APP_VERSION           = "0.5.0"
	DB_SCHEMA             = "nsiv29test" // CAUTION! CHANGING TO A PRODUCTION SCHEMA CAN OVERRIDE DATA
	BASE_META_XLSX_PATH   = "./assets/baseMetaData.xlsx"
	COPY_XLSX_PATH        = "./metadata.xlsx"
	NATIONAL_MAP_BASE_URL = "https://tnmaccess.nationalmap.gov/api/v1/products?"
)

// Config is for the general app
type Config struct {
	Mode types.Mode
	PathConfig
	StoreConfig
	AccessConfig
	ElevationConfig
}

type PathConfig struct {
	ShpPath string
	XlsPath string
}

// StoreConfig holds only params required for database connection
type StoreConfig struct {
	ConnStr string
	Dbuser  string
	Dbpass  string
	Dbname  string
	Dbhost  string
	Dbport  string
}

type AccessConfig struct {
	Group  string
	Role   types.Role
	UserId string
}

type ElevationConfig struct {
	Dataset string
	Version string
	Quality types.Quality
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
func NewConfig(c *cli.Context, mode types.Mode) (Config, error) {

	// validate for valid mode
	if mode != types.Access && mode != types.Prep && mode != types.Upload && mode != types.Elevation {
		return Config{}, errors.New(fmt.Sprintf(
			"invalid mode, --mode can only be %s, %s, %s, or %s",
			types.Access,
			types.Prep,
			types.Upload,
			types.Elevation,
		))
	}

	var storeCfg StoreConfig
	var pathCfg PathConfig
	var accessCfg AccessConfig
	var elevCfg ElevationConfig

	// validate sql connection creds
	if mode == types.Access || mode == types.Upload || mode == types.Elevation {
		sqlConn := c.String("sqlConn")
		if sqlConn == "" {
			return Config{}, errors.New("invalid sql connection string, --sqlConn should not be empty")
		}
		var user, pass, database, host, port string
		// std lib regex doesn't support lookahead and lookbehind
		var re *regexp2.Regexp
		var m *regexp2.Match
		var err error
		sqlConnParamsMap := map[string]string{}
		sqlConnParams := []string{"user", "password", "host", "port", "database"}
		for _, param := range sqlConnParams {
			re = regexp2.MustCompile(fmt.Sprintf(`(?<=%s=).+?(?=\s|$)`, param), 0)
			m, err = re.FindStringMatch(sqlConn)
			if err != nil || m == nil || m.String() == "" {
				return Config{}, errors.New(fmt.Sprintf("invalid sql connection string, unable to parse '%s' argument", param))
			}
			sqlConnParamsMap[param] = m.String()
		}

		if util.StrContains([]string{sqlConn, user, pass, database, host, port}, "") {
			return Config{}, errors.New("invalid sql connection string, respecify --sqlConn")
		}
		storeCfg = StoreConfig{
			ConnStr: sqlConn,
			Dbuser:  sqlConnParamsMap["user"],
			Dbpass:  sqlConnParamsMap["password"],
			Dbname:  sqlConnParamsMap["database"],
			Dbhost:  sqlConnParamsMap["host"],
			Dbport:  sqlConnParamsMap["port"],
		}
	}

	// validate file pathings
	if mode == types.Prep || mode == types.Upload {
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

	// validate access mod params
	if mode == types.Access {
		role := types.Role(c.String("role"))
		if !util.StrContains([]string{string(types.Admin), string(types.Owner), string(types.User)}, string(role)) {
			return Config{}, errors.New(fmt.Sprintf(
				"invalid role, --role accepts only %s, %s, or %s",
				types.Admin,
				types.Owner,
				types.User,
			))
		}
		group := c.String("group")
		if group == "" {
			return Config{}, errors.New("invalid group, --group must not be empty")
		}
		user := c.String("user")
		if user == "" {
			return Config{}, errors.New("invalid user id, --user must not be empty")
		}
		accessCfg = AccessConfig{
			Group:  group,
			Role:   role,
			UserId: user,
		}
	}

	//validate elevation params
	if mode == types.Elevation {
		role := types.Role(c.String("role"))
		if !util.StrContains([]string{string(types.Admin), string(types.Owner), string(types.User)}, string(role)) {
			return Config{}, errors.New(fmt.Sprintf(
				"invalid role, --role accepts only %s, %s, or %s",
				types.Admin,
				types.Owner,
				types.User,
			))
		}
		d := c.String("dataset")
		if d == "" {
			return Config{}, errors.New("invalid dataset, --dataset must not be empty")
		}
		v := c.String("version")
		if v == "" {
			return Config{}, errors.New("invalid dataset version, --version must not be empty")
		}
		q := c.String("quality")
		if q == "" {
			return Config{}, errors.New("invalid quality, --quality must not be empty")
		}
		elevCfg = ElevationConfig{
			Dataset: d,
			Version: v,
			Quality: types.QualityReverse[q],
		}
	}

	return Config{
		Mode:            mode,
		PathConfig:      pathCfg,
		StoreConfig:     storeCfg,
		AccessConfig:    accessCfg,
		ElevationConfig: elevCfg,
	}, nil
}
