package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/files"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/model"
	shape "github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/shp"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/xls"
	"github.com/google/uuid"
	"github.com/jonas-p/go-shp"
	"github.com/urfave/cli/v2"
)

func Core(c *cli.Context) error {
	//  pre - generate config xls from shp
	//  upload - upload based on data and metadata from xls and shp
	//  access - change access group and role
	cfg, err := config.NewConfig(c)
	if err != nil {
		panic(err)
	}

	// Prep mode generates the metadata xls required by Upload
	if cfg.Mode == types.Prep {
		err = Prep(cfg)
	}
	// Upload mode uploads the dataset and associated metadata
	if cfg.Mode == types.Upload {
		err = Upload(cfg)
	}
	// Access mode change access permission of groups
	if cfg.Mode == types.Access {
		err = ChangeAccess(cfg)
	}
	return err
}

// PreUpload generates an xls template from shp file fields
func Prep(cfg config.Config) error {

	// copy xls file
	const baseXlsSrc = "./assets/baseMetadata.xlsx"
	const cpXlsDest = "./metadata.xlsx"
	err := files.Copy(baseXlsSrc, cpXlsDest)
	if err != nil {
		return err
	}

	xlsF, err := xls.NewXls(cpXlsDest)
	if err != nil {
		return err
	}
	shpf, err := shape.NewShp(cfg.ShpPath)
	if err != nil {
		return err
	}
	fields := shpf.Fields()
	var loc, val string
	for j, f := range fields {
		loc = "B" + fmt.Sprint(j+2)
		val = f.String()
		err = xlsF.SetCellValue("field-domain", loc, val)
		if err != nil {
			return err
		}
	}
	xlsF.Save()
	wd, err := os.Getwd()
	fmt.Println("Metadata template file successfully created at:", filepath.Join(wd, cpXlsDest))
	return err
}

//// Upload populates metadata from the config xls and upload data from shp file
func Upload(cfg config.Config) error {
	st, err := store.NewStore(cfg)
	if err != nil {
		return err
	}

	// fields := shpf.Fields()

	// var insertTable = dq.TableDataSet{
	// 	Name:       cfg.Dbtablename,
	// 	Statements: map[string]string{},
	// 	Fields:     model.Point{},
	// }

	// batchSize := 20000
	// lastRecordIdx := shpf.AttributeCount() - 1
	// var records []model.Point
	// for shpf.Next() {
	// 	i, _ := shpf.Shape()

	// 	// construct data struct from point
	// 	var newPoint model.Point
	// 	for j, f := range fields {
	// 		val := shpf.ReadAttribute(i, j)
	// 		fieldStr := strings.Title(strings.ToLower(f.String()))
	// 		structutil.SetField(&newPoint, fieldStr, val)
	// 	}
	// 	records = append(records, newPoint)
	// 	// Batch upload on reaching batchsize limit
	// 	if (i != 0 && i%batchSize == 0) || i == lastRecordIdx {
	// 		if i == lastRecordIdx {
	// 			// batching the last odd lot records doesn't work for some reason
	// 			err = st.DS.Insert(&insertTable).
	// 				Records(&records).
	// 				Execute()
	// 		} else {
	// 			err = st.DS.Insert(&insertTable).
	// 				Records(&records).
	// 				Batch(true).
	// 				BatchSize(len(records)).
	// 				Execute()
	// 		}
	// 		if err != nil {
	// 			return err
	// 		}
	// 		fmt.Println("Proccessed " + fmt.Sprint(i+1) + " records")
	// 		records = []model.Point{}
	// 	}
	// }
	// fmt.Println("Processing finished.")

	/////////////////////////////////////////////////////////
	// Data insertion procedure:
	//  Check if schema exists based on unique(name, version)
	//      Yes -> reference id from existing schema
	//      No -> insert new schema into store
	//  Check if field exists based on unique(name, type)
	//      Yes -> reference id for existing field
	//      No ->
	//          insert new field into store
	//          Required domain?
	//              Yes -> insert new domain referencing field
	//              No -> continue
	//  Check if field is already associated to schema
	//      Yes -> reference association
	//      No -> add association
	//  Check if dataset exists
	//      Yes -> check if dataset has referenced schema
	//          Yes -> reference id
	//          No -> panic
	//      No -> create new dataset
	//  Insert inventory table using ogr2ogr
	log.Printf("Reading metadata from: %s\n", cfg.XlsPath)
	xlsF, err := xls.NewXls(cfg.XlsPath)
	if err != nil {
		return err
	}

	/////////////////////////////////////////////////
	//  SCHEMA
	var schemaId uuid.UUID
	schemaName, err := xlsF.GetCellValue("schema", "C1")
	schemaVersion, err := xlsF.GetCellValue("schema", "C2")
	schemaNotes, err := xlsF.GetCellValue("schema", "C3")

	schema := model.Schema{
		Name:    schemaName,
		Version: schemaVersion,
		Notes:   schemaNotes,
	}
	log.Printf("Retrieving id for unique schema=%s version=%s\n", schema.Name, schema.Version)
	schemaId, err = st.GetSchemaId(schema)
	if err != nil {
		return err
	}
	if schemaId == uuid.Nil {
		log.Printf("schema=%s version=%s do not exists. Adding to schema table...\n", schema.Name, schema.Version)
		schemaId, err = st.AddSchema(schema)
		if err != nil {
			panic(err)
		}
	}
	schema.Id = schemaId

	// init data from shp file
	log.Printf("Reading shapefile from: %s\n", cfg.ShpPath)
	shpf, err := shp.Open(cfg.ShpPath)
	if err != nil {
		return err
	}
	defer shpf.Close()

	fields := shpf.Fields()
	for j, f := range fields {

		///////////////////////////////
		//   FIELD
		var fieldId uuid.UUID
		fieldDescription, err := xlsF.GetCellValue("field-domain", "E"+fmt.Sprint(j+2))
		if err != nil {
			return err
		}
		sIsDomain, err := xlsF.GetCellValue("field-domain", "C"+fmt.Sprint(j+2))
		if err != nil {
			return err
		}
		bIsDomain, err := strconv.ParseBool(sIsDomain)
		if err != nil {
			return err
		}
		sPrivate, err := xlsF.GetCellValue("field-domain", "D"+fmt.Sprint(j+2))
		if err != nil {
			return err
		}
		bPrivate, err := strconv.ParseBool(sPrivate)
		if err != nil {
			return err
		}
		field := model.Field{
			Name:        f.String(),
			Type:        types.DatatypeReverse[string(f.Fieldtype)],
			Description: fieldDescription,
			IsDomain:    bIsDomain,
		}
		log.Printf("Retrieving id for unique field=%s type=%s\n", field.Name, field.Type)
		fieldId, err = st.GetFieldId(field)
		if err != nil {
			return err
		}
		// If no id -> field is not in db -> add field + add association to schema + domain
		if fieldId != uuid.Nil {
			field.Id = fieldId
		} else {
			log.Printf("field=%s type=%s do not exists. Adding to field table...\n", field.Name, field.Type)
			fieldId, err = st.AddField(field)
			if err != nil {
				panic(err)
			}
			field.Id = fieldId
			///////////////////////////////
			//   DOMAIN
			// Process domain only if specified by field ie. field holds a discrete categorical variable
			// Currently this is specified from the metadata xls, could be a TODO to automatically detect field based only on the shp file
			if bIsDomain {
				log.Printf("field=%s holds discrete categorical variables. Adding to domain table...\n", field.Name)
				fieldVals, err := shape.UniqueValues(shpf, f)
				if err != nil {
					return err
				}
				for _, v := range fieldVals {
					domain := model.Domain{
						FieldId: fieldId,
						Value:   v,
					}
					// This location can only be reached for new field inserts,
					// can assume that domain, has not yet exists
					domainId, err := st.AddDomain(domain)
					if err != nil {
						return err
					}
					domain.Id = domainId
				}
			}
		}
		///////////////////////////////
		//   SCHEMA_FIELD_ASSOCIATION check for both cases - field already exists or new insert
		//      the same field can be associated to multiple schemas
		sf := model.SchemaField{
			Id:         schemaId,
			NsiFieldId: fieldId,
			IsPrivate:  bPrivate,
		}
		flagAssociation, err := st.SchemaFieldAssociationExists(sf)
		if err != nil {
			return err
		}
		if !flagAssociation {
			log.Printf("Unable to find association between schema=%s and field=%s. Adding to schema_field table...\n", schema.Name, field.Name)
			_, err = st.AddSchemaFieldAssociation(sf)
			if err != nil {
				return err
			}
		}
	}

	////////////////////////////////////////
	//  DATASET comes after the fields loop
	datasetName, err := xlsF.GetCellValue("dataset", "C1")
	if err != nil {
		return err
	}
	datasetVersion, err := xlsF.GetCellValue("dataset", "C2")
	if err != nil {
		return err
	}
	datasetDescription, err := xlsF.GetCellValue("dataset", "C3")
	if err != nil {
		return err
	}
	datasetPurpose, err := xlsF.GetCellValue("dataset", "C4")
	if err != nil {
		return err
	}
	datasetCreatedBy, err := xlsF.GetCellValue("dataset", "C5")
	if err != nil {
		return err
	}
	sDatasetQuality, err := xlsF.GetCellValue("dataset", "C6")
	if err != nil {
		return err
	}

	// TODO validate this quality input
	var qualityId uuid.UUID
	quality := model.Quality{
		Value: types.QualityReverse[sDatasetQuality],
	}
	qualityId, err = st.GetQualityId(quality)
	if err != nil {
		return err
	}
	if qualityId == uuid.Nil {
		qualityId, err = st.AddQuality(quality)
	}

	dataset := model.Dataset{
		Name:        datasetName,
		Version:     datasetVersion,
		SchemaId:    schemaId,
		Description: datasetDescription,
		Purpose:     datasetPurpose,
		CreatedBy:   datasetCreatedBy,
		QualityId:   qualityId,
	}
	err = st.GetDataset(&dataset)
	if err != nil {
		return err
	}
	var datasetId uuid.UUID
	if dataset.Id == uuid.Nil {
		tableName := "inventory_" + strings.ReplaceAll(uuid.New().String(), "-", "_")
		dataset.TableName = tableName
		datasetId, err = st.AddDataset(dataset)
		if err != nil {
			return err
		}
		dataset.Id = datasetId
		// create new table
		log.Printf("Creating table=%s for dataset=%s", dataset.TableName, dataset.Name)
		execStr := fmt.Sprintf(`ogr2ogr -f "PostgreSQL" PG:"%s" %s -lco precision=no -lco fid=fd_id -lco geometry_name=shape -nln %s.%s`,
			strings.ReplaceAll(cfg.StoreConfig.ConnStr, "database=", "dbname="),
			cfg.ShpPath, store.DbSchema, dataset.TableName,
		)
		fmt.Println(execStr)
		_, err = exec.Command(
			"sh", "-c", execStr,
		).Output()
	} else {
		// dataset already exists
		flagDataInStore, err := st.ShpDataInStore(dataset, shpf)
		if err != nil {
			return err
		}
		if !flagDataInStore { // data has not yet been added to store
			log.Printf("table=%s exists for dataset=%s. Appending rows...", dataset.TableName, dataset.Name)
			execStr := fmt.Sprintf(`ogr2ogr -append -update -f "PostgreSQL" PG:"%s" %s -lco precision=no -nln %s.%s`,
				strings.ReplaceAll(cfg.StoreConfig.ConnStr, "database=", "dbname="),
				cfg.ShpPath, store.DbSchema, dataset.TableName,
			)
			_, err = exec.Command(
				"sh", "-c", execStr,
			).Output()
		} else {
			return errors.New("shp file has already been uploaded")
		}
	}
	if err != nil {
		panic(err)
	} else {
		err = st.UpdateDatasetBBox(dataset)
		if err != nil {
			return err
		}
		log.Printf("Data uploaded to dataset.name= %s dataset.id=%s", dataset.Name, dataset.Id)
	}
	return err
}

func ChangeAccess(cfg config.Config) error {
	st, err := store.NewStore(cfg)
	var accessId uuid.UUID
	access := model.Access{
		DatasetId:  cfg.AccessConfig.DatasetId,
		Group:      cfg.AccessConfig.Group,
		Role:       cfg.AccessConfig.Role,
		Permission: types.RolePermission[cfg.AccessConfig.Role],
	}
	accessId, err = st.GetAccessId(access)
	if err != nil {
		return err
	}
	if accessId == uuid.Nil {
		_, err = st.AddAccess(access)
	}

	return err
}

func AddElevation(cfg config.Config) error {
	return nil
}
