package models

import (
	"reflect"
)

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
