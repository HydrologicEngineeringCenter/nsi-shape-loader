package config

import (
	"errors"

	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

type Config struct {
	Dbuser      string
	Dbpass      string
	Dbname      string
	Dbtablename string
	Dbhost      string
	Dbport      string
	FilePath    string
}

func (c *Config) Rdbmsconfig() dq.RdbmsConfig {
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
	if c.NumFlags() < 5 {
		return Config{}, errors.New("newconfig: not enough input flags")
	}
	return Config{
		Dbuser:      c.String("dbuser"),
		Dbpass:      c.String("dbpass"),
		Dbhost:      c.String("dbhost"),
		Dbport:      c.String("dbport"),
		Dbname:      c.String("dbname"),
		Dbtablename: c.String("dbtname"),
		FilePath:    c.String("filepath"),
	}, nil
}
