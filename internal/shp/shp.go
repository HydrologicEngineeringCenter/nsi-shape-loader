package shp

import (
	"github.com/jonas-p/go-shp"
)

func NewShp(src string) (*shp.Reader, error) {
	shpf, err := shp.Open(src)
	defer shpf.Close()
	return shpf, err
}

// UniqueValues determines all unique values from field
func UniqueValues(shpF *shp.Reader, f shp.Field) []string {
	var valSlice []string
	var fIdx int
	vals := make(map[string]bool)
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
		// using hacky way to constraint values to a unique (ie non-repeated) set here
		// map grows on the heap, could be a performance sink
		if !vals[val] {
			vals[val] = true
		}
	}
	for key := range vals {
		valSlice = append(valSlice, key)
	}
	return valSlice
}
