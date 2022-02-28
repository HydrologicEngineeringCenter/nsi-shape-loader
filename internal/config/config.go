package config

import (
	"errors"

	"github.com/urfave/cli/v2"
	dq "github.com/usace/goquery"
)

type Config struct {
	Dbuser    string
	Dbpass    string
	Dbname    string
	Dbhost    string
	Dbstore   string
	Dbdriver  string
	DBSSLMode string
	Dbport    string
}

func (c *Config) Rdbmsconfig() dq.RdbmsConfig {
	return dq.RdbmsConfig{
		Dbuser:   c.Dbuser,
		Dbpass:   c.Dbpass,
		Dbhost:   c.Dbhost,
		Dbport:   c.Dbport,
		Dbname:   c.Dbname,
		DbDriver: c.Dbdriver,
		DbStore:  c.Dbstore,
	}
}

// NewConfig generates new config from cli args context
func NewConfig(c *cli.Context) (Config, error) {
	if c.NArg() < 0 {
		return Config{}, errors.New("newconfig: not enough arguments")
	}
	return Config{}, nil
}
