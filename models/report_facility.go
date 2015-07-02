package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения сервиса отчета
type ViewApiReportFacility struct {
	Facility_ID int64 `json:"id" db:"service_id" validate:"nonzero"` // Идентификатор сервиса
}

type DtoReportFacility struct {
	Report_ID   int64 `db:"report_id"`  // Идентификатор отчета
	Facility_ID int64 `db:"service_id"` // Идентификатор сервиса
}

// Конструктор создания объекта проекта отчета в api
func NewViewApiReportFacility(facility_id int64) *ViewApiReportFacility {
	return &ViewApiReportFacility{
		Facility_ID: facility_id,
	}
}

// Конструктор создания объекта проекта отчета в бд
func NewDtoReportFacility(report_id int64, facility_id int64) *DtoReportFacility {
	return &DtoReportFacility{
		Report_ID:   report_id,
		Facility_ID: facility_id,
	}
}

func (facility *ViewApiReportFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(facility, errors, req)
}
