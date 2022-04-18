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

const (
	DB_SCHEMA           = "nsiv29test" // CAUTION! CHANGING TO A PRODUCTION SCHEMA CAN OVERRIDE DATA
	APP_VERSION         = "0.5.0"
	BASE_META_XLSX_PATH = "./assets/baseMetaData.xlsx"
	COPY_XLSX_PATH      = "./metadata.xlsx"
)

// Config is for the general app
type Config struct {
	Mode types.Mode
	PathConfig
	StoreConfig
	AccessConfig
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
		return Config{}, errors.New(fmt.Sprintf(
			"invalid mode, --mode can only be %s, %s, or %s",
			types.Access,
			types.Prep,
			types.Upload,
		))
	}

	var storeCfg StoreConfig
	var pathCfg PathConfig
	var accessCfg AccessConfig

	// validate sql connection creds
	if mode == types.Access || mode == types.Upload {
		sqlConn := c.String("sqlConn")
		if sqlConn == "" {
			return Config{}, errors.New("invalid sql connection string, --sqlConn should not be empty")
		}
		var user, pass, database, host, port string
		// std lib regex doesn't support lookahead and lookbehind
		var re *regexp2.Regexp
		var m *regexp2.Match
		var err error
		re = regexp2.MustCompile(`(?<=user=).+?(?=\s|$)`, 0)
		m, err = re.FindStringMatch(sqlConn)
		if err != nil || m == nil {
			return Config{}, errors.New("invalid sql connection string, unable to parse 'user' argument")
		}
		user = m.String()
		re = regexp2.MustCompile(`(?<=password=).+?(?=\s|$)`, 0)
		m, err = re.FindStringMatch(sqlConn)
		if err != nil || m == nil {
			return Config{}, errors.New("invalid sql connection string, unable to parse 'password' argument")
		}
		pass = m.String()
		re = regexp2.MustCompile(`(?<=host=).+?(?=\s|$)`, 0)
		m, err = re.FindStringMatch(sqlConn)
		if err != nil || m == nil {
			return Config{}, errors.New("invalid sql connection string, unable to parse 'host' argument")
		}
		host = m.String()
		re = regexp2.MustCompile(`(?<=port=).+?(?=\s|$)`, 0)
		m, err = re.FindStringMatch(sqlConn)
		if err != nil || m == nil {
			return Config{}, errors.New("invalid sql connection string, unable to parse 'port' argument")
		}
		port = m.String()
		re = regexp2.MustCompile(`(?<=database=).+?(?=\s|$)`, 0)
		m, err = re.FindStringMatch(sqlConn)
		if err != nil || m == nil {
			return Config{}, errors.New("invalid sql connection string, unable to parse 'database' argument")
		}
		database = m.String()

		if util.StrContains([]string{sqlConn, user, pass, database, host, port}, "") {
			return Config{}, errors.New("invalid sql connection string, respecify --sqlConn")
		}
		storeCfg = StoreConfig{
			ConnStr: sqlConn,
			Dbuser:  user,
			Dbpass:  pass,
			Dbname:  database,
			Dbhost:  host,
			Dbport:  port,
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

	return Config{
		Mode:         mode,
		PathConfig:   pathCfg,
		StoreConfig:  storeCfg,
		AccessConfig: accessCfg,
	}, nil
}
