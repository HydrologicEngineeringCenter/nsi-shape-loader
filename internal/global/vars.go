package global

const (
	APP_NAME                            = "seahorse"
	APP_VERSION                         = "0.5.0"
	DB_SCHEMA                           = "nsiv29test" // CAUTION! CHANGING TO A PRODUCTION SCHEMA CAN OVERRIDE DATA
	BASE_META_XLSX_PATH                 = "./assets/baseMetaData.xlsx"
	COPY_XLSX_PATH                      = "./metadata.xlsx"
	NATIONAL_MAP_SCHEME                 = "https"
	NATIONAL_MAP_HOST                   = "tnmaccess.nationalmap.gov"
	NATIONAL_MAP_PATH                   = "api/v1/products"
	ELEVATION_COLUMN_NAME               = "ground_elev"
	NO_PARALLEL_REQUEST_TO_NATIONAL_MAP = 100
)
