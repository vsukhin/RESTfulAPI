package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения заказа отчета
type ViewApiReportOrder struct {
	Order_ID int64 `json:"id" db:"order_id" validate:"nonzero"` // Идентификатор заказа
}

type DtoReportOrder struct {
	Report_ID int64 `db:"report_id"` // Идентификатор отчета
	Order_ID  int64 `db:"order_id"`  // Идентификатор заказа
}

// Конструктор создания объекта заказа отчета в api
func NewViewApiReportOrder(order_id int64) *ViewApiReportOrder {
	return &ViewApiReportOrder{
		Order_ID: order_id,
	}
}

// Конструктор создания объекта заказа отчета в бд
func NewDtoReportOrder(report_id int64, order_id int64) *DtoReportOrder {
	return &DtoReportOrder{
		Report_ID: report_id,
		Order_ID:  order_id,
	}
}

func (order *ViewApiReportOrder) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(order, errors, req)
}
