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
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/model"
	shape "github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/shp"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/store"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/types"
	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/xls"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"
)

func Core(c *cli.Context) error {
	log.Printf("\n\n")
	log.Printf("============================================================")
	log.Printf("                     SEAHORSE %s                            ", config.APP_VERSION)
	log.Printf("============================================================")
	//  pre - generate config xls from shp
	//  upload - upload based on data and metadata from xls and shp
	//  access - change access group and role
	cfg, err := config.NewConfig(c)
	if err != nil {
		log.Fatal(err)
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
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// PreUpload generates an xls template from shp file fields
func Prep(cfg config.Config) error {

	// copy xls file
	const baseXlsSrc = config.BASE_META_XLSX_PATH
	const cpXlsDest = config.COPY_XLSX_PATH
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
		err = xlsF.F.SetCellRichText("field-domain", loc, []excelize.RichTextRun{
			{
				Text: val,
				Font: &excelize.Font{
					Bold:  false,
					Color: "FF0000",
				},
			},
		})
		if err != nil {
			return err
		}
	}
	xlsF.F.Save()
	wd, err := os.Getwd()
	log.Println("Metadata template file successfully created at:", filepath.Join(wd, cpXlsDest))
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
	// log.Printf("Retrieving id for unique schema=%s version=%s\n", s.Name, s.Version)
	err = st.GetSchemaId(&s)
	if err != nil {
		return err
	}
	if s.Id == uuid.Nil {
		// log.Printf("schema=%s version=%s do not exists. Adding to schema table...\n", s.Name, s.Version)
		err = st.AddSchema(&s)
		if err != nil {
			return err
		}
	}

	fields, err := metaAccessor.GetFields()
	if err != nil {
		return err
	}
	for _, f := range fields {
		// log.Printf("Retrieving id for unique field=%s type=%s\n", f.Name, f.Type)
		err = st.GetFieldId(&f)
		if err != nil {
			return err
		}
		// If no id -> field is not in db -> add field + add association to schema + domain
		if f.Id == uuid.Nil {
			// log.Printf("field=%s type=%s do not exists. Adding to field table...\n", f.Name, f.Type)
			err = st.AddField(&f)
			if err != nil {
				return err
			}
			///////////////////////////////
			//   DOMAIN
			// Process domain only if specified by field ie. field holds a discrete categorical variable
			// Currently this is specified from the metadata xls, could be a TODO to automatically detect field based only on the shp file
			if f.IsDomain {
				// log.Printf("field=%s holds discrete categorical variables. Adding to domain table...\n", f.Name)
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
			// log.Printf("Unable to find association between schema=%s and field=%s. Adding to schema_field table...\n", s.Name, f.DbName)
			err = st.AddSchemaFieldAssociation(sf)
			if err != nil {
				return err
			}
		}
	}

	// quality
	q, err := metaAccessor.GetQuality(st)
	if err != nil {
		return err
	}
	err = st.GetQuality(&q)
	if err != nil {
		return err
	}

	// group
	g, err := metaAccessor.GetGroup()
	if err != nil {
		return err
	}
	err = st.GetGroupId(&g)
	if err != nil {
		return err
	}
	if g.Id == uuid.Nil {
		err = st.AddGroup(&g)
		if err != nil {
			return err
		}
	}

	// dataset
	// quality handling is implicit within the GetDataset call on metaAccessor
	d, err := metaAccessor.GetDataset(st, s, g)
	if err != nil {
		return err
	}
	err = st.GetDataset(&d)
	if err != nil {
		return err
	}
	// map field name from shp to what will be in postgis
	shp2DbName, err := metaAccessor.GetShpDbFieldNameMap()
	if err != nil {
		return err
	}
	sqlArg := store.GenerateSqlArg(shp2DbName, strings.TrimSuffix(
		filepath.Base(cfg.ShpPath),
		filepath.Ext(cfg.ShpPath),
	))
	var execStr string
	if d.Id == uuid.Nil {
		// creating new dataset
		d.TableName = "inventory_" + strings.ReplaceAll(uuid.New().String(), "-", "_")
		err = st.AddDataset(&d)
		if err != nil {
			return err
		}
		// create new table
		log.Printf("Creating table=%s for dataset=%s", d.TableName, d.Name)
		execStr = fmt.Sprintf(`ogr2ogr -f "PostgreSQL" PG:"%s" %s -lco precision=no -lco fid=fd_id -lco geometry_name=shape -nln %s.%s %s`,
			strings.ReplaceAll(cfg.StoreConfig.ConnStr, "database=", "dbname="),
			cfg.ShpPath, store.DbSchema, d.TableName, sqlArg,
		)
	} else {
		// dataset already exists
		flagDataInStore, err := st.ShpDataInStore(d, metaAccessor.S)
		if err != nil {
			return err
		}
		if !flagDataInStore { // data has not yet been added to store
			log.Printf("table=%s exists for dataset=%s. Appending rows...", d.TableName, d.Name)
			execStr = fmt.Sprintf(`ogr2ogr -append -update -f "PostgreSQL" PG:"%s" %s -lco precision=no -nln %s.%s`,
				strings.ReplaceAll(cfg.StoreConfig.ConnStr, "database=", "dbname="),
				cfg.ShpPath, store.DbSchema, d.TableName,
			)
		} else {
			return errors.New("Upload failed - shp file has already been uploaded")
		}
	}
	// fmt.Println(execStr)
	_, err = exec.Command(
		"sh", "-c", execStr,
	).Output()
	if err != nil {
	} else {
		err = st.UpdateDatasetBBox(d)
		if err != nil {
			return err
		}
		log.Printf("Data uploaded to dataset.name=%s dataset.table_name=%s", d.Name, d.TableName)
	}
	return err
}

func ChangeAccess(cfg config.Config) error {
	st, err := store.NewStore(cfg)
	if err != nil {
		return err
	}

	// group
	g := model.Group{
		Name: cfg.AccessConfig.Group,
	}
	err = st.GetGroupId(&g)
	if err != nil {
		return err
	}
	if g.Id == uuid.Nil {
		return errors.New(fmt.Sprintf("Changing access role failed - group.name=%s does not exists in the database", cfg.AccessConfig.Group))
	}

	// member
	m := model.Member{
		GroupId: g.Id,
		Role:    cfg.AccessConfig.Role,
		UserId:  cfg.AccessConfig.UserId,
	}
	err = st.GetMemberId(&m)
	if err != nil {
		return err
	}
	if m.Id == uuid.Nil {
		// user has no association to the group
		err = st.AddMember(&m)
	} else {
		// user association exists
		err = st.UpdateMemberRole(&m)
	}
	log.Printf("member.user_id=%s now exists as member.role=%s for group.name=%s", m.UserId, m.Role, g.Name)
	return err
}

func AddElevation(cfg config.Config) error {
	return nil
}
