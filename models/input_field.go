package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения прогноза количества вводимого типа поля
type ViewApiInputField struct {
	Product_ID int `json:"fieldTypeId" db:"product_id" validate:"nonzero"` // Идентификатор типа поля
	Count      int `json:"count" db:"count" validate:"min=0"`              // Количество
}

type DtoInputField struct {
	Order_ID   int64 `db:"order_id"`   // Идентификатор заказа
	Product_ID int   `db:"product_id"` // Идентификатор типа поля
	Count      int   `db:"count"`      // Количество
}

// Конструктор создания объекта прогноза количества вводимого типа поля в api
func NewViewApiInputField(product_id int, count int) *ViewApiInputField {
	return &ViewApiInputField{
		Product_ID: product_id,
		Count:      count,
	}
}

// Конструктор создания объекта прогноза количества вводимого типа поля в бд
func NewDtoInputField(order_id int64, product_id int, count int) *DtoInputField {
	return &DtoInputField{
		Order_ID:   order_id,
		Product_ID: product_id,
		Count:      count,
	}
}

func (field *ViewApiInputField) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(field, errors, req)
}
