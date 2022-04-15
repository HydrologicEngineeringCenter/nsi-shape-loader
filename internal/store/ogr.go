package store

import "fmt"

func GenerateSqlArg(shp2DbColMap map[string]string) string {
	sqlArg := `-sql "SELECT `
	i := 0
	noElements := len(shp2DbColMap)
	for k, v := range shp2DbColMap {
		sqlArg += fmt.Sprintf(` %s AS %s`, k, v)
		if i < noElements {
			sqlArg += ","
		}
		i++
	}
	sqlArg += `"`
	return sqlArg
}
