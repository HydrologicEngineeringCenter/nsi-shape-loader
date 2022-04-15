package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/files"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/ingest"
	shape "github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/shp"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/xls"
	"github.com/google/uuid"
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
		err = xlsF.F.SetCellValue("field-domain", loc, val)
		if err != nil {
			return err
		}
	}
	xlsF.F.Save()
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
	metaAccessor, err := ingest.NewMetaAccessor(cfg)
	if err != nil {
		return err
	}
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

	/////////////////////////////////////////////////
	//  SCHEMA
	s, err := metaAccessor.GetSchema()
	if err != nil {
		return err
	}
	log.Printf("Retrieving id for unique schema=%s version=%s\n", s.Name, s.Version)
	err = st.GetSchemaId(&s)
	if err != nil {
		return err
	}
	if s.Id == uuid.Nil {
		log.Printf("schema=%s version=%s do not exists. Adding to schema table...\n", s.Name, s.Version)
		err = st.AddSchema(&s)
		if err != nil {
			return err
		}
	}

	fields, err := metaAccessor.GetFields()
	if err != nil {
		return err
	}
	shp2DbName, err := metaAccessor.GetShpDbFieldNameMap()
	if err != nil {
		return err
	}

	for _, f := range fields {
		log.Printf("Retrieving id for unique field=%s type=%s\n", f.Name, f.Type)
		err = st.GetFieldId(&f)
		if err != nil {
			return err
		}
		// If no id -> field is not in db -> add field + add association to schema + domain
		if f.Id == uuid.Nil {
			log.Printf("field=%s type=%s do not exists. Adding to field table...\n", f.Name, f.Type)
			err = st.AddField(&f)
			if err != nil {
				return err
			}
			///////////////////////////////
			//   DOMAIN
			// Process domain only if specified by field ie. field holds a discrete categorical variable
			// Currently this is specified from the metadata xls, could be a TODO to automatically detect field based only on the shp file
			if f.IsDomain {
				log.Printf("field=%s holds discrete categorical variables. Adding to domain table...\n", f.Name)
				if err != nil {
					return err
				}
				domains, err := metaAccessor.GetDomainsForField(f)
				if err != nil {
					return err
				}
				for _, d := range domains {
					// This location can only be reached for new field inserts,
					// can assume that domain, has not yet exists
					err = st.AddDomain(&d)
					if err != nil {
						return err
					}
				}
			}
		}
		///////////////////////////////
		//   SCHEMA_FIELD_ASSOCIATION
		//      check for both cases - field already exists or new insert
		//      since the same field can be associated to multiple schemas
		sf, err := metaAccessor.GetSchemaFieldAssociation(s, f)
		if err != nil {
			return err
		}
		flagAssociation, err := st.SchemaFieldAssociationExists(sf)
		if err != nil {
			return err
		}
		if !flagAssociation {
			log.Printf("Unable to find association between schema=%s and field=%s. Adding to schema_field table...\n", s.Name, f.Name)
			err = st.AddSchemaFieldAssociation(sf)
			if err != nil {
				return err
			}
		}
	}

	////////////////////////////////////////
	//  DATASET comes after the fields loop
	//      Quality handling is implicit within the GetDataset call on metaAccessor
	dataset, err := metaAccessor.GetDataset(st, s)
	if err != nil {
		return err
	}
	err = st.GetDataset(&dataset)
	if err != nil {
		return err
	}

	sqlArg := store.GenerateSqlArg(shp2DbName)
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
		execStr := fmt.Sprintf(`ogr2ogr -f "PostgreSQL" PG:"%s" %s -lco precision=no -lco fid=fd_id -lco geometry_name=shape -nln %s.%s %s`,
			strings.ReplaceAll(cfg.StoreConfig.ConnStr, "database=", "dbname="),
			cfg.ShpPath, store.DbSchema, dataset.TableName, sqlArg,
		)
		fmt.Println(execStr)
		_, err = exec.Command(
			"sh", "-c", execStr,
		).Output()
	} else {
		// dataset already exists
		flagDataInStore, err := st.ShpDataInStore(dataset, metaAccessor.S)
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
		log.Printf("Data uploaded to dataset.name=%s dataset.id=%s", dataset.Name, dataset.Id)
	}
	return err
}

func ChangeAccess(cfg config.Config) error {
	// st, err := store.NewStore(cfg)
	// var accessId uuid.UUID
	// group := model.Group{
	// 	Id:   accessId,
	// 	Name: "",
	// }
	// member := model.Member{
	// 	DatasetId:  cfg.AccessConfig.DatasetId,
	// 	Group:      cfg.AccessConfig.Group,
	// 	Role:       cfg.AccessConfig.Role,
	// 	Permission: types.RolePermission[cfg.AccessConfig.Role],
	// }
	// accessId, err = st.GetAccessId(access)
	// if err != nil {
	// 	return err
	// }
	// if accessId == uuid.Nil {
	// 	_, err = st.AddAccess(access)
	// }
	// return err
	return nil
}

func AddElevation(cfg config.Config) error {
	return nil
}
