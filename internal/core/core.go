package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

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
	//  Check if field is associated to schema
	//      Yes -> reference association
	//      No -> add association
	//  Check if dataset exists
	//      Yes -> check if dataset has referenced schema
	//          Yes -> reference id
	//          No -> panic
	//      No -> create new dataset
	//  Insert inventory table using bash script
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
	schemaId, err = st.GetSchemaId(schema)
	if err != nil {
		return err
	}
	if schemaId == uuid.Nil {
		schemaId, err = st.AddSchema(schema)
		if err != nil {
			panic(err)
		}
	}

	// init data from shp file
	fmt.Println("Reading shapefile from: " + cfg.ShpPath)
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
		fieldDescription, err := xlsF.GetCellValue("field-domain", "D"+fmt.Sprint(j+2))
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
		field := model.Field{
			Name:        f.String(),
			Type:        types.DatatypeReverse[string(f.Fieldtype)],
			Description: fieldDescription,
			IsDomain:    bIsDomain,
		}
		fieldId, err = st.GetFieldId(field)
		if err != nil {
			return err
		}
		if fieldId == uuid.Nil {
			fieldId, err = st.AddField(field)
			if err != nil {
				panic(err)
			}

			///////////////////////////////
			//   DOMAIN
			// Process domain only if specified by field ie. field holds a discrete categorical variable
			// Currently this is specified from the metadata xls, could be a TODO to automatically detect field based only on the shp file
			if bIsDomain {
				fieldVals := shape.UniqueValues(shpf, f)
				if err != nil {
					return err
				}
				for _, v := range fieldVals {
					domain := model.Domain{
						FieldId: fieldId,
						Value:   v,
					}
					// This location can only be reached for new field inserts,
					// can assume that domain do not exists
					_, err := st.AddDomain(domain)
					if err != nil {
						return err
					}
				}
			}
		}
		///////////////////////////////
		//   SCHEMA_FIELD ASSOCIATION
		flagAssociation, err := st.SchemaFieldAssociationExists(schemaId, fieldId)
		if err != nil {
			return err
		}
		if !flagAssociation {
			_, err = st.AddSchemaFieldAssociation(schemaId, fieldId)
			if err != nil {
				return err
			}
		}
	}

	////////////////////////////////////////
	//  DATASET

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

	tableName := "inventory_" + uuid.New().String()
	dataset := model.Dataset{
		Name:        datasetName,
		Version:     datasetVersion,
		SchemaId:    schemaId,
		TableName:   tableName,
		Shape:       types.ShapeReverse[shpf.GeometryType],
		Description: datasetDescription,
		Purpose:     datasetPurpose,
		CreatedBy:   datasetCreatedBy,
		QualityId:   qualityId,
	}
	datasetId, err := st.GetDatasetId(dataset)
	if err != nil {
		return err
	}
	if datasetId == uuid.Nil {
		_, err = st.AddDataset(dataset)
		// Upload data
		// cmd = exec.Command("/bin/bash", "./createtable", "-d", "-s", cfg.StoreConfig.ConnStr, "-c", "-t")
		_, err = exec.Command("/bin/bash", "./createtable", "-s", cfg.StoreConfig.ConnStr, "-c", store.DbSchema, "-t", dataset.TableName).Output()
	} else {
		// dataset already exists
		_, err = exec.Command("/bin/bash", "./appendtable", "-d", "-s", cfg.StoreConfig.ConnStr, "-c", "-t").Output()
	}

	// return string(cmd), err
	return err
}

func ChangeAccess(cfg config.Config) error {
	st, err := store.NewStore(cfg)
	////////////////////////////////////////
	//  ACCESS

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
