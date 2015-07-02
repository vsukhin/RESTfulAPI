package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения позиции прайс листа верификации данных
type ViewApiDataProduct struct {
	Product_ID int `json:"id" db:"product_id" validate:"nonzero"` // Идентификатор позиции
}

type DtoDataProduct struct {
	Order_ID   int64 `db:"order_id"`   // Идентификатор заказа
	Product_ID int   `db:"product_id"` // Идентификатор позиции
}

// Конструктор создания объекта позиции прайс листа верификации данных в api
func NewViewApiDataProduct(product_id int) *ViewApiDataProduct {
	return &ViewApiDataProduct{
		Product_ID: product_id,
	}
}

// Конструктор создания объекта позиции прайс листа верификации данных в бд
func NewDtoDataProduct(order_id int64, product_id int) *DtoDataProduct {
	return &DtoDataProduct{
		Order_ID:   order_id,
		Product_ID: product_id,
	}
}

func (field *ViewApiDataProduct) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(field, errors, req)
}
