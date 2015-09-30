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

// Структура для организации хранения документа
type ViewShortDocument struct {
	Company_ID int64  `json:"organisationId" validate:"nonzero"` // Идентификатор компании
	Begin_Date string `json:"begin" validate:"max=255"`          // Начало периода
	End_Date   string `json:"end" validate:"max=255"`            // Окончание периода
}

type ViewMiddleDocument struct {
	Company_ID int64  `json:"organisationId"  validate:"nonzero"` // Идентификатор компании
	Name       string `json:"name" validate:"max=255"`            // Название
	File_ID    int64  `json:"fileId" validate:"nonzero"`          // Идентификатор файла
}

type ViewLongDocument struct {
	Document_Type_ID int    `json:"categoryId" validate:"nonzero"` // Идентификатор типа документа
	Company_ID       int64  `json:"organisationId"`                // Идентификатор компании
	Name             string `json:"name" validate:"min=1,max=255"` // Название
	File_ID          int64  `json:"fileId" validate:"nonzero"`     // Идентификатор файла
}

type ApiMetaDocument struct {
	Total int64 `json:"count"` // Общее число документов
}

type ApiShortDocument struct {
	ID int64 `json:"id" db:"id"` // Уникальный идентификатор документа
}

type ApiLongDocument struct {
	ID               int64     `json:"id" db:"id"`                         // Уникальный идентификатор документа
	Document_Type_ID int       `json:"categoryId" db:"categoryId"`         // Идентификатор типа документа
	Unit_ID          int64     `json:"unitId" db:"unitId"`                 // Идентификатор объединения
	Company_ID       int64     `json:"organisationId" db:"organisationId"` // Идентификатор компании
	Name             string    `json:"name" db:"name"`                     // Название
	Created          time.Time `json:"created" db:"created"`               // Время создания
	Updated          time.Time `json:"edited" db:"edited"`                 // Время изменения
	Locked           bool      `json:"lock" db:"lock"`                     // Неизменяемый
	Pending          bool      `json:"pending" db:"pending"`               // Ожидающий
	File_ID          int64     `json:"fileId" db:"fileId"`                 // Идентификатор файла
}

type DocumentSearch struct {
	ID               int64     `query:"id" search:"id"`                                               // Уникальный идентификатор документа
	Document_Type_ID int       `query:"categoryId" search:"document_type_id"`                         // Идентификатор типа документа
	Unit_ID          int64     `query:"unitId" search:"unit_id"`                                      // Идентификатор объединения
	Company_ID       int64     `query:"organisationId" search:"company_id"`                           // Идентификатор компании
	Name             string    `query:"name" search:"name" group:"name"`                              // Название
	Created          time.Time `query:"created" search:"created" group:"convert(created using utf8)"` // Время создания
	Updated          time.Time `query:"edited" search:"updated" group:"convert(updated using utf8)"`  // Время изменения
	Locked           bool      `query:"lock" search:"locked"`                                         // Неизменяемый
	Pending          bool      `query:"pending" search:"pending"`                                     // Ожидающий
	File_ID          int64     `query:"fileId" search:"file_id"`                                      // Идентификатор файла
}

type DtoDocument struct {
	ID               int64     `db:"id"`               // Уникальный идентификатор документа
	Document_Type_ID int       `db:"document_type_id"` // Идентификатор типа компании
	Unit_ID          int64     `db:"unit_id"`          // Идентификатор объединения
	Company_ID       int64     `db:"company_id"`       // Идентификатор компании
	Name             string    `db:"name"`             // Название
	Locked           bool      `db:"locked"`           // Неизменяемый
	Pending          bool      `db:"pending"`          // Ожидающий
	File_ID          int64     `db:"file_id"`          // Идентификатор файла
	Created          time.Time `db:"created"`          // Время создания
	Updated          time.Time `db:"updated"`          // Время изменения
	Active           bool      `db:"active"`           // Aктивен
}

// Конструктор создания объекта документа в api
func NewApiMetaDocument(total int64) *ApiMetaDocument {
	return &ApiMetaDocument{
		Total: total,
	}
}

func NewApiShortDocument(id int64) *ApiShortDocument {
	return &ApiShortDocument{
		ID: id,
	}
}

func NewApiLongDocument(id int64, document_type_id int, unit_id int64, company_id int64, name string,
	created time.Time, updated time.Time, locked bool, pending bool, file_id int64) *ApiLongDocument {
	return &ApiLongDocument{
		ID:               id,
		Document_Type_ID: document_type_id,
		Unit_ID:          unit_id,
		Company_ID:       company_id,
		Name:             name,
		Created:          created,
		Updated:          updated,
		Locked:           locked,
		Pending:          pending,
		File_ID:          file_id,
	}
}

// Конструктор создания объекта документа в бд
func NewDtoDocument(id int64, document_type_id int, unit_id int64, company_id int64, name string,
	locked bool, pending bool, file_id int64, created time.Time, updated time.Time, active bool) *DtoDocument {
	return &DtoDocument{
		ID:               id,
		Document_Type_ID: document_type_id,
		Unit_ID:          unit_id,
		Company_ID:       company_id,
		Name:             name,
		Locked:           locked,
		Pending:          pending,
		File_ID:          file_id,
		Created:          created,
		Updated:          updated,
		Active:           active,
	}
}

func (document *DocumentSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, document), nil
}

func (document *DocumentSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, document)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "unitId":
		fallthrough
	case "organisationId":
		fallthrough
	case "fileId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "categoryId":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		fallthrough
	case "created":
		fallthrough
	case "edited":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "lock":
		fallthrough
	case "pending":
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

func (document *DocumentSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(document)
}

func (document *ViewShortDocument) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(document, errors, req)
}

func (document *ViewMiddleDocument) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(document, errors, req)
}

func (document *ViewLongDocument) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(document, errors, req)
}
