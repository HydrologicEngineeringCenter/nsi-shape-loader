package types

import "github.com/jonas-p/go-shp"

type DataFields map[string]string

type Data string

const (
	Integer Data = "Integer"
)

type Shape string

const (
	Point   Shape = "Point"
	Polygon Shape = "Polygon"
)

var (
	GeometryReverse = map[shp.ShapeType]string{
		shp.NULL:        "NULL",
		shp.POINT:       "POINT",
		shp.POLYLINE:    "POLYLINE",
		shp.POLYGON:     "POLYGON",
		shp.MULTIPOINT:  "MULTIPOINT",
		shp.POINTZ:      "POINTZ",
		shp.POLYLINEZ:   "POLYLINEZ",
		shp.POLYGONZ:    "POLYGONZ",
		shp.MULTIPOINTZ: "MULTIPOINTZ",
		shp.POINTM:      "POINTM",
		shp.POLYLINEM:   "POLYLINEM",
		shp.POLYGONM:    "POLYGONM",
		shp.MULTIPOINTM: "MULTIPOINTM",
		shp.MULTIPATCH:  "MULTIPATCH",
	}
)

type NsiField string

// Field type uses mapping from go-shape
const (
	String NsiField = "C"
	Number NsiField = "N"
	Float  NsiField = "F"
	Date   NsiField = "D"
)

var (
	NsiFieldReverse = map[string]NsiField{
		"C": String,
		"N": Number,
		"F": Float,
		"D": Date,
	}
)

type Quality string

const (
	High Quality = "High"
)

type Role string

const (
	Admin Role = "Admin"
)

type Permission string

const (
	Upload Permission = "Upload"
)

// FROM go-shp
// //  is a identifier for the the type of shapes.
// type  int32

// // These are the possible shape types.
// const (
// 	NULL         = 0
// 	POINT        = 1
// 	POLYLINE     = 3
// 	POLYGON      = 5
// 	MULTIPOINT   = 8
// 	POINTZ       = 11
// 	POLYLINEZ    = 13
// 	POLYGONZ     = 15
// 	MULTIPOINTZ  = 18
// 	POINTM       = 21
// 	POLYLINEM    = 23
// 	POLYGONM     = 25
// 	MULTIPOINTM  = 28
// 	MULTIPATCH   = 31
// )
