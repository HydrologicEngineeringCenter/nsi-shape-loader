package global

// APP
const (
	APP_NAME    = "seahorse"
	APP_VERSION = "0.9.1"
	DB_SCHEMA   = "nsiv291test" // CAUTION! CHANGING TO A PRODUCTION SCHEMA CAN OVERRIDE DATA
)

// PREP
const (
	BASE_META_XLSX_PATH = "./assets/baseTemplate.xlsx"
	COPY_XLSX_PATH      = "./assets/metadata.xlsx"
)

// ELEVATION
const (
	ELEVATION_COLUMN_NAME          = "ground_elev" // ground_elev is hardwired into struct tags, there are multiple source of truth for this value
	ELEVATION_BATCHSIZE            = 10000
	ELEVATION_NO_PARALLEL_ROUTINES = 4
	NATIONAL_MAP_DATASET           = "National Elevation Dataset (NED) 1/3 arc-second"
	NATIONAL_MAP_SCHEME            = "https"
	NATIONAL_MAP_HOST              = "tnmaccess.nationalmap.gov"
	NATIONAL_MAP_PATH              = "api/v1/products"
	NATIONAL_MAP_CACHE_BASEPATH    = "./assets/dem/"
)
