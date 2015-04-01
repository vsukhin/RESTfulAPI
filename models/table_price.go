package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

//Структура для организации хранения свойств прайс-листов
type DtoPriceProperties struct {
	Customer_Table_ID int64     `db:"customer_table_id"` // Идентификатор пользовательской таблицы
	Facility_ID       int64     `db:"service_id"`        // Идентификатор сервиса
	After_ID          int64     `db:"after_id"`          // Идентификатор предыдущего прайс-листа
	Begin             time.Time `db:"begin"`             // Время начало действия прайс-листа
	End               time.Time `db:"end"`               // Время окончания действия прайс-листа
	Created           time.Time `db:"created"`           // Время создания
}

type ViewApiPriceProperties struct {
	Facility_ID int64     `json:"serviceId" db:"service_id" validate:"nonzero"` // Уникальный идентификатор сервиса
	After_ID    int64     `json:"afterPriceId" db:"after_id"`                   // Уникальный идентификатор предыдущего прайс-листа
	Begin       time.Time `json:"begin" db:"begin"`                             // Время начало действия прайс-листа
	End         time.Time `json:"end" db:"end"`                                 // Время окончания действия прайс-листа
}

// Конструктор создания объекта свойств прайс-листа в api
func NewViewApiPriceProperties(facility_id int64, after_id int64, begin time.Time, end time.Time) *ViewApiPriceProperties {
	return &ViewApiPriceProperties{
		Facility_ID: facility_id,
		After_ID:    after_id,
		Begin:       begin,
		End:         end,
	}
}

// Конструктор создания объекта свойств прайс-листа в бд
func NewDtoPriceProperties(customer_table_id int64, facility_id int64, after_id int64,
	begin time.Time, end time.Time, created time.Time) *DtoPriceProperties {
	return &DtoPriceProperties{
		Customer_Table_ID: customer_table_id,
		Facility_ID:       facility_id,
		After_ID:          after_id,
		Begin:             begin,
		End:               end,
		Created:           created,
	}
}

func (priceproperties *ViewApiPriceProperties) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(priceproperties, errors, req)
}
