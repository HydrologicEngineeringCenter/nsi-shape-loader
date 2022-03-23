package types

import "github.com/jonas-p/go-shp"

type Shape string

// const (
// 		shp.NULL Shape =        "NULL"
// 		shp.POINT =       "POINT"
// 		shp.POLYLINE =    "POLYLINE"
// 		shp.POLYGON =     "POLYGON"
// 		shp.MULTIPOINT =  "MULTIPOINT"
// 		shp.POINTZ =      "POINTZ"
// 		shp.POLYLINEZ =   "POLYLINEZ"
// 		shp.POLYGONZ =    "POLYGONZ"
// 		shp.MULTIPOINTZ = "MULTIPOINTZ"
// 		shp.POINTM =      "POINTM"
// 		shp.POLYLINEM =   "POLYLINEM"
// 		shp.POLYGONM =    "POLYGONM"
// 		shp.MULTIPOINTM = "MULTIPOINTM"
// 		shp.MULTIPATCH =  "MULTIPATCH"
// )

var (
	ShapeReverse = map[shp.ShapeType]Shape{
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

type Datatype string

// Field type uses mapping from go-shp
const (
	Char   Datatype = "C"
	Number          = "N"
	Float           = "F"
	Date            = "D"
)

var (
	DatatypeReverse = map[string]Datatype{
		"C": Char,
		"N": Number,
		"F": Float,
		"D": Date,
	}
)

type Quality string

const (
	High   Quality = "High"
	Medium         = "Medium"
	Low            = "Low"
)

var (
	QualityReverse = map[string]Quality{
		"High":   High,
		"Medium": Medium,
		"Low":    Low,
	}
)

type Role string

const (
	Admin Role = "admin"
	Owner      = "owner"
	User       = "user"
)

var (
	RolePermission = map[Role]string{
		Admin: "read add delete update",
		Owner: "read add update",
		User:  "read",
	}
)

// type Permission string

// const (
// 	Read   Permission = "Read"
// 	Add               = "Add"
// 	Edit              = "Edit"
// 	Delete            = "Delete"
// 	All               = "All"
// )

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

type Mode string

const (
	Prep   Mode = "prep"
	Upload      = "upload"
	Access      = "access"
)

var (
	ModeReverse = map[string]Mode{
		"prep":   Prep,
		"upload": Upload,
		"access": Access,
	}
)
