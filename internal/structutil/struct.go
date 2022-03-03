package structutil

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func SetField(item interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	fieldNames := map[string]int{}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		fieldNames[typeField.Name] = i
	}

	fieldNum, ok := fieldNames[fieldName]
	if !ok {
		return fmt.Errorf("field %s does not exist within the provided item", fieldName)
	}
	fieldVal := v.Field(fieldNum)
	valueKind := fieldVal.Kind()
	var err error
	switch valueKind {
	case reflect.Bool:
		coercedVal, err := strconv.ParseBool(value.(string))
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(coercedVal))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		coercedVal, err := strconv.ParseInt(value.(string), 0, 64)
		if err != nil {
			return err
		}
		// fieldVal.SetString(reflect.ValueOf(coercedVal))
		fieldVal.SetInt(coercedVal)
	case reflect.Float32, reflect.Float64:
		coercedVal, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		fieldVal.SetFloat(coercedVal)
	case reflect.String:
		fieldVal.Set(reflect.ValueOf(value))
	default:
		err = errors.New("value not coercible")
	}
	return err
}
