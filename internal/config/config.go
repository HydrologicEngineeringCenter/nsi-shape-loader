package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

type Mode string

const (
	Pre    Mode = "P"
	Upload Mode = "U"
)

var (
	ModeReverse = map[string]Mode{
		"P": Pre,
		"U": Upload,
	}
)

// StoreConfig holds only params required for database connection
type StoreConfig struct {
	Dbuser string
	Dbpass string
	Dbname string
	Dbhost string
	Dbport string
}

// Config is for the general app
type Config struct {
	ShpPath string
	XlsPath string // excel file that maps field to description
	Mode
	StoreConfig
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

	// Init StoreConfig from env variables
	var storeCfg StoreConfig
	if err := envconfig.Process("", &storeCfg); err != nil {
		return Config{}, err
	}

	mode := ModeReverse[c.String("mode")]
	if mode == Upload && c.NumFlags() < 6 {
		return Config{}, errors.New("newconfig: not enough input flags")
	}
	return Config{
		Mode:        ModeReverse[c.String("mode")],
		StoreConfig: storeCfg,
		ShpPath:     c.String("shppath"),
		XlsPath:     c.String("xlsmeta"),
	}, nil
}
