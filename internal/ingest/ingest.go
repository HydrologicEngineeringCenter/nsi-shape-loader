package ingest

import (
	"fmt"
	"log"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/model"
	shape "github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/shp"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/xls"
	"github.com/jonas-p/go-shp"
	"github.com/usace/xlscellreader"
)

func NewMetaAccessor(c config.Config) (MetaAccessor, error) {
	log.Printf("Reading metadata from: %s\n", c.XlsPath)
	xlsF, err := xls.NewXls(c.XlsPath)
	if err != nil {
		return MetaAccessor{}, err
	}
	// init data from shp file
	log.Printf("Reading shapefile from: %s\n", c.ShpPath)
	shpf, err := shp.Open(c.ShpPath)
	if err != nil {
		return MetaAccessor{}, err
	}
	defer shpf.Close()

	meta := MetaAccessor{
		X: xlsF,
		S: shpf,
	}
	return meta, nil
}

// MetaAccessor wraps around the xls and shp readers and acts as a Data Access
// Object. It is similar to the store, only the store is used to access the
// PostGIS database. If a method requires store access, the store Data Access
// Object must be passed explicitly as an argument.
type MetaAccessor struct {
	X *xlscellreader.CellReader // xls reader
	S *shp.Reader               // shp reader
}

func (a MetaAccessor) GetSchema() (model.Schema, error) {
	schemaName, err := a.X.GetString("schema", "C1")
	if err != nil {
		return model.Schema{}, err
	}
	schemaVersion, err := a.X.GetString("schema", "C2")
	if err != nil {
		return model.Schema{}, err
	}
	schemaNotes, err := a.X.GetString("schema", "C3")
	if err != nil {
		return model.Schema{}, err
	}
	schema := model.Schema{
		Name:    schemaName,
		Version: schemaVersion,
		Notes:   schemaNotes,
	}
	return schema, err
}

// GetShpDbFieldNameMap maps shp field name to db field name. A filter is applied
// to keep only fields indicated in the xls
func (a MetaAccessor) GetShpDbFieldNameMap() (map[string]string, error) {
	shp2DbName := map[string]string{} // map field name from shp file to db table col name
	var err error
	fieldsShape := a.S.Fields()
	fieldsModel, err := a.GetFields()
	if err != nil {
		return map[string]string{}, err
	}
	for j, f := range fieldsModel {
		if f.IsInDb {
			fSName := fieldsShape[j].String()
			fXName, err := a.X.GetString("field-domain", "F"+fmt.Sprint(j+2))
			if err != nil {
				return map[string]string{}, err
			}
			shp2DbName[fSName] = fXName
		}
	}
	return shp2DbName, nil
}

func (a MetaAccessor) GetFields() ([]model.Field, error) {
	var fieldsModel []model.Field
	fields := a.S.Fields()
	for j, f := range fields {
		fieldDescription, err := a.X.GetString("field-domain", "G"+fmt.Sprint(j+2))
		if err != nil {
			return []model.Field{}, err
		}
		isDomain, err := a.X.GetBool("field-domain", "D"+fmt.Sprint(j+2))
		if err != nil {
			return []model.Field{}, err
		}
		isInDb, err := a.X.GetBool("field-domain", "C"+fmt.Sprint(j+2))
		if err != nil {
			return []model.Field{}, err
		}
		shpName, err := a.X.GetString("field-domain", "B"+fmt.Sprint(j+2))
		if err != nil {
			return []model.Field{}, err
		}
		dbName, err := a.X.GetString("field-domain", "F"+fmt.Sprint(j+2))
		if err != nil {
			return []model.Field{}, err
		}
		if isInDb {
			field := model.Field{
				ShpName:     shpName,
				DbName:      dbName,
				Type:        types.DatatypeReverse[string(f.Fieldtype)],
				Description: fieldDescription,
				IsDomain:    isDomain,
				IsInDb:      isInDb,
			}
			fieldsModel = append(fieldsModel, field)
		}
	}
	return fieldsModel, nil
}

func (a MetaAccessor) GetGroup() (model.Group, error) {
	groupName, err := a.X.GetString("dataset", "C7")
	if err != nil {
		return model.Group{}, err
	}
	g := model.Group{
		Name: groupName,
	}
	return g, nil
}

func (a MetaAccessor) GetDomainsForField(f model.Field) ([]model.Domain, error) {
	idx, err := shape.FieldIdx(a.S, f.ShpName)
	if err != nil {
		return []model.Domain{}, err
	}
	fieldShape, err := shape.Field(a.S, idx)
	if err != nil {
		return []model.Domain{}, err
	}
	vals, err := shape.UniqueValues(a.S, fieldShape)
	var domains []model.Domain
	for _, val := range vals {
		d := model.Domain{
			FieldId: f.Id,
			Value:   val,
		}
		domains = append(domains, d)
	}
	return domains, nil
}

func (a MetaAccessor) GetDataset(s *store.PSStore, schema model.Schema, g model.Group) (model.Dataset, error) {
	datasetName, err := a.X.GetString("dataset", "C1")
	if err != nil {
		return model.Dataset{}, err
	}
	datasetVersion, err := a.X.GetString("dataset", "C2")
	if err != nil {
		return model.Dataset{}, err
	}
	datasetDescription, err := a.X.GetString("dataset", "C3")
	if err != nil {
		return model.Dataset{}, err
	}
	datasetPurpose, err := a.X.GetString("dataset", "C4")
	if err != nil {
		return model.Dataset{}, err
	}
	datasetCreatedBy, err := a.X.GetString("dataset", "C5")
	if err != nil {
		return model.Dataset{}, err
	}
	q, err := a.GetQuality(s)
	if err != nil {
		return model.Dataset{}, err
	}
	dataset := model.Dataset{
		Name:        datasetName,
		Version:     datasetVersion,
		SchemaId:    schema.Id,
		Description: datasetDescription,
		Purpose:     datasetPurpose,
		CreatedBy:   datasetCreatedBy,
		GroupId:     g.Id,
		QualityId:   q.Id,
	}
	return dataset, nil
}

func (a MetaAccessor) GetQuality(s *store.PSStore) (model.Quality, error) {
	qs, err := a.X.GetString("dataset", "C6")
	if err != nil {
		return model.Quality{}, err
	}
	q := model.Quality{
		Value: types.QualityReverse[qs],
	}
	err = s.GetQuality(&q)
	if err != nil {
		return model.Quality{}, err
	}
	return q, nil
}

func (a MetaAccessor) GetSchemaFieldAssociation(s model.Schema, f model.Field) (model.SchemaField, error) {
	fIdx, err := shape.FieldIdx(a.S, f.ShpName)
	if err != nil {
		return model.SchemaField{}, err
	}
	isPrivate, err := a.X.GetBool("field-domain", "E"+fmt.Sprint(fIdx+2))
	if err != nil {
		return model.SchemaField{}, err
	}
	sf := model.SchemaField{
		Id:         s.Id,
		NsiFieldId: f.Id,
		IsPrivate:  isPrivate,
	}
	return sf, nil
}
