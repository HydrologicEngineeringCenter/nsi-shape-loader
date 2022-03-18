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
	return nil
}

func (st *PSStore) AddField(f model.Field) (uuid.UUID, error) {
	var fId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&fieldTable).
			Tx(&tx).
			StatementKey("insert").
			Params().
			Dest(&fId).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return fId, err
}

func (st *PSStore) AddSchema(schema model.Schema) (uuid.UUID, error) {
	var schemaId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&domainTable).
			Tx(&tx).
			StatementKey("insert").
			Params(schema.Name, schema.Version, schema.Notes).
			Dest(&schemaId).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return schemaId, err
}

func (st *PSStore) AddDataset(schema model.Schema) error {
	return nil
}

func (st *PSStore) AddAccess(schema model.Schema) error {
	return nil
}

func (st *PSStore) AddQuality(schema model.Schema) error {
	return nil
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
		return uuid.UUID{}, errors.New("more than 1 id exists for schema.name=" + s.Name + " and schema.version=" + s.Version)
	}
	return ids[0], err
}

func (st *PSStore) GetFieldId(f model.Field) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(schemaTable.Statements["select"]).
		Params(f.Name).
		Dest(&ids).
		Fetch()
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for field.name=" + f.Name)
	}
	return ids[0], err
}

// Check if table exists in database
func (st *PSStore) TableExists(schema string, table string) (bool, error) {
	var result bool
	err := st.DS.Select(`
    SELECT EXISTS (
        SELECT FROM pg_tables
        WHERE
            schemaname='$1' AND
            tablename='$2'
    )
    `).
		Params(schema, table).
		Dest(&result).
		Fetch()
	return result, err
}

func (st *PSStore) SchemaFieldAssociationExists(fieldId uuid.UUID, schemaID uuid.UUID) (bool, error) {
	var result bool
	var association []model.SchemaField
	err := st.DS.Select(`
    SELECT * FROM schema_field
    WHERE id=$1 AND field_id=$2
    `).
		Params(schema, table).
		Dest(&association).
		Fetch()
	if len(association) > 0 {
		result = true
	} else {
		result = false
	}
	return result, err
}
