package models

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
