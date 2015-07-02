package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения события sms
type ViewApiSMSEvent struct {
	Event_ID int `json:"id" db:"event_id" validate:"nonzero"` // Идентификатор события
}

type DtoSMSEvent struct {
	Order_ID int64 `db:"order_id"` // Идентификатор заказа
	Event_ID int   `db:"event_id"` // Идентификатор события
}

// Конструктор создания объекта события sms в api
func NewViewApiSMSEvent(event_id int) *ViewApiSMSEvent {
	return &ViewApiSMSEvent{
		Event_ID: event_id,
	}
}

// Конструктор создания объекта события sms в бд
func NewDtoSMSEvent(order_id int64, event_id int) *DtoSMSEvent {
	return &DtoSMSEvent{
		Order_ID: order_id,
		Event_ID: event_id,
	}
}

func (event *ViewApiSMSEvent) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(event, errors, req)
}
