package models

import (
	"errors"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Структура для организации хранения пользовательских таблиц
type ViewShortCustomerTable struct {
	Name   string `json:"name" validate:"nonzero,min=1,max=255"` // Название пользовательской таблицы
	UnitID int64  `json:"unitId"`                                // Объединение
}

type ViewLongCustomerTable struct {
	Name   string `json:"name" validate:"nonzero,min=1,max=255"` // Название пользовательской таблицы
	Type   string `json:"type" validate:"nonzero,min=1,max=255"` // Название пользовательской таблицы
	UnitID int64  `json:"unitId" validate:"nonzero"`             // Объединение
}

type ApiShortCustomerTable struct {
	Name   string `json:"name" db:"name"`      // Название пользовательской таблицы
	Type   string `json:"type" db:"type"`      // Тип
	UnitID int64  `json:"unitId" db:"unit_id"` // Идентификатор объединения
}

type ApiLongCustomerTable struct {
	ID     int64  `json:"id" db:"id"`          // Уникальный идентификатор пользовательской таблицы
	Name   string `json:"name" db:"name"`      // Название пользовательской таблицы
	Type   string `json:"type" db:"type"`      // Тип
	UnitID int64  `json:"unitId" db:"unit_id"` // Идентификатор объединения
}

type ApiMetaCustomerTable struct {
	NumOfRows      int64 `json:"rows"`          // Число строк
	NumOfCols      int64 `json:"columns"`       // Число колонок
	Checked        bool  `json:"verified"`      // Выполнялась проверка
	QaulityPer     byte  `json:"quality"`       // Качество данных
	NumOfWrongRows int64 `json:"incorrectRows"` // Количество строк с неверными данными
}

type TableSearch struct {
	ID     int64  `query:"id" search:"c.id"`          // Уникальный идентификатор пользовательской таблицы
	Name   string `query:"name" search:"c.name"`      // Название пользовательской таблицыя
	Type   string `query:"type" search:"t.name"`      // Тип
	UnitID int64  `query:"unitId" search:"c.unit_id"` // Идентификатор объединения
}

type DtoCustomerTable struct {
	ID        int64     `db:"id"`        // Уникальный идентификатор пользовательской таблицы
	Name      string    `db:"name"`      // Название
	Created   time.Time `db:"created"`   // Время создания
	TypeID    int64     `db:"type_id"`   // Идентификатор типа
	UnitID    int64     `db:"unit_id"`   // Идентификатор объединения
	Active    bool      `db:"active"`    // Активная
	Permanent bool      `db:"permanent"` // Постоянная
}

// Конструктор создания объекта пользовательской таблицы в api
func NewViewShortCustomerTable(name string, unitid int64) *ViewShortCustomerTable {
	return &ViewShortCustomerTable{
		Name:   name,
		UnitID: unitid,
	}
}

func NewViewLongCustomerTable(name string, typevalue string, unitid int64) *ViewLongCustomerTable {
	return &ViewLongCustomerTable{
		Name:   name,
		Type:   typevalue,
		UnitID: unitid,
	}
}

func NewApiShortCustomerTable(name string, typevalue string, unitid int64) *ApiShortCustomerTable {
	return &ApiShortCustomerTable{
		Name:   name,
		Type:   typevalue,
		UnitID: unitid,
	}
}

func NewApiLongCustomerTable(id int64, name string, typevalue string, unitid int64) *ApiLongCustomerTable {
	return &ApiLongCustomerTable{
		ID:     id,
		Name:   name,
		Type:   typevalue,
		UnitID: unitid,
	}
}

func NewApiMetaCustomerTable(numofrows int64, numofcols int64, checked bool, qualityper byte, numofwrongrows int64) *ApiMetaCustomerTable {
	return &ApiMetaCustomerTable{
		NumOfRows:      numofrows,
		NumOfCols:      numofcols,
		Checked:        checked,
		QaulityPer:     qualityper,
		NumOfWrongRows: numofwrongrows,
	}
}

// Конструктор создания объекта пользовательской таблицы в бд
func NewDtoCustomerTable(id int64, name string, created time.Time, typeid int64, unitid int64, active bool, permanent bool) *DtoCustomerTable {
	return &DtoCustomerTable{
		ID:        id,
		Name:      name,
		Created:   created,
		TypeID:    typeid,
		UnitID:    unitid,
		Active:    active,
		Permanent: permanent,
	}
}

func (customertable *ViewShortCustomerTable) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(customertable, errors, req)
}

func (customertable *ViewLongCustomerTable) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(customertable, errors, req)
}

func (table *TableSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, table), nil
}

func (table *TableSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, table)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "unitId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		fallthrough
	case "type":
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

func (table *TableSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(table)
}
