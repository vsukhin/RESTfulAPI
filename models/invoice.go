package models

import (
	"errors"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"time"
)

// Структура для организации хранения счета
type ViewInvoice struct {
	Company_ID int64   `json:"organisationId"`         // Идентификатор компании
	Total      float64 `json:"total" validate:"min=0"` // Всего
}

type ApiMetaInvoice struct {
	Total     int64 `json:"count"`         // Всего
	Unpaid    int64 `json:"unpaid"`        // Число неоплаченных
	Companies int64 `json:"organisations"` // Число компаний
	Deleted   int64 `json:"deleted"`       // Число удаленных
}

type ApiShortInvoice struct {
	ID         int64   `json:"id" db:"id"`                         // Уникальный идентификатор счета
	Company_ID int64   `json:"organisationId" db:"organisationId"` // Идентификатор компании
	Total      float64 `json:"total" db:"total"`                   // Всего
	Paid       bool    `json:"paid" db:"paid"`                     // Оплачен
	Deleted    bool    `json:"del" db:"del"`                       // Удален
}

type ApiLongInvoice struct {
	Company_ID   int64            `json:"organisationId" db:"company_id"` // Идентификатор компании
	VAT          float64          `json:"vat" db:"vat"`                   // НДС
	Total        float64          `json:"total" db:"total"`               // Всего
	InvoiceItems []ApiInvoiceItem `json:"goods,omitempty" db:"-"`         // Позиции счета
	Paid         bool             `json:"paid" db:"paid"`                 // Оплачен
	Deleted      bool             `json:"del" db:"del"`                   // Удален
}

type ApiFullInvoice struct {
	ID           int64            `json:"id" db:"id"`                     // Уникальный идентификатор счета
	Company_ID   int64            `json:"organisationId" db:"company_id"` // Идентификатор компании
	VAT          float64          `json:"vat" db:"vat"`                   // НДС
	Total        float64          `json:"total" db:"total"`               // Всего
	InvoiceItems []ApiInvoiceItem `json:"goods,omitempty" db:"-"`         // Позиции счета
	Paid         bool             `json:"paid" db:"paid"`                 // Оплачен
	Deleted      bool             `json:"del" db:"del"`                   // Удален
}

type InvoiceSearch struct {
	ID         int64   `query:"id" search:"id"`                     // Уникальный идентификатор счета
	Company_ID int64   `query:"organisationId" search:"company_id"` // Идентификатор компании
	Total      float64 `query:"total" search:"total"`               // Всего
	Paid       bool    `query:"paid" search:"paid"`                 // Оплачен
	Deleted    bool    `query:"del" search:"(not active)"`          // Удален

}

type DtoInvoice struct {
	ID           int64            `db:"id"`         // Уникальный идентификатор счета
	Company_ID   int64            `db:"company_id"` // Идентификатор компании
	VAT          float64          `db:"vat"`        // НДС
	Total        float64          `db:"total"`      // Всего
	InvoiceItems []DtoInvoiceItem `db:"-"`          // Позиции счета
	Paid         bool             `db:"paid"`       // Оплачен
	Created      time.Time        `db:"created"`    // Время создания
	Active       bool             `db:"active"`     // Aктивен
}

// Конструктор создания объекта счета в api
func NewApiMetaInvoice(total int64, unpaid int64, companies int64, deleted int64) *ApiMetaInvoice {
	return &ApiMetaInvoice{
		Total:     total,
		Unpaid:    unpaid,
		Companies: companies,
		Deleted:   deleted,
	}
}

func NewApiShortInvoice(id int64, company_id int64, vat float64, total float64, paid bool, deleted bool) *ApiShortInvoice {
	return &ApiShortInvoice{
		ID:         id,
		Company_ID: company_id,
		Total:      total,
		Paid:       paid,
		Deleted:    deleted,
	}
}

func NewApiLongInvoice(company_id int64, vat float64, total float64, invoiceitems []ApiInvoiceItem, paid bool, deleted bool) *ApiLongInvoice {
	return &ApiLongInvoice{
		Company_ID:   company_id,
		VAT:          vat,
		Total:        total,
		InvoiceItems: invoiceitems,
		Paid:         paid,
		Deleted:      deleted,
	}
}

func NewApiFullInvoice(id int64, company_id int64, vat float64, total float64, invoiceitems []ApiInvoiceItem, paid bool, deleted bool) *ApiFullInvoice {
	return &ApiFullInvoice{
		ID:           id,
		Company_ID:   company_id,
		VAT:          vat,
		Total:        total,
		InvoiceItems: invoiceitems,
		Paid:         paid,
		Deleted:      deleted,
	}
}

// Конструктор создания объекта счета в бд
func NewDtoInvoice(id int64, company_id int64, vat float64, total float64, invoiceitems []DtoInvoiceItem, paid bool, created time.Time, active bool) *DtoInvoice {
	return &DtoInvoice{
		ID:           id,
		Company_ID:   company_id,
		VAT:          vat,
		Total:        total,
		InvoiceItems: invoiceitems,
		Paid:         paid,
		Created:      created,
		Active:       active,
	}
}

func (invoice *InvoiceSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, invoice), nil
}

func (invoice *InvoiceSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, invoice)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "organisationId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "total":
		_, errConv := strconv.ParseFloat(invalue, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "paid":
		fallthrough
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

func (invoice *InvoiceSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(invoice)
}

func (invoice *ViewInvoice) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(invoice, errors, req)
}
