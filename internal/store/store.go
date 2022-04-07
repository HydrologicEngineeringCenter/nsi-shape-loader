package store

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	shape "github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/shp"
	"github.com/google/uuid"
	"github.com/jonas-p/go-shp"
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

func (st *PSStore) AddDomain(d model.Domain) (uuid.UUID, error) {
	var dId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&domainTable).
			Tx(&tx).
			StatementKey("insert").
			Params(d.FieldId, d.Value).
			Dest(&dId).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return dId, err
}

func (st *PSStore) AddField(f model.Field) (uuid.UUID, error) {
	var fId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&fieldTable).
			Tx(&tx).
			StatementKey("insert").
			Params(f.Name, f.Type, f.Description, f.IsDomain).
			Dest(&fId).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return fId, err
}

func (st *PSStore) AddSchemaFieldAssociation(sf model.SchemaField) (uuid.UUID, error) {
	var schemaId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&schemaFieldTable).
			Tx(&tx).
			StatementKey("insert").
			Params(sf.Id, sf.NsiFieldId, sf.IsPrivate).
			Dest(&schemaId).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return schemaId, err
}

func (st *PSStore) AddSchema(schema model.Schema) (uuid.UUID, error) {
	var schemaId uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&schemaTable).
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

func (st *PSStore) AddDataset(d model.Dataset) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(datasetTable.Statements["insertNullShape"]).
		Params(
			d.Name,
			d.Version,
			d.SchemaId,
			d.TableName,
			d.Description,
			d.Purpose,
			d.CreatedBy,
			d.QualityId,
		).
		Dest(&ids).
		Fetch()
	if err != nil {
		panic(err)
	}
	return ids[0], err
}

func (st *PSStore) AddAccess(a model.Access) (uuid.UUID, error) {
	var id uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&accessTable).
			Tx(&tx).
			StatementKey("insert").
			Params(a.DatasetId, a.Group, a.Role, a.Permission).
			Dest(&id).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return id, err
}

func (st *PSStore) AddQuality(q model.Quality) (uuid.UUID, error) {
	var id uuid.UUID
	err := goquery.Transaction(st.DS, func(tx goquery.Tx) {
		err := st.DS.Select().
			DataSet(&domainTable).
			Tx(&tx).
			StatementKey("insert").
			Params(q.Value, q.Description).
			Dest(&id).
			Fetch()
		if err != nil {
			panic(err)
		}
	})
	return id, err
}

func (st *PSStore) GetDomainId(d model.Domain) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(schemaTable.Statements["selectId"]).
		Params(d.FieldId, d.Value).
		Dest(&ids).
		Fetch()
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for domain.field_id=" + d.FieldId.String() + ", domain.value=" + d.Value)
	}
	return ids[0], err
}

func (st *PSStore) GetAccessId(a model.Access) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(schemaTable.Statements["selectId"]).
		Params(a.DatasetId, a.Group).
		Dest(&ids).
		Fetch()
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for access.dataset_id=" + a.DatasetId.String() + " and access.access_group=" + a.Group)
	}
	return ids[0], err
}

func (st *PSStore) GetDatasetId(d model.Dataset) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(datasetTable.Statements["selectId"]).
		Params(d.Name, d.Version, d.Purpose, d.QualityId).
		Dest(&ids).
		Fetch()
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New(fmt.Sprintf(`more than 1 id exists for
        dataset.name=%s
        dataset.version=%s
        dataset.shape=%s
        dataset.purpose=%s
        dataset.quality_id=%s`,
			d.Name,
			d.Version,
			d.Shape,
			d.Purpose,
			d.QualityId,
		))
	}
	return ids[0], err
}

// GetDataset queries based on its Name, Version, Purpose, and QualityId
func (st *PSStore) GetDataset(d *model.Dataset) error {
	var ds []model.Dataset
	err := st.DS.
		Select(datasetTable.Statements["select"]).
		Params(d.Name, d.Version, d.Purpose, d.QualityId).
		Dest(&ds).
		Fetch()
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		return nil
	}
	if len(ds) > 1 {
		return errors.New(fmt.Sprintf(`more than 1 dataset exists for
        dataset.name=%s
        dataset.version=%s
        dataset.purpose=%s
        dataset.quality_id=%s`,
			d.Name,
			d.Version,
			d.Purpose,
			d.QualityId,
		))
	}
	// replace data at location referenced by pointer
	*d = ds[0]
	return err
}

func (st *PSStore) GetFieldId(f model.Field) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&fieldTable).
		StatementKey("select").
		Params(f.Name).
		Dest(&ids).
		Fetch()
	if len(ids) == 0 {
		return uuid.UUID{}, err
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for field.name=" + f.Name + " and field.type=" + string(f.Type))
	}
	return ids[0], err
}

func (st *PSStore) GetSchemaId(s model.Schema) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(schemaTable.Statements["selectId"]).
		Params(s.Name, s.Version).
		Dest(&ids).
		Fetch()
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for schema.name=" + s.Name + " and schema.version=" + s.Version)
	}
	return ids[0], err
}

func (st *PSStore) GetQualityId(q model.Quality) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select(qualityTable.Statements["selectId"]).
		Params(q.Value).
		Dest(&ids).
		Fetch()
	if err != nil {
		return uuid.UUID{}, err
	}
	if len(ids) == 0 {
		return uuid.UUID{}, nil
	}
	if len(ids) > 1 {
		return uuid.UUID{}, errors.New("more than 1 id exists for quality.value=" + string(q.Value))
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

func (st *PSStore) SchemaFieldAssociationExists(sf model.SchemaField) (bool, error) {
	var ids []uuid.UUID
	var result bool
	err := st.DS.
		Select().
		DataSet(&schemaFieldTable).
		StatementKey("selectId").
		Params(sf.Id, sf.NsiFieldId).
		Dest(&ids).
		Fetch()
	if err != nil {
		panic(err)
	}
	if len(ids) > 0 {
		result = true
	} else {
		result = false
	}
	return result, err
}

func (st *PSStore) UpdateDatasetBBox(d model.Dataset) error {
	// hacky way to dynamically generate table_name since identifiers cannot be used as variables
	// should be safe from sql injection since all table names are generated internally from guids
	var ids []interface{}
	err := st.DS.
		Select(strings.ReplaceAll(datasetTable.Statements["updateBBox"], "{table_name}", d.TableName)).
		Params(d.Id).
		Dest(&ids). // interface doesn't work without a dest sink
		Fetch()
	return err
}

// ShpDataInStore checks if shp file has already been uploaded to database
func (st *PSStore) ShpDataInStore(d model.Dataset, s *shp.Reader) (bool, error) {
	// algo takes a set of random sample points, if any sample matches with
	// an entry in the db, return true
	var ids []int
	sampleSize := 50

	xIdx, err := shape.FieldIdx(s, "X")
	if err != nil {
		return false, err
	}
	yIdx, err := shape.FieldIdx(s, "Y")
	if err != nil {
		return false, err
	}

	for i := 0; i < sampleSize; i++ {
		sampleIdx := rand.Int() % s.AttributeCount()
		x := s.ReadAttribute(sampleIdx, xIdx)
		y := s.ReadAttribute(sampleIdx, yIdx)
		err := st.DS.
			Select(strings.ReplaceAll(datasetTable.Statements["structureInInventory"], "{table_name}", d.TableName)).
			Params(x, y).
			Dest(&ids).
			Fetch()
		if err != nil {
			return false, err
		}
		if len(ids) > 0 {
			return true, nil
		}
	}
	return false, nil
}
