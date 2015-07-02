package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения поставщика отчета
type ViewApiReportSupplier struct {
	Supplier_ID int64 `json:"id" db:"supplier_id" validate:"nonzero"` // Идентификатор поставщика
}

type DtoReportSupplier struct {
	Report_ID   int64 `db:"report_id"`   // Идентификатор отчета
	Supplier_ID int64 `db:"supplier_id"` // Идентификатор поставщика
}

// Конструктор создания объекта проекта отчета в api
func NewViewApiReportSupplier(supplier_id int64) *ViewApiReportSupplier {
	return &ViewApiReportSupplier{
		Supplier_ID: supplier_id,
	}
}

// Конструктор создания объекта поставщика отчета в бд
func NewDtoReportSupplier(report_id int64, supplier_id int64) *DtoReportSupplier {
	return &DtoReportSupplier{
		Report_ID:   report_id,
		Supplier_ID: supplier_id,
	}
}

func (supplier *ViewApiReportSupplier) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(supplier, errors, req)
}
