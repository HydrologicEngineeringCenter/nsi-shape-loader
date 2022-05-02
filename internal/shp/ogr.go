package shp

import "fmt"

// GenerateSqlArg generates the -sql argument required for ogr2ogr
func GenerateSqlArg(shp2DbColMap map[string]string, shpFileName string) string {
	sqlArg := `-sql "SELECT `
	i := 0
	noElements := len(shp2DbColMap)
	for k, v := range shp2DbColMap {
		sqlArg += fmt.Sprintf(` %s AS %s`, k, v)
		if i < noElements-1 {
			sqlArg += ","
		}
		i++
	}
	sqlArg += (` FROM \"` + shpFileName + `\"`)
	sqlArg += `"`
	return sqlArg
}
