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
	CUSTOMER_TABLE_SIGNATURE_DEFAULT = "Система"
)

// Структура для организации хранения пользовательских таблиц
type ViewShortCustomerTable struct {
	Name   string `json:"name" validate:"min=1,max=255"` // Название пользовательской таблицы
	UnitID int64  `json:"unitId"`                        // Объединение
}

type ViewLongCustomerTable struct {
	Name   string `json:"name" validate:"min=1,max=255"` // Название пользовательской таблицы
	Type   int    `json:"type"`                          // Тип
	UnitID int64  `json:"unitId" validate:"nonzero"`     // Объединение
}

type ApiMetaCustomerTable struct {
	Total int64 `json:"count"` // Общее число пользовательских таблиц
}

type ApiShortCustomerTable struct {
	Name   string `json:"name" db:"name"`      // Название пользовательской таблицы
	Type   int    `json:"type" db:"type"`      // Тип
	UnitID int64  `json:"unitId" db:"unit_id"` // Идентификатор объединения
}

type ApiMiddleCustomerTable struct {
	ID   int64  `json:"id" db:"id"`     // Уникальный идентификатор пользовательской таблицы
	Name string `json:"name" db:"name"` // Название пользовательской таблицы
	Type int    `json:"type" db:"type"` // Тип
}

type ApiLongCustomerTable struct {
	ID     int64  `json:"id" db:"id"`          // Уникальный идентификатор пользовательской таблицы
	Name   string `json:"name" db:"name"`      // Название пользовательской таблицы
	Type   int    `json:"type" db:"type"`      // Тип
	UnitID int64  `json:"unitId" db:"unit_id"` // Идентификатор объединения
}

type ApiSearchCustomerTable struct {
	ID        int64     `json:"id" db:"id"`                   // Уникальный идентификатор пользовательской таблицы
	Name      string    `json:"name" db:"name"`               // Название пользовательской таблицы
	NumOfRows int64     `json:"rows" db:"rows"`               // Число строк
	Created   time.Time `json:"created" db:"created"`         // Время создания
	Signature string    `json:"createdName" db:"createdName"` // Автор
	Type      int       `json:"type" db:"type"`               // Тип
	UnitID    int64     `json:"unitId" db:"unit_id"`          // Идентификатор объединения
}

type ApiFullMetaCustomerTable struct {
	NumOfRows      int64     `json:"rows"`          // Число строк
	Created        time.Time `json:"created"`       // Время создания
	Signature      string    `json:"createdName"`   // Автор
	NumOfCols      int64     `json:"columns"`       // Число колонок
	Checked        bool      `json:"verified"`      // Выполнялась проверка
	QaulityPer     byte      `json:"quality"`       // Качество данных
	NumOfWrongRows int64     `json:"incorrectRows"` // Количество строк с неверными данными
}

type TableSearch struct {
	ID        int64     `query:"id" search:"c.id"`                                                                              // Уникальный идентификатор пользовательской таблицы
	Name      string    `query:"name" search:"c.name" group:"c.name"`                                                           // Название пользовательской таблицы
	NumOfRows int64     `query:"rows" search:"(select count(*) from table_data where customer_table_id = c.id and active = 1)"` // Число строк
	Created   time.Time `query:"created" search:"c.created" group:"convert(c.created using utf8)"`                              // Время создания
	Signature string    `query:"createdName" search:"c.signature" group:"c.signature"`                                          // Автор
	Type      int       `query:"type" search:"c.type_id"`                                                                       // Тип
	UnitID    int64     `query:"unitId" search:"c.unit_id"`                                                                     // Идентификатор объединения
}

type DtoCustomerTable struct {
	ID                      int64     `db:"id"`                      // Уникальный идентификатор пользовательской таблицы
	Name                    string    `db:"name"`                    // Название
	Created                 time.Time `db:"created"`                 // Время создания
	TypeID                  int       `db:"type_id"`                 // Идентификатор типа
	UnitID                  int64     `db:"unit_id"`                 // Идентификатор объединения
	Active                  bool      `db:"active"`                  // Активная
	Permanent               bool      `db:"permanent"`               // Постоянная
	Import_Ready            bool      `db:"import_ready"`            // Готовность импорта
	Import_Percentage       byte      `db:"import_percentage"`       // Процент импорта
	Import_Columns          int64     `db:"import_columns"`          // Количество импортированных колонок
	Import_Rows             int64     `db:"import_rows"`             // Количество импортированных строк
	Import_WrongRows        int64     `db:"import_wrongrows"`        // Количество импортированных ошибочных строк
	Signature               string    `db:"signature"`               // Автор
	Import_Error            bool      `db:"import_error"`            // Ошибка импорта
	Import_ErrorDescription string    `db:"import_errordescription"` // Описание ошибки импорта
}

// Конструктор создания объекта пользовательской таблицы в api
func NewApiMetaCustomerTable(total int64) *ApiMetaCustomerTable {
	return &ApiMetaCustomerTable{
		Total: total,
	}
}

func NewApiShortCustomerTable(name string, typevalue int, unitid int64) *ApiShortCustomerTable {
	return &ApiShortCustomerTable{
		Name:   name,
		Type:   typevalue,
		UnitID: unitid,
	}
}

func NewApiMiddleCustomerTable(id int64, name string, typevalue int) *ApiMiddleCustomerTable {
	return &ApiMiddleCustomerTable{
		ID:   id,
		Name: name,
		Type: typevalue,
	}
}

func NewApiLongCustomerTable(id int64, name string, typevalue int, unitid int64) *ApiLongCustomerTable {
	return &ApiLongCustomerTable{
		ID:     id,
		Name:   name,
		Type:   typevalue,
		UnitID: unitid,
	}
}

func NewApiSearchCustomerTable(id int64, name string, numofrows int64, created time.Time, signature string, typevalue int,
	unitid int64) *ApiSearchCustomerTable {
	return &ApiSearchCustomerTable{
		ID:        id,
		Name:      name,
		NumOfRows: numofrows,
		Created:   created,
		Signature: signature,
		Type:      typevalue,
		UnitID:    unitid,
	}
}

func NewApiFullMetaCustomerTable(numofrows int64, created time.Time, signature string, numofcols int64, checked bool, qualityper byte,
	numofwrongrows int64) *ApiFullMetaCustomerTable {
	return &ApiFullMetaCustomerTable{
		NumOfRows:      numofrows,
		Created:        created,
		Signature:      signature,
		NumOfCols:      numofcols,
		Checked:        checked,
		QaulityPer:     qualityper,
		NumOfWrongRows: numofwrongrows,
	}
}

// Конструктор создания объекта пользовательской таблицы в бд
func NewDtoCustomerTable(id int64, name string, created time.Time, typeid int, unitid int64, active bool, permanent bool,
	import_ready bool, import_percentage byte, import_columns int64, import_rows int64, import_wrongrows int64, signature string,
	import_error bool, import_errordescription string) *DtoCustomerTable {
	return &DtoCustomerTable{
		ID:                      id,
		Name:                    name,
		Created:                 created,
		TypeID:                  typeid,
		UnitID:                  unitid,
		Active:                  active,
		Permanent:               permanent,
		Import_Ready:            import_ready,
		Import_Percentage:       import_percentage,
		Import_Columns:          import_columns,
		Import_Rows:             import_rows,
		Import_WrongRows:        import_wrongrows,
		Signature:               signature,
		Import_Error:            import_error,
		Import_ErrorDescription: import_errordescription,
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
	case "type":
		fallthrough
	case "unitId":
		fallthrough
	case "rows":
		precision := 64
		if infield == "type" {
			precision = 32
		}
		_, errConv := strconv.ParseInt(invalue, 0, precision)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue

	case "name":
		fallthrough
	case "created":
		fallthrough
	case "createdName":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (table *TableSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(table)
}
