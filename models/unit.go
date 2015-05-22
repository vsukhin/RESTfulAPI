package models

import (
	"errors"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	UNIT_NAME_DEFAULT = "Название объединения по умолчанию"
)

// Структура для организации хранения объединений
type ViewShortUnit struct {
	Name string `json:"name" validate:"nonzero,min=1,max=255"` // Название объединения
}

type ViewLongUnit struct {
	Name    string `json:"name" validate:"nonzero,min=1,max=255"` // Название объединения
	Deleted bool   `json:"del"`                                   // Удален
}

type ApiShortMetaUnit struct {
	Total int64 `json:"count"` // Общее количество объединений
}

type ApiLongMetaUnit struct {
	NumOfUsers      int64 `json:"users"`         // Общее количество пользователей
	NumOfTables     int64 `json:"tables"`        // Общее количество таблиц
	NumOfProjects   int64 `json:"projects"`      // Общее количество проектов
	NumOfOrders     int64 `json:"orders"`        // Общее количество заказов
	NumOfFacilities int64 `json:"services"`      // Общее количество услуг
	NumOfCompanies  int64 `json:"organisations"` // Общее количество компаний
	NumOfSMSSenders int64 `json:"smsfroms"`      // Общее количество отправителей
}

type ApiShortUnit struct {
	ID   int64  `json:"id"`   // Уникальный идентификатор объединения
	Name string `json:"name"` // Название объединения
}

type ApiLongUnit struct {
	ID      int64     `json:"id"`      // Уникальный идентификатор объединения
	Created time.Time `json:"created"` // Время создания объединения
	Name    string    `json:"name"`    // Название объединения
	Deleted bool      `json:"del"`     // Удален
}

type DtoUnit struct {
	ID      int64     `db:"id"`      // Уникальный идентификатор объединения
	Created time.Time `db:"created"` // Время создания объединения
	Name    string    `db:"name"`    // Название объединения
	Active  bool      `db:"active"`  // Активен
}

// Конструктор создания объекта объединения в api
func NewApiShortMetaUnit(total int64) *ApiShortMetaUnit {
	return &ApiShortMetaUnit{
		Total: total,
	}
}

func NewApiLongMetaUnit(numofusers int64, numoftables int64, numofprojects int64, numoforders int64,
	numoffacilities int64, numofcompanies int64, numofsmssenders int64) *ApiLongMetaUnit {
	return &ApiLongMetaUnit{
		NumOfUsers:      numofusers,
		NumOfTables:     numoftables,
		NumOfProjects:   numofprojects,
		NumOfOrders:     numoforders,
		NumOfFacilities: numoffacilities,
		NumOfCompanies:  numofcompanies,
		NumOfSMSSenders: numofsmssenders,
	}
}

func NewApiShortUnit(id int64, name string) *ApiShortUnit {
	return &ApiShortUnit{
		ID:   id,
		Name: name,
	}
}

func NewApiLongUnit(id int64, created time.Time, name string, deleted bool) *ApiLongUnit {
	return &ApiLongUnit{
		ID:      id,
		Created: created,
		Name:    name,
		Deleted: deleted,
	}
}

type UnitSearch struct {
	ID   int64  `query:"id" search:"id"`     // Уникальный идентификатор объединения
	Name string `query:"name" search:"name"` // Название
}

// Конструктор создания объекта объединения в бд
func NewDtoUnit(id int64, created time.Time, name string, active bool) *DtoUnit {
	return &DtoUnit{
		ID:      id,
		Created: created,
		Name:    name,
		Active:  active,
	}
}

func (unit *UnitSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, unit), nil
}

func (unit *UnitSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, unit)
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
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (unit *UnitSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(unit)
}

func (unit *ViewShortUnit) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(unit, errors, req)
}

func (unit *ViewLongUnit) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(unit, errors, req)
}
