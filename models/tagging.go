package models

import (
	"reflect"
)

func GetFieldValue(field string, object interface{}) (value interface{}, found bool) {
	found = false
	value = nil
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		if field == structAddr.Type().Field(i).Tag.Get("json") {
			found = true
			value = structAddr.Field(i).Interface()
			break
		}
	}

	return value, found
}

func GetSearchTag(param string, object interface{}) (search string) {
	search = ""
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("query")
		if param == fieldTag {
			search = structAddr.Type().Field(i).Tag.Get("search")
			break
		}
	}

	return search
}

func GetJsonTag(field string, object interface{}) (name string) {
	name = ""
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("json")
		if field == structAddr.Type().Field(i).Name {
			name = fieldTag
			break
		}
	}

	return name
}

func CheckQueryTag(param string, object interface{}) (found bool) {
	found = false
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("query")
		if param == fieldTag {
			found = true
			break
		}
	}

	return found
}

func GetAllSearchTags(object interface{}) (tags *[]string) {
	tags = new([]string)
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("search")
		if fieldTag != "" && fieldTag != "-" {
			*tags = append(*tags, fieldTag)
		}
	}

	return tags
}

func GetAllGroupTags(object interface{}) (tags *[]string) {
	tags = new([]string)
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("group")
		if fieldTag != "" && fieldTag != "-" {
			*tags = append(*tags, fieldTag)
		}
	}

	return tags
}

func GetDbTagValue(param string, object interface{}) (value interface{}, found bool) {
	found = false
	value = nil
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("db")
		if param == fieldTag {
			found = true
			value = structAddr.Field(i).Interface()
			break
		}
	}

	return value, found
}

func SetDbTagValue(param string, object interface{}, value interface{}) (found bool) {
	found = false
	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get("db")
		if param == fieldTag {
			found = true
			switch structAddr.Field(i).Type().String() {
			case "bool":
				data, ok := value.(bool)
				if ok {
					structAddr.Field(i).SetBool(data)
				}
			case "string":
				data, ok := value.(string)
				if ok {
					structAddr.Field(i).SetString(data)
				}
			}
			break
		}
	}

	return found
}
