package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения статуса заказа отчета
type ViewApiReportComplexStatus struct {
	ComplexStatus_ID int `json:"id" db:"complex_status_id" validate:"nonzero"` // Идентификатор статуса
}

type DtoReportComplexStatus struct {
	Report_ID        int64 `db:"report_id"`         // Идентификатор отчета
	ComplexStatus_ID int   `db:"complex_status_id"` // Идентификатор статуса
}

// Конструктор создания объекта заказа отчета в api
func NewViewApiReportComplexStatus(complexstatus_id int) *ViewApiReportComplexStatus {
	return &ViewApiReportComplexStatus{
		ComplexStatus_ID: complexstatus_id,
	}
}

// Конструктор создания объекта заказа отчета в бд
func NewDtoReportComplexStatus(report_id int64, complexstatus_id int) *DtoReportComplexStatus {
	return &DtoReportComplexStatus{
		Report_ID:        report_id,
		ComplexStatus_ID: complexstatus_id,
	}
}

func (complexstatus *ViewApiReportComplexStatus) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(complexstatus, errors, req)
}
