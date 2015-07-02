package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

type BudgetedBy byte

const (
	TYPE_BUDGETEDBY_UNKNOWN BudgetedBy = iota
	TYPE_BUDGETEDBY_FACILITY
	TYPE_BUDGETEDBY_COMPLEX_STATUS
	TYPE_BUDGETEDBY_SUPPLIER
)

const (
	TYPE_BUDGETEDBY_FACILITY_VALUE       = "service"
	TYPE_BUDGETEDBY_COMPLEX_STATUS_VALUE = "orderstatus"
	TYPE_BUDGETEDBY_SUPPLIER_VALUE       = "supplier"
)

// Структура для организации хранения отчета
type ViewReport struct {
	Periods         []ViewApiReportPeriod        `json:"periods"`                     // Периоды
	Projects        []ViewApiReportProject       `json:"projects"`                    // Проекты
	Orders          []ViewApiReportOrder         `json:"orders"`                      // Заказы
	Budgeted        string                       `json:"budgetBy" validate:"nonzero"` // Вид бюджетированния
	Facilities      []ViewApiReportFacility      `json:"services"`                    // Сервисы
	ComplexStatuses []ViewApiReportComplexStatus `json:"statuses"`                    // Статусы
	Suppliers       []ViewApiReportSupplier      `json:"suppliers"`                   // Поставщики
	Settings        ViewApiReportSettings        `json:"settingsOrders"`              // Настройки
}

type ApiMetaReport struct {
	Access bool  `json:"enable" db:"reportAccess` // Доступность отчетов
	ID     int64 `json:"aggregateId" db:"id"`     // Идентификатор отчета
}

type ApiReport struct {
	ID              int64                        `json:"id" db:"id"`                 // Уникальный идентификатор отчета
	Created         string                       `json:"created" db:"created"`       // Время создания
	Unit_ID         int64                        `json:"unitId" db:"unit_id"`        // Идентификатор объединения
	User_ID         int64                        `json:"userId" db:"user_id"`        // Идентификатор пользователя
	Periods         []ViewApiReportPeriod        `json:"periods,omitempty" db:"-"`   // Периоды
	Projects        []ViewApiReportProject       `json:"projects,omitempty" db:"-"`  // Проекты
	Orders          []ViewApiReportOrder         `json:"orders,omitempty" db:"-"`    // Заказы
	Budgeted        string                       `json:"budgetBy" db:"budgeted"`     // Вид бюджетированния
	Facilities      []ViewApiReportFacility      `json:"services,omitempty" db:"-"`  // Сервисы
	ComplexStatuses []ViewApiReportComplexStatus `json:"statuses,omitempty" db:"-"`  // Статусы
	Suppliers       []ViewApiReportSupplier      `json:"suppliers,omitempty" db:"-"` // Поставщики
	Settings        ViewApiReportSettings        `json:"settingsOrders" db:"-"`      // Настройки
}

type DtoReport struct {
	ID              int64                    `db:"id"`       // Уникальный идентификатор отчета
	User_ID         int64                    `db:"user_id"`  // Идентификатор пользователя
	Unit_ID         int64                    `db:"unit_id"`  // Идентификатор объединения
	Periods         []DtoReportPeriod        `db:"-"`        // Периоды
	Projects        []DtoReportProject       `db:"-"`        // Проекты
	Orders          []DtoReportOrder         `db:"-"`        // Заказы
	Budgeted        BudgetedBy               `db:"budgeted"` // Вид бюджетированния
	Facilities      []DtoReportFacility      `db:"-"`        // Сервисы
	ComplexStatuses []DtoReportComplexStatus `db:"-"`        // Статусы
	Suppliers       []DtoReportSupplier      `db:"-"`        // Поставщики
	Settings        DtoReportSettings        `db:"-"`        // Настройки
	Created         time.Time                `db:"created"`  // Время создания
	Active          bool                     `db:"active"`   // Aктивен
}

// Конструктор создания объекта отчета в api
func NewApiMetaReport(access bool, id int64) *ApiMetaReport {
	return &ApiMetaReport{
		Access: access,
		ID:     id,
	}
}

func NewApiReport(id int64, created string, unit_id int64, user_id int64, periods []ViewApiReportPeriod,
	projects []ViewApiReportProject, orders []ViewApiReportOrder, budgeted string, facilities []ViewApiReportFacility,
	complexstatuses []ViewApiReportComplexStatus, suppliers []ViewApiReportSupplier, settings ViewApiReportSettings) *ApiReport {
	return &ApiReport{
		ID:              id,
		Created:         created,
		Unit_ID:         unit_id,
		User_ID:         user_id,
		Periods:         periods,
		Projects:        projects,
		Orders:          orders,
		Budgeted:        budgeted,
		Facilities:      facilities,
		ComplexStatuses: complexstatuses,
		Suppliers:       suppliers,
		Settings:        settings,
	}
}

// Конструктор создания объекта отчета в бд
func NewDtoReport(id int64, user_id int64, unit_id int64, periods []DtoReportPeriod, projects []DtoReportProject,
	orders []DtoReportOrder, budgeted BudgetedBy, facilities []DtoReportFacility, complexstatuses []DtoReportComplexStatus,
	suppliers []DtoReportSupplier, settings DtoReportSettings, created time.Time, active bool) *DtoReport {
	return &DtoReport{
		ID:              id,
		User_ID:         user_id,
		Unit_ID:         unit_id,
		Periods:         periods,
		Projects:        projects,
		Orders:          orders,
		Budgeted:        budgeted,
		Facilities:      facilities,
		ComplexStatuses: complexstatuses,
		Suppliers:       suppliers,
		Settings:        settings,
		Created:         created,
		Active:          active,
	}
}

func (report *ViewReport) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, period := range report.Periods {
		errors = Validate(&period, errors, req)
	}
	for _, project := range report.Projects {
		errors = Validate(&project, errors, req)
	}
	for _, order := range report.Orders {
		errors = Validate(&order, errors, req)
	}
	for _, facility := range report.Facilities {
		errors = Validate(&facility, errors, req)
	}
	for _, complexstatus := range report.ComplexStatuses {
		errors = Validate(&complexstatus, errors, req)
	}
	for _, supplier := range report.Suppliers {
		errors = Validate(&supplier, errors, req)
	}

	return Validate(report, errors, req)
}
