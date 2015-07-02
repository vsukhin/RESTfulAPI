package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения позиции прайс листа ввода данных
type ViewApiInputProduct struct {
	Product_ID int `json:"fieldTypeId" db:"product_id" validate:"nonzero"` // Идентификатор позиции
}

type DtoInputProduct struct {
	Order_ID   int64 `db:"order_id"`   // Идентификатор заказа
	Product_ID int   `db:"product_id"` // Идентификатор позиции
}

// Конструктор создания объекта позиции прайс листа ввода данных в api
func NewViewApiInputProduct(product_id int) *ViewApiInputProduct {
	return &ViewApiInputProduct{
		Product_ID: product_id,
	}
}

// Конструктор создания объекта позиции прайс листа ввода данных в бд
func NewDtoInputProduct(order_id int64, product_id int) *DtoInputProduct {
	return &DtoInputProduct{
		Order_ID:   order_id,
		Product_ID: product_id,
	}
}

func (field *ViewApiInputProduct) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(field, errors, req)
}
