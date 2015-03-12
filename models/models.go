/*
Mapping models on database.
*/
package models

import (
	logging "github.com/op/go-logging"
)

type IDs interface {
	GetIDs() []int64
}

type Checker interface {
	Check(field string) (valid bool, err error)
}

type Extractor interface {
	Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error)
	GetAllFields(parameter interface{}) (fields *[]string)
}

type OrderExp struct {
	Field string
	Order string
}

type FilterExp struct {
	Fields []string
	Op     string
	Value  string
}

type UserRole int

const (
	USER_ROLE_DEVELOPER UserRole = iota + 1
	USER_ROLE_ADMINISTRATOR
	USER_ROLE_SUPPLIER
	USER_ROLE_CUSTOMER
)

var (
	log = logging.MustGetLogger("models")
)
