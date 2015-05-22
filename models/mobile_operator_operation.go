package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения действия мобильного оператора
type ViewApiMobileOperatorOperation struct {
	MobileOperator_ID int  `json:"mobileOperatorId" db:"mobileoperator_id" validate:"nonzero"` // Идентификатор мобильного оператора
	Percent           byte `json:"percent" db:"percent" validate:"min=0,max=100"`              // Процент
	Count             int  `json:"count" db:"count" validate:"min=0"`                          // Количество
}

type DtoMobileOperatorOperation struct {
	Order_ID          int64 `db:"order_id"`          // Идентификатор заказа
	MobileOperator_ID int   `db:"mobileoperator_id"` // Идентификатор мобильного оператора
	Percent           byte  `db:"percent"`           // Процент
	Count             int   `db:"count"`             // Количество
}

// Конструктор создания объекта действия мобильного оператора в api
func NewViewApiMobileOperatorOperation(mobileoperator_id int, percent byte, count int) *ViewApiMobileOperatorOperation {
	return &ViewApiMobileOperatorOperation{
		MobileOperator_ID: mobileoperator_id,
		Percent:           percent,
		Count:             count,
	}
}

// Конструктор создания объекта действия мобильного оператора в бд
func NewDtoMobileOperatorOperation(order_id int64, mobileoperator_id int, percent byte, count int) *DtoMobileOperatorOperation {
	return &DtoMobileOperatorOperation{
		Order_ID:          order_id,
		MobileOperator_ID: mobileoperator_id,
		Percent:           percent,
		Count:             count,
	}
}

func (operation *ViewApiMobileOperatorOperation) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(operation, errors, req)
}
