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

// Структура для организации хранения отправителя
type ViewSMSSender struct {
	Planned_End string `json:"createEnd" validate:"max=255"` // Планируемое окончание
	Renew       bool   `json:"autoRenew"`                    // Продлевать
}

type ApiMetaSMSSenderPerSupplier struct {
	Supplier_ID int64 `json:"supplierId" db:"supplier_id"` // Идентификатор поставщика
	Total       int64 `json:"count" db:"total"`            // Число действующих отправителей
}

type ApiMetaSMSSender struct {
	Total        int64                         `json:"count"`       // Число действующих отправителей
	NumOfDeleted int64                         `json:"deleted"`     // Число удаленных отправителей
	NumOfNew     int64                         `json:"new"`         // Число незарегистрированных отправителей
	Suppliers    []ApiMetaSMSSenderPerSupplier `json:"bySuppliers"` // По поставщикам
}

type ApiShortSMSSender struct {
	ID            int64  `json:"id" db:"id"`                     // Уникальный идентификатор отправителя
	Name          string `json:"name" db:"name"`                 // Название
	Supplier_ID   int64  `json:"supplierId" db:"suppplier_id"`   // Идентификатор поставщика
	Registered    bool   `json:"registered" db:"registered"`     // Зарегистрирован
	Planned_Begin string `json:"createBegin" db:"planned_begin"` // Планируемое начало
	Planned_End   string `json:"createEnd" db:"planned_end"`     // Планируемое окончание
	Actual_Begin  string `json:"begin" db:"actual_begin"`        // Фактическое начало
	Actual_End    string `json:"end" db:"actual_end"`            // Фактическое окончание
	Renew         bool   `json:"autoRenew" db:"renew"`           // Продлевать
}

type ApiLongSMSSender struct {
	ID            int64  `json:"id" db:"id"`                   // Уникальный идентификатор отправителя
	Name          string `json:"name" db:"name"`               // Название
	Supplier_ID   int64  `json:"supplierId" db:"supplierId"`   // Идентификатор поставщика
	Registered    bool   `json:"registered" db:"registered"`   // Зарегистрирован
	Planned_Begin string `json:"createBegin" db:"createBegin"` // Планируемое начало
	Planned_End   string `json:"createEnd" db:"createEnd"`     // Планируемое окончание
	Actual_Begin  string `json:"begin" db:"begin"`             // Фактическое начало
	Actual_End    string `json:"end" db:"end"`                 // Фактическое окончание
	Renew         bool   `json:"autoRenew" db:"autoRenew"`     // Продлевать
	Deleted       bool   `json:"del" db:"del"`                 // Удален
}

type SMSSenderSearch struct {
	ID            int64  `query:"id" search:"id"`                                                                                        // Уникальный идентификатор отправителя
	Name          string `query:"name" search:"name" group:"name"`                                                                       // Название
	Supplier_ID   int64  `query:"supplierId" search:"supplier_id"`                                                                       // Идентификатор поставщика
	Registered    bool   `query:"registered" search:"registered"`                                                                        // Зарегистрирован
	Planned_Begin string `query:"createBegin" search:"planned_begin" group:"convert(date_format(planned_begin, '%Y-%m-%d') using utf8)"` // Планируемое начало
	Planned_End   string `query:"createEnd" search:"planned_end" group:"convert(date_format(planned_end, '%Y-%m-%d') using utf8)"`       // Планируемое окончание
	Actual_Begin  string `query:"begin" search:"actual_begin" group:"convert(date_format(actual_begin, '%Y-%m-%d') using utf8)"`         // Фактическое начало
	Actual_End    string `query:"end" search:"actual_end" group:"convert(date_format(actual_end, '%Y-%m-%d') using utf8)"`               // Фактическое окончание
	Renew         bool   `query:"autoRenew" search:"renew"`                                                                              // Продлевать
	Deleted       bool   `query:"del" search:"(not active)"`                                                                             // Удален
}

type DtoSMSSender struct {
	ID              int64     `db:"id"`              // Уникальный идентификатор отправителя
	Unit_ID         int64     `db:"unit_id"`         // Идентификатор объединения
	Name            string    `db:"name"`            // Название
	Created         time.Time `db:"created"`         // Время создания
	Registered      bool      `db:"registered"`      // Зарегистрирован
	Active          bool      `db:"active"`          // Aктивен
	Planned_Begin   time.Time `db:"planned_begin"`   // Планируемое начало
	Planned_End     time.Time `db:"planned_end"`     // Планируемое окончание
	Actual_Begin    time.Time `db:"actual_begin"`    // Фактическое начало
	Actual_End      time.Time `db:"actual_end"`      // Фактическое окончание
	Rejected        bool      `db:"rejected"`        // Отказано
	Rejected_Reason string    `db:"rejected_reason"` // Причина отказа
	Supplier_ID     int64     `db:"supplier_id"`     // Идентификатор поставщика
	Renew           bool      `db:"renew"`           // Продлевать
}

// Конструктор создания объекта отправителя в api
func NewApiMetaSMSSenderPerSupplier(supplier_id int64, total int64) *ApiMetaSMSSenderPerSupplier {
	return &ApiMetaSMSSenderPerSupplier{
		Supplier_ID: supplier_id,
		Total:       total,
	}
}

func NewApiMetaSMSSender(total int64, numofdeleted int64, numofnew int64, suppliers []ApiMetaSMSSenderPerSupplier) *ApiMetaSMSSender {
	return &ApiMetaSMSSender{
		Total:        total,
		NumOfDeleted: numofdeleted,
		NumOfNew:     numofnew,
		Suppliers:    suppliers,
	}
}

func NewApiShortSMSSender(id int64, name string, supplier_id int64, registered bool, planned_begin, planned_end, actual_begin, actual_end string,
	renew bool) *ApiShortSMSSender {
	return &ApiShortSMSSender{
		ID:            id,
		Name:          name,
		Supplier_ID:   supplier_id,
		Registered:    registered,
		Planned_Begin: planned_begin,
		Planned_End:   planned_end,
		Actual_Begin:  actual_begin,
		Actual_End:    actual_end,
		Renew:         renew,
	}
}

func NewApiLongSMSSender(id int64, name string, supplier_id int64, registered bool, planned_begin, planned_end, actual_begin, actual_end string,
	renew bool, deleted bool) *ApiLongSMSSender {
	return &ApiLongSMSSender{
		ID:            id,
		Name:          name,
		Supplier_ID:   supplier_id,
		Registered:    registered,
		Planned_Begin: planned_begin,
		Planned_End:   planned_end,
		Actual_Begin:  actual_begin,
		Actual_End:    actual_end,
		Renew:         renew,
		Deleted:       deleted,
	}
}

// Конструктор создания объекта отправителя в бд
func NewDtoSMSSender(id int64, unit_id int64, name string, created time.Time, registered bool,
	active bool, planned_begin, planned_end, actual_begin, actual_end time.Time,
	rejected bool, rejected_reason string, supplier_id int64, renew bool) *DtoSMSSender {
	return &DtoSMSSender{
		ID:              id,
		Unit_ID:         unit_id,
		Name:            name,
		Created:         created,
		Registered:      registered,
		Active:          active,
		Planned_Begin:   planned_begin,
		Planned_End:     planned_end,
		Actual_Begin:    actual_begin,
		Actual_End:      actual_end,
		Rejected:        rejected,
		Rejected_Reason: rejected_reason,
		Supplier_ID:     supplier_id,
		Renew:           renew,
	}
}

func (smssender *SMSSenderSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, smssender), nil
}

func (smssender *SMSSenderSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, smssender)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "supplierId":
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
	case "registered":
		fallthrough
	case "autoRenew":
		fallthrough
	case "del":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	case "createBegin":
		fallthrough
	case "createEnd":
		fallthrough
	case "begin":
		fallthrough
	case "end":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (smssender *SMSSenderSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(smssender)
}

func (smssender *ViewSMSSender) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(smssender, errors, req)
}
