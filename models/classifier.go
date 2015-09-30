package models

import (
	"errors"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	CLASSIFIER_TYPE_MAIN = 0
)

// Структура для организации хранения классификатора
type ViewClassifier struct {
	Name string `json:"name" validate:"min=1,max=255"` // Название
}

type ViewUpdateClassifier struct {
	Name    string `json:"name" validate:"min=1,max=255"` // Название
	Deleted bool   `json:"del"`                           // Удален
}

type ApiShortClassifier struct {
	ID   int    `json:"id" db:"id"`     // Уникальный идентификатор классификатора
	Name string `json:"name" db:"name"` // Название
}

type ApiLongClassifier struct {
	ID      int    `json:"id" db:"id"`     // Уникальный идентификатор классификатора
	Name    string `json:"name" db:"name"` // Название
	Deleted bool   `json:"del" db:"del"`   // Удален
}

type ClassifierSearch struct {
	ID      int    `query:"id" search:"id"`                  // Уникальный идентификатор классификатора
	Name    string `query:"name" search:"name" group:"name"` // Название
	Deleted bool   `query:"del" search:"(not active)"`       // Удален
}

type DtoClassifier struct {
	ID      int       `db:"id"`      // Уникальный идентификатор классификатора
	Name    string    `db:"name"`    // Название
	Active  bool      `db:"active"`  // Aктивен
	Created time.Time `db:"created"` // Время создания
}

// Конструктор создания объекта классификатора в api
func NewApiShortClassifier(id int, name string) *ApiShortClassifier {
	return &ApiShortClassifier{
		ID:   id,
		Name: name,
	}
}

func NewApiLongClassifier(id int, name string, deleted bool) *ApiLongClassifier {
	return &ApiLongClassifier{
		ID:      id,
		Name:    name,
		Deleted: deleted,
	}
}

// Конструктор создания объекта классификатора в бд
func NewDtoClassifier(id int, name string, active bool, created time.Time) *DtoClassifier {
	return &DtoClassifier{
		ID:      id,
		Name:    name,
		Active:  active,
		Created: created,
	}
}

func (classifier *ClassifierSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, classifier), nil
}

func (classifier *ClassifierSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, classifier)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "del":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (classifier *ClassifierSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(classifier)
}

func (classifier *ViewClassifier) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(classifier, errors, req)
}

func (classifier *ViewUpdateClassifier) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(classifier, errors, req)
}
