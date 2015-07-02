package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения периода sms
type ViewApiSMSPeriod struct {
	Period_ID int `json:"id" db:"period_id" validate:"nonzero"` // Идентификатор периода
}

type DtoSMSPeriod struct {
	Order_ID  int64 `db:"order_id"`  // Идентификатор заказа
	Period_ID int   `db:"period_id"` // Идентификатор периода
}

// Конструктор создания объекта периода sms в api
func NewViewApiSMSPeriod(period_id int) *ViewApiSMSPeriod {
	return &ViewApiSMSPeriod{
		Period_ID: period_id,
	}
}

// Конструктор создания объекта периода sms в бд
func NewDtoSMSPeriod(order_id int64, period_id int) *DtoSMSPeriod {
	return &DtoSMSPeriod{
		Order_ID:  order_id,
		Period_ID: period_id,
	}
}

func (period *ViewApiSMSPeriod) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(period, errors, req)
}
