package shp

import (
	"github.com/HydrologicEngineeringCenter/nsi-shape-loader/internal/util"
	"github.com/jonas-p/go-shp"
)

func NewShp(src string) (*shp.Reader, error) {
	shpf, err := shp.Open(src)
	defer shpf.Close()
	return shpf, err
}

// UniqueValues determines all unique values from field
func UniqueValues(shpF *shp.Reader, f shp.Field) []string {
	var vals []string
	var fIdx int
	fields := shpF.Fields()

	// Loop to find field index
	for i, field := range fields {
		if field.String() == f.String() {
			fIdx = i
			break
		}
	}
	// Loop to find all unique values
	for i := 0; i < shpF.AttributeCount(); i++ {
		val := shpF.ReadAttribute(i, fIdx)
		if !util.StrContains(vals, val) {
			vals = append(vals, val)
		}
	}
	return vals
}
