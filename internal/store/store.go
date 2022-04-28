package store

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/global"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/model"
	shape "github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/shp"
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

func (st *PSStore) AddDomain(d *model.Domain) error {
	var dId uuid.UUID
	err := st.DS.Select().
		DataSet(&domainTable).
		StatementKey("insert").
		Params(d.FieldId, d.Value).
		Dest(&dId).
		Fetch()
	if err != nil {
		return err
	}
	d.Id = dId
	return nil
}

func (st *PSStore) AddField(f *model.Field) error {
	var fId uuid.UUID
	err := st.DS.Select().
		DataSet(&fieldTable).
		StatementKey("insert").
		Params(f.DbName, f.Type, f.Description, f.IsDomain).
		Dest(&fId).
		Fetch()
	if err != nil {
		return err
	}
	f.Id = fId
	return nil
}

func (st *PSStore) AddMember(m *model.Member) error {
	var mId uuid.UUID
	err := st.DS.Select().
		DataSet(&memberTable).
		StatementKey("insert").
		Params(m.GroupId, m.Role, m.UserId).
		Dest(&mId).
		Fetch()
	if err != nil {
		return err
	}
	m.Id = mId
	return nil
}

func (st *PSStore) AddSchemaFieldAssociation(sf model.SchemaField) error {
	var schemaId uuid.UUID
	err := st.DS.Select().
		DataSet(&schemaFieldTable).
		StatementKey("insert").
		Params(sf.Id, sf.NsiFieldId, sf.IsPrivate).
		Dest(&schemaId).
		Fetch()
	if err != nil {
		return err
	}
	return nil
}

func (st *PSStore) AddSchema(schema *model.Schema) error {
	var schemaId uuid.UUID
	err := st.DS.Select().
		DataSet(&schemaTable).
		StatementKey("insert").
		Params(schema.Name, schema.Version, schema.Notes).
		Dest(&schemaId).
		Fetch()
	if err != nil {
		return err
	}
	schema.Id = schemaId
	return err
}

func (st *PSStore) AddDataset(d *model.Dataset) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&datasetTable).
		StatementKey("insertNullShape").
		Params(
			d.Name,
			d.Version,
			d.SchemaId,
			d.TableName,
			d.Description,
			d.Purpose,
			d.CreatedBy,
			d.QualityId,
			d.GroupId,
		).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	d.Id = ids[0]
	return err
}

func (st *PSStore) AddGroup(g *model.Group) error {
	var id uuid.UUID
	err := st.DS.Select().
		DataSet(&groupTable).
		StatementKey("insert").
		Params(g.Name).
		Dest(&id).
		Fetch()
	if err != nil {
		return err
	}
	g.Id = id
	return err
}

func (st *PSStore) GetDomainId(d model.Domain) (uuid.UUID, error) {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&schemaTable).
		StatementKey("selectId").
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

func (st *PSStore) GetGroupId(g *model.Group) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&groupTable).
		StatementKey("selectId").
		Params(g.Name).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if len(ids) > 1 {
		return errors.New("more than 1 id exists for group.name=" + g.Name)
	}
	g.Id = ids[0]
	return nil
}

func (st *PSStore) GetMemberId(m *model.Member) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&memberTable).
		StatementKey("selectId").
		Params(m.GroupId, m.UserId).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if len(ids) > 1 {
		return errors.New(fmt.Sprintf("more than 1 id exists for group_member.group_id=%s and group_member.user_id=%s", m.GroupId.String(), m.UserId))
	}
	m.Id = ids[0]
	return nil
}

func (st *PSStore) GetDatasetId(d *model.Dataset) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&datasetTable).
		StatementKey("selectId").
		Params(d.Name, d.Version, d.Purpose, d.QualityId).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if len(ids) > 1 {
		return errors.New(fmt.Sprintf(`more than 1 id exists for
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
	d.Id = ids[0]
	return err
}

// GetDataset queries based on its Name, Version, Purpose, and QualityId
func (st *PSStore) GetDataset(d *model.Dataset) error {
	var ds []model.Dataset
	err := st.DS.
		Select().
		DataSet(&datasetTable).
		StatementKey("select").
		Params(d.Name, d.Version, d.QualityId).
		Dest(&ds).
		Fetch()
	if err != nil {
		return err
	}
	if len(ds) == 0 {
		d.Id = uuid.Nil
	} else {
		*d = ds[0]
	}
	return nil
}

func (st *PSStore) GetFieldId(f *model.Field) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&fieldTable).
		StatementKey("select").
		Params(f.DbName).
		Dest(&ids).
		Fetch()
	if len(ids) == 0 {
		f.Id = uuid.Nil
		return err
	}
	if len(ids) > 1 {
		return errors.New("more than 1 id exists for field.name=" + f.DbName + " and field.type=" + string(f.Type))
	}
	f.Id = ids[0]
	return err
}

// GetSchemaId queries the database based on the supplied schema name and version.
// Replaces Id field if a corresponding entry exists, otherwise change Id field to uuid.Nil
func (st *PSStore) GetSchemaId(s *model.Schema) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&schemaTable).
		StatementKey("selectId").
		Params(s.Name, s.Version).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		s.Id = uuid.Nil
		return nil
	}
	if len(ids) > 1 {
		return errors.New("more than 1 id exists for schema.name=" + s.Name + " and schema.version=" + s.Version)
	}
	s.Id = ids[0]
	return nil
}

func (st *PSStore) GetQuality(q *model.Quality) error {
	var qDb model.Quality
	err := st.DS.
		Select().
		DataSet(&qualityTable).
		StatementKey("select").
		Params(q.Value).
		Dest(&qDb).
		Fetch()
	if err != nil {
		return err
	}
	*q = qDb
	return nil
}

func (st *PSStore) GetQualityId(q *model.Quality) error {
	var ids []uuid.UUID
	err := st.DS.
		Select().
		DataSet(&qualityTable).
		StatementKey("selectId").
		Params(q.Value).
		Dest(&ids).
		Fetch()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	if len(ids) > 1 {
		return errors.New("more than 1 id exists for quality.value=" + string(q.Value))
	}
	q.Id = ids[0]
	return nil
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
		return false, err
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

func (st *PSStore) UpdateMemberRole(m *model.Member) error {
	var ids []interface{}
	err := st.DS.
		Select().
		DataSet(&memberTable).
		StatementKey("updateRole").
		Params(m.Id, m.Role).
		Dest(&ids). // interface doesn't work without a dest sink
		Fetch()
	return err
}

// ElevationColumnExists tests if elevation column exists for inventory table
func (st *PSStore) ElevationColumnExists(d model.Dataset) (bool, error) {
	var res bool
	err := st.DS.
		Select(datasetTable.Statements["elevationColumnExists"]).
		Params(global.DB_SCHEMA, d.TableName, global.ELEVATION_COLUMN_NAME).
		Dest(&res).
		Fetch()
	if err != nil {
		return false, err
	}
	return res, nil
}

//////////////////////////////////////////////////
// Add row to sql table generically
//////////////////////////////////////////////////

// appModels contrainsts generic database access to only a list of structs
type appModels interface {
	model.Field | model.Schema | model.Domain | model.SchemaField | model.Dataset | model.Group | model.Member | model.Quality
}

func getAppModelId[T appModels](m *T) uuid.UUID {
	return reflect.ValueOf(*m).FieldByName("Id").Interface().(uuid.UUID)
}

func setAppModelId[T appModels](m *T, id uuid.UUID) {
	reflect.ValueOf(*m).FieldByName("Id").SetBytes(id[:])
}

type insertConfig struct {
	StatementKey string
	FieldOrder   []string
	QueryTable   *goquery.TableDataSet
}

var (
	// this mapper should be app specific
	insertConfigMapper = map[reflect.Type]insertConfig{
		// quality should not be insertable
		reflect.TypeOf(model.Dataset{}): {
			StatementKey: "insertNullShape",
			FieldOrder: []string{
				"Name", "Version", "SchemaId", "TableName", "Description", "Purpose", "CreatedBy", "QualityId", "GroupId",
			},
			QueryTable: &datasetTable,
		},
		reflect.TypeOf(model.Domain{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"FieldId", "Value",
			},
			QueryTable: &domainTable,
		},
		reflect.TypeOf(model.Field{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"DbName", "Type", "Description", "IsDomain",
			},
			QueryTable: &fieldTable,
		},
		reflect.TypeOf(model.Schema{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"Name", "Version", "Notes",
			},
			QueryTable: &schemaTable,
		},
		reflect.TypeOf(model.SchemaField{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"NsiFieldId", "IsPrivate",
			},
			QueryTable: &schemaFieldTable,
		},
		reflect.TypeOf(model.Group{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"Name",
			},
			QueryTable: &groupTable,
		},
		reflect.TypeOf(model.Member{}): {
			StatementKey: "insert",
			FieldOrder: []string{
				"GroupId", "Role", "UserId",
			},
			QueryTable: &memberTable,
		},
	}
)

// AddRow adds row to table based on a model struct
func AddRow[T appModels](st *PSStore, m *T) error {
	var params []interface{}
	modelType := reflect.TypeOf(*m)
	cfg := insertConfigMapper[modelType]
	// loop over all insertable fields
	for _, f := range cfg.FieldOrder {
		fieldVal := reflect.ValueOf(*m).FieldByName(f)
		valKind := fieldVal.Kind()
		switch valKind {
		case reflect.Bool:
			params = append(params, fieldVal.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			params = append(params, fieldVal.Int())
		case reflect.Float32, reflect.Float64:
			params = append(params, fieldVal.Float())
		case reflect.String:
			params = append(params, fieldVal.String())
		case reflect.TypeOf(uuid.Nil).Kind():
			// uuid  Kind() is Ox17 Array, potentially not safe
			params = append(params, fieldVal.Interface().(uuid.UUID))
		default:
			return errors.New(fmt.Sprintf("Generic AddRow does not support param of type: %s", valKind))
		}
	}

	var id uuid.UUID
	err := st.DS.
		Select().
		DataSet(cfg.QueryTable).
		StatementKey(cfg.StatementKey).
		Params(params...).
		Dest(&id).
		Fetch()
	if err != nil {
		return err
	}
	if getAppModelId(m) == uuid.Nil && id != uuid.Nil {
		setAppModelId(m, id)
	}
	return nil
}
