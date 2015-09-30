package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для организации хранения сервиса имени отправителя заказа
type ViewHeaderFacility struct {
	CreateBegin string `json:"createBegin" validate:"max=255"` // Планируемое начало
	CreateEnd   string `json:"createEnd" validate:"max=255"`   // Планируемое окончание
	Name        string `json:"name" validate:"min=1,max=255"`  // Название
	Renew       bool   `json:"autoRenew"`                      // Продлевать
}

type ApiHeaderFacility struct {
	CreateBegin string  `json:"createBegin" db:"createBegin"` // Планируемое начало
	CreateEnd   string  `json:"createEnd" db:"createEnd"`     // Планируемое окончание
	Name        string  `json:"name" db:"name"`               // Название
	Begin       string  `json:"begin" db:"begin"`             // Фактическое начало
	End         string  `json:"end" db:"end"`                 // Фактическое окончание
	AutoRenew   bool    `json:"autoRenew" db:"autoRenew"`     // Продлевать
	Cost        float64 `json:"cost" db:"cost"`               // Сумма заказа исходя из расчётных показателей заказа
	CostFactual float64 `json:"costFactual" db:"costFactual"` // Текущая стоимость заказа
}

type DtoHeaderFacility struct {
	Order_ID    int64     `db:"order_id"`    // Идентификатор заказа
	CreateBegin time.Time `db:"createBegin"` // Планируемое начало
	CreateEnd   time.Time `db:"createEnd"`   // Планируемое окончание
	Name        string    `db:"name"`        // Название
	Begin       time.Time `db:"begin"`       // Фактическое начало
	End         time.Time `db:"end"`         // Фактическое окончание
	AutoRenew   bool      `db:"autoRenew"`   // Продлевать
	Cost        float64   `db:"cost"`        // Сумма заказа исходя из расчётных показателей заказа
	CostFactual float64   `db:"costFactual"` // Текущая стоимость заказа
}

// Конструктор создания объекта сервиса имени отправителя заказа в api
func NewApiHeaderFacility(createbegin, createend, name, begin, end string, autorenew bool, cost float64,
	costFactual float64) *ApiHeaderFacility {
	return &ApiHeaderFacility{
		CreateBegin: createbegin,
		CreateEnd:   createend,
		Name:        name,
		Begin:       begin,
		End:         end,
		AutoRenew:   autorenew,
		Cost:        cost,
		CostFactual: costFactual,
	}
}

// Конструктор создания объекта сервиса имени отправителя заказа в бд
func NewDtoHeaderFacility(order_id int64, createbegin, createend time.Time, name string, begin, end time.Time, autorenew bool,
	cost float64, costFactual float64) *DtoHeaderFacility {
	return &DtoHeaderFacility{
		Order_ID:    order_id,
		CreateBegin: createbegin,
		CreateEnd:   createend,
		Name:        name,
		Begin:       begin,
		End:         end,
		AutoRenew:   autorenew,
		Cost:        cost,
		CostFactual: costFactual,
	}
}

func (facility *ViewHeaderFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(facility, errors, req)
}
