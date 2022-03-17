package core

import (
	"fmt"

	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/config"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/files"
	shape "github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/shp"
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/xls"
	"github.com/urfave/cli/v2"
)

func Core(c *cli.Context) error {
	// 2 modes to core functionalities
	//  pre - generate config xls from shp
	//  upload - upload based on data and metadata from xls and shp
	cfg, err := config.NewConfig(c)
	if cfg.Mode == config.Pre {
		err = PreUpload(cfg)
	} else {
		err = PreUpload(cfg)
	}
	return err
}

// PreUpload generates an xls template from shp file fields
func PreUpload(cfg config.Config) error {

	// copy xls file
	const baseXlsSrc = "./assets/baseMetadata.xlsx"
	const cpXlsDest = "./metadata.xlsx"
	err := files.Copy(baseXlsSrc, cpXlsDest)
	if err != nil {
		return err
	}

	// populate fields worksheet based on shp fields
	xls, err := xls.NewXls(cpXlsDest)
	if err != nil {
		return err
	}
	shpf, err := shape.NewShp(cfg.ShpPath)
	if err != nil {
		return err
	}
	fields := shpf.Fields()
	var loc string
	for j, f := range fields {
		loc = "B" + fmt.Sprint(j+2)
		// val = f.String()
		err = xls.SetCellValue("fields", loc, f.Name)
		if err != nil {
			return err
		}
	}
	return err
}

//// Upload populates metadata from the config xls and upload data from shp file
//func Upload(cfg config.Config) error {
//	st, err := store.NewStore(cfg)
//	if err != nil {
//		return err
//	}

//	// fields := shpf.Fields()

//	// var insertTable = dq.TableDataSet{
//	// 	Name:       cfg.Dbtablename,
//	// 	Statements: map[string]string{},
//	// 	Fields:     model.Point{},
//	// }

//	// batchSize := 20000
//	// lastRecordIdx := shpf.AttributeCount() - 1
//	// var records []model.Point
//	// for shpf.Next() {
//	// 	i, _ := shpf.Shape()

//	// 	// construct data struct from point
//	// 	var newPoint model.Point
//	// 	for j, f := range fields {
//	// 		val := shpf.ReadAttribute(i, j)
//	// 		fieldStr := strings.Title(strings.ToLower(f.String()))
//	// 		structutil.SetField(&newPoint, fieldStr, val)
//	// 	}
//	// 	records = append(records, newPoint)
//	// 	// Batch upload on reaching batchsize limit
//	// 	if (i != 0 && i%batchSize == 0) || i == lastRecordIdx {
//	// 		if i == lastRecordIdx {
//	// 			// batching the last odd lot records doesn't work for some reason
//	// 			err = st.DS.Insert(&insertTable).
//	// 				Records(&records).
//	// 				Execute()
//	// 		} else {
//	// 			err = st.DS.Insert(&insertTable).
//	// 				Records(&records).
//	// 				Batch(true).
//	// 				BatchSize(len(records)).
//	// 				Execute()
//	// 		}
//	// 		if err != nil {
//	// 			return err
//	// 		}
//	// 		fmt.Println("Proccessed " + fmt.Sprint(i+1) + " records")
//	// 		records = []model.Point{}
//	// 	}
//	// }
//	// fmt.Println("Processing finished.")

//	/////////////////////////////////////////////////////////
//	// Fill out field + domain from included XLS
//	file, err := excelize.OpenFile(cfg.FieldMap)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		// Close the spreadsheet.
//		err = file.Close()
//	}()

//	//  Check if schema already exists
//	//  If it is, then use id to populate dataset meta
//	//  Otherwise, add a new schema
//	var schemaId uuid.UUID
//	schemaName, err := file.GetCellValue("schema", "A2")
//	schemaVersion, err := file.GetCellValue("schema", "B2")
//	schemaNotes, err := file.GetCellValue("schema", "C2")

//	schema := model.Schema{
//		Name:    schemaName,
//		Version: schemaVersion,
//		Notes:   schemaNotes,
//	}
//	schemaId, err = st.GetSchemaId(schema)
//	if err != nil {
//		return err
//	}
//	if schemaId == uuid.Nil { // schema do not exists
//		schemaId, err = st.AddSchema(schema)

//		// // Start adding in fields
//		// rows, err := file.GetRows()
//		// if err != nil {
//		// 	fmt.Println(err)
//		// 	return "", err
//		// }
//		// for _, row := range rows {
//		// 	for _, colCell := range row {
//		// 		fmt.Print(colCell, "\t")
//		// 	}
//		// 	fmt.Println()
//		// }

//		// Populate meta data
//		fmt.Println("Reading shapefile from: " + cfg.FilePath)
//		shpf, err := shp.Open(cfg.FilePath)
//		if err != nil {
//			return err
//		}
//		defer shpf.Close()
//		shapeType := types.GeometryReverse[shpf.GeometryType]

//		fields := shpf.Fields()
//		for j, f := range fields {
//			uuidNsiField := uuid.New()
//			uuidDomain := uuid.New()

//			// Field quantitative vs qualitative
//			field := model.Field{
//				Id:          uuidNsiField,
//				Name:        f.String(),
//				Type:        types.NsiFieldReverse[string(f.Fieldtype)],
//				Description: "",
//				IsDomain:    false,
//			}

//			domain := model.Domain{
//				Id:         uuidDomain,
//				NsiFieldId: uuidNsiField,
//				Value:      0,
//				Datatype:   "",
//			}
//		}
//	}

//	// 	schemaField := model.SchemaField{
//	// 		Id:         [16]byte{},
//	// 		NsiFieldId: uuidNsiField,
//	// 	}

//	// 	schema := model.Schema{
//	// 		Id:      [16]byte{},
//	// 		Name:    schemaName,
//	// 		Version: schemaVersion,
//	// 		Notes:   schemaNotes,
//	// 	}
//	// 	err = st.AddSchema(schema)

//	// 	dataset := model.Dataset{
//	// 		Id:          [16]byte{},
//	// 		Name:        "",
//	// 		Version:     "",
//	// 		NsiSchemaId: [16]byte{},
//	// 		TableName:   "",
//	// 		Shape:       types.Shape(shapeType),
//	// 		Description: "",
//	// 		Purpose:     "",
//	// 		DateCreated: time.Time{},
//	// 		CreatedBy:   "",
//	// 		QualityId:   [16]byte{},
//	// 	}

//	// }

//	// // Upload data
//	// cmd, err := exec.Command("/bin/bash", "./upload").Output()

//	// return string(cmd), err
//	return err
//}
