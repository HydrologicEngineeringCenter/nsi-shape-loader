package store

import (
	"log"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/usace/goquery"
)

type PSStore struct {
	DS goquery.DataStore
}

func NewStore(c config.Config) (*PSStore, error) {
	dbconf := c.Rdbmsconfig()
	ds, err := goquery.NewRdbmsDataStore(&dbconf)
	if err != nil {
		log.Printf("Unable to connect to database during startup: %s", err)
	} else {
		log.Printf("Connected as %s to database %s:%s/%s", c.Dbuser, c.Dbhost, c.Dbport, c.Dbname)
	}

	st := PSStore{ds}
	return &st, nil
}
