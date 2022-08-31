package util

import "reflect"

func GetFieldValue(v interface{}, fieldName string) (value string, ok bool) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	fieldValue := val.FieldByName(fieldName)

	if fieldValue.IsValid() {
		value = fieldValue.String()
		ok = true
	}

	return value, ok
}
