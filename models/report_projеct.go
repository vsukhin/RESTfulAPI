package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения проекта отчета
type ViewApiReportProject struct {
	Project_ID int64 `json:"id" db:"project_id" validate:"nonzero"` // Идентификатор проекта
}

type DtoReportProject struct {
	Report_ID  int64 `db:"report_id"`  // Идентификатор отчета
	Project_ID int64 `db:"project_id"` // Идентификатор проекта
}

// Конструктор создания объекта проекта отчета в api
func NewViewApiReportProject(project_id int64) *ViewApiReportProject {
	return &ViewApiReportProject{
		Project_ID: project_id,
	}
}

// Конструктор создания объекта проекта отчета в бд
func NewDtoReportProject(report_id int64, project_id int64) *DtoReportProject {
	return &DtoReportProject{
		Report_ID:  report_id,
		Project_ID: project_id,
	}
}

func (project *ViewApiReportProject) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(project, errors, req)
}
