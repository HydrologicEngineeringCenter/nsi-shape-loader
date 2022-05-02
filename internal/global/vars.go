package global

// APP
const (
	APP_NAME              = "seahorse"
	APP_VERSION           = "0.5.0"
	DB_SCHEMA             = "nsiv29test" // CAUTION! CHANGING TO A PRODUCTION SCHEMA CAN OVERRIDE DATA
	BASE_META_XLSX_PATH   = "./assets/baseMetaData.xlsx"
	COPY_XLSX_PATH        = "./metadata.xlsx"
	ELEVATION_COLUMN_NAME = "ground_elev"
)

// NATIONAL MAP REQUEST
const (
	NATIONAL_MAP_DATASET           = "National Elevation Dataset (NED) 1/9 arc-second"
	NATIONAL_MAP_SCHEME            = "https"
	NATIONAL_MAP_HOST              = "tnmaccess.nationalmap.gov"
	NATIONAL_MAP_PATH              = "api/v1/products"
	NATIONAL_MAP_CACHE_BASEPATH    = "test/dem/"
	NATIONAL_MAP_QUERY_RESULT_JSON = "assets/TNMQuery13ArcSec.json"
)
