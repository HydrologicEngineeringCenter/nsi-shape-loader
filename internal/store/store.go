package store

import (
	"errors"
	"log"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	"github.com/google/uuid"
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

func (st *PSStore) AddDomain(schema model.Schema) error {
}

func (st *PSStore) AddField(schema model.Schema) error {
}

func (st *PSStore) AddSchema(schema model.Schema) error {
	return s.DS.Exec(goquery.NoTx, domainTable.Statements["insert"], schema.Name, schema.Version, schema.Notes)
}

func (st *PSStore) AddDataset(schema model.Schema) error {
}

func (st *PSStore) AddAccess(schema model.Schema) error {
}

func (st *PSStore) AddQuality(schema model.Schema) error {
}

func (st *PSStore) GetSchemaId(s model.Schema) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(schemaTable.Statements["selectId"]).
		Params(s.Name, s.Version).
		Dest(&ids).
		Fetch()
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for schema=" + s.Name + " and version=" + s.Version)
	}
	return ids[0], err
}
