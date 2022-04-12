package shp

import (
	"errors"
	"fmt"

	"github.com/jonas-p/go-shp"
)

func NewShp(src string) (*shp.Reader, error) {
	shpf, err := shp.Open(src)
	defer shpf.Close()
	return shpf, err
}

// UniqueValues determines all unique values from field
func UniqueValues(shpF *shp.Reader, f shp.Field) ([]string, error) {
	var valSlice []string
	vals := make(map[string]bool)

	fIdx, err := FieldIdx(shpF, f.String())
	if err != nil {
		return []string{}, err
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
	return valSlice, nil
}

// FieldIdx loops to find field index in the shp file
func FieldIdx(shpF *shp.Reader, fs string) (int, error) {
	fields := shpF.Fields()
	for i, field := range fields {
		if field.String() == fs {
			return i, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("shp file does not contain field=%s", fs))
}
