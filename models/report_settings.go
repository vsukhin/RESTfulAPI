package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения настроек отчета
type ViewApiReportSettings struct {
	Field string `json:"field" db:"field" validate:"min=0,max=255"`  // Поле сортировки
	Order string `json:"orderBy" db:"order" validate:"min=0,max=15"` // Порядок сортировки
	Page  int64  `json:"page" db:"page" validate:"min=0"`            // Номер страницы
	Count int64  `json:"count" db:"count"`                           // Записей на странице
}

type DtoReportSettings struct {
	Report_ID int64  `db:"report_id"` // Идентификатор отчета
	Field     string `db:"field"`     // Поле сортировки
	Order     string `db:"order"`     // Порядок сортировки
	Page      int64  `db:"page"`      // Страница
	Count     int64  `db:"count"`     // Записей на странице
}

// Конструктор создания объекта настроек отчета в api
func NewViewApiReportSettings(field string, order string, page int64, count int64) *ViewApiReportSettings {
	return &ViewApiReportSettings{
		Field: field,
		Order: order,
		Page:  page,
		Count: count,
	}
}

// Конструктор создания объекта настроек отчета в бд
func NewDtoReportSettings(report_id int64, field string, order string, page int64, count int64) *DtoReportSettings {
	return &DtoReportSettings{
		Report_ID: report_id,
		Field:     field,
		Order:     order,
		Page:      page,
		Count:     count,
	}
}

func (settings *ViewApiReportSettings) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(settings, errors, req)
}
