package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct; got %s", v.Kind())
	}

	// Iterate through struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		// Set the map key to the field name and value to the field value
		result[field.Name] = value.Interface()
	}

	return result, nil
}
