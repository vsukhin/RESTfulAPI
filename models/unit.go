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
	UNIT_NAME_DEFAULT = "Название компании"
)

// Структура для организации хранения объединений
type ViewShortUnit struct {
	Name string `json:"name" validate:"min=1,max=255"` // Название объединения
}

type ViewLongUnit struct {
	Name    string `json:"name" validate:"min=1,max=255"` // Название объединения
	Deleted bool   `json:"del"`                           // Удален
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
	NumOfInvoices   int64 `json:"invoices"`      // Общее количество счетов
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

type ApiFullUnit struct {
	ID         int64     `json:"id" db:"id"`                   // Уникальный идентификатор объединения
	Name       string    `json:"name" db:"name"`               // Название объединения
	Created    time.Time `json:"created" db:"created"`         // Время создания объединения
	Subscribed bool      `json:"subscription" db:"subscribed"` // Платный режим
	Paid       bool      `json:"paid" db:"paid"`               // Оплачен
	Begin_Paid time.Time `json:"paidBegin" db:"begin_paid"`    // Начало оплаченного периода
	End_Paid   time.Time `json:"paidEnd" db:"end_paid"`        // Окончание оплаченного периода
}

type UnitSearch struct {
	ID   int64  `query:"id" search:"id"`                  // Уникальный идентификатор объединения
	Name string `query:"name" search:"name" group:"name"` // Название
}

type DtoUnit struct {
	ID         int64     `db:"id"`         // Уникальный идентификатор объединения
	Created    time.Time `db:"created"`    // Время создания объединения
	Name       string    `db:"name"`       // Название объединения
	Active     bool      `db:"active"`     // Активен
	Subscribed bool      `db:subscribed`   // Платный режим
	Paid       bool      `db:paid`         // Оплачен
	Begin_Paid time.Time `db:"begin_paid"` // Начало оплаченного периода
	End_Paid   time.Time `db:"end_paid"`   // Окончание оплаченного периода
	UUID       string    `db:"uuid"`       // UUID объединения
}

// Конструктор создания объекта объединения в api
func NewApiShortMetaUnit(total int64) *ApiShortMetaUnit {
	return &ApiShortMetaUnit{
		Total: total,
	}
}

func NewApiLongMetaUnit(numofusers int64, numoftables int64, numofprojects int64, numoforders int64,
	numoffacilities int64, numofcompanies int64, numofsmssenders int64, numofinvoices int64) *ApiLongMetaUnit {
	return &ApiLongMetaUnit{
		NumOfUsers:      numofusers,
		NumOfTables:     numoftables,
		NumOfProjects:   numofprojects,
		NumOfOrders:     numoforders,
		NumOfFacilities: numoffacilities,
		NumOfCompanies:  numofcompanies,
		NumOfSMSSenders: numofsmssenders,
		NumOfInvoices:   numofinvoices,
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

func NewApiFullUnit(id int64, name string, created time.Time, subscribed bool, paid bool,
	begin_paid time.Time, end_paid time.Time) *ApiFullUnit {
	return &ApiFullUnit{
		ID:         id,
		Name:       name,
		Created:    created,
		Subscribed: subscribed,
		Paid:       paid,
		Begin_Paid: begin_paid,
		End_Paid:   end_paid,
	}
}

// Конструктор создания объекта объединения в бд
func NewDtoUnit(id int64, created time.Time, name string, active bool, subscribed bool, paid bool,
	begin_paid time.Time, end_paid time.Time, uuid string) *DtoUnit {
	return &DtoUnit{
		ID:         id,
		Created:    created,
		Name:       name,
		Active:     active,
		Subscribed: subscribed,
		Paid:       paid,
		Begin_Paid: begin_paid,
		End_Paid:   end_paid,
		UUID:       uuid,
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
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (unit *UnitSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(unit)
}

func (unit *ViewShortUnit) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(unit, errors, req)
}

func (unit *ViewLongUnit) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(unit, errors, req)
}
