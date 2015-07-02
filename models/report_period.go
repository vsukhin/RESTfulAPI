package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для хранения периода отчета
type ViewApiReportPeriod struct {
	Begin string `json:"begin" db:"begin"` // Начало периода
	End   string `json:"end" db:"end"`     // Окончание периода
}

type DtoReportPeriod struct {
	ID        int64     `db:"id"`        // Уникальный идентификатор периода
	Report_ID int64     `db:"report_id"` // Идентификатор отчета
	Begin     time.Time `db:"begin"`     // Начало периода
	End       time.Time `db:"end"`       // Окончание периода
}

// Конструктор создания объекта периода отчета в api
func NewViewApiReportPeriod(begin string, end string) *ViewApiReportPeriod {
	return &ViewApiReportPeriod{
		Begin: begin,
		End:   end,
	}
}

// Конструктор создания объекта периода отчета в бд
func NewDtoReportPeriod(id int64, report_id int64, begin time.Time, end time.Time) *DtoReportPeriod {
	return &DtoReportPeriod{
		ID:        id,
		Report_ID: report_id,
		Begin:     begin,
		End:       end,
	}
}

func (period *ViewApiReportPeriod) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(period, errors, req)
}
