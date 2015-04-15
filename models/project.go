package models

import (
	"errors"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Структура для организации хранения проекта
type ViewProject struct {
	Name string `json:"name" validate:"nonzero,min=1,max=255"` // Название
}

type ViewUpdateProject struct {
	Name    string `json:"name" validate:"nonzero,min=1,max=255"` // Название
	Archive bool   `json:"archive"`                               // Aрхивирован
}

type ApiMetaProject struct {
	Total        int64 `json:"count"`   // Общее число проектов
	NumOfArchive int64 `json:"archive"` // Число архивированных проектов
}

type ApiShortProject struct {
	ID   int64  `json:"id" db:"id"`     // Уникальный идентификатор проекта
	Name string `json:"name" db:"name"` // Название
}

type ApiLongProject struct {
	ID      int64     `json:"id" db:"id"`           // Уникальный идентификатор проекта
	Name    string    `json:"name" db:"name"`       // Название
	Archive bool      `json:"archive" db:"archive"` // Aрхивирован
	Created time.Time `json:"created" db:"created"` // Время создания
}

type ProjectSearch struct {
	ID   int64  `query:"id" search:"id"`     // Уникальный идентификатор проекта
	Name string `query:"name" search:"name"` // Название
}

type DtoProject struct {
	ID      int64     `db:"id"`      // Уникальный идентификатор проекта
	Unit_ID int64     `db:"unit_id"` // Идентификатор объединения
	Name    string    `db:"name"`    // Название
	Active  bool      `db:"active"`  // Aктивен
	Created time.Time `db:"created"` // Время создания
}

// Конструктор создания объекта проекта в api
func NewApiMetaProject(total int64, numofarchive int64) *ApiMetaProject {
	return &ApiMetaProject{
		Total:        total,
		NumOfArchive: numofarchive,
	}
}

func NewApiShortProject(id int64, name string) *ApiShortProject {
	return &ApiShortProject{
		ID:   id,
		Name: name,
	}
}

func NewApiLongProject(id int64, name string, archive bool, created time.Time) *ApiLongProject {
	return &ApiLongProject{
		ID:      id,
		Name:    name,
		Archive: archive,
		Created: created,
	}
}

// Конструктор создания объекта проекта в бд
func NewDtoProject(id int64, unit_id int64, name string, active bool, created time.Time) *DtoProject {
	return &DtoProject{
		ID:      id,
		Unit_ID: unit_id,
		Name:    name,
		Active:  active,
		Created: created,
	}
}

func (project *ProjectSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, project), nil
}

func (project *ProjectSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, project)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		if strings.Contains(invalue, "'") {
			errValue = errors.New("Wrong field value")
			break
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Uknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (project *ProjectSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(project)
}

func (project *ViewProject) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(project, errors, req)
}

func (project *ViewUpdateProject) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(project, errors, req)
}