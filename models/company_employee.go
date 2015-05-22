package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

const (
	EMPLOYEE_TYPE_CEO        = "ceo"
	EMPLOYEE_TYPE_ACCOUNTANT = "accountant"
)

// Структура для организации хранения сотрудника компании
type ViewApiCompanyEmployee struct {
	Employee_Type string `json:"type" db:"employee_type" validate:"nonzero,min=1,max=50"` // Тип сотрудника
	Ditto         string `json:"ditto" db:"ditto" validate:"max=50"`                      // Повторный тип сотрудника
	Surname       string `json:"surname" db:"surname" validate:"max=100"`                 // Фамилия
	Name          string `json:"name" db:"name" validate:"max=100"`                       // Имя
	MiddleName    string `json:"middleName" db:"middlename" validate:"max=100"`           // Отчество
	Base          string `json:"basis" db:"base" validate:"max=255"`                      // Основание
	Deleted       bool   `json:"del" db:"del"`                                            // Удален
}

type DtoCompanyEmployee struct {
	ID            int64  `db:"id"`            // Уникальный идентификатор сотрудника
	Company_ID    int64  `db:"company_id"`    // Идентификатор компании
	Employee_Type string `db:"employee_type"` // Тип сотрудника
	Ditto         string `db:"ditto"`         // Повторный тип сотрудника
	Surname       string `db:"surname"`       // Фамилия
	Name          string `db:"name"`          // Имя
	MiddleName    string `db:"middlename"`    // Отчество
	Base          string `db:"base"`          // Основание
	Active        bool   `db:"active"`        // Активный
}

// Конструктор создания объекта сотрудника компании в api
func NewViewApiCompanyEmployee(employee_type string, ditto string, surname string, name string, middlename string,
	base string, deleted bool) *ViewApiCompanyEmployee {
	return &ViewApiCompanyEmployee{
		Employee_Type: employee_type,
		Ditto:         ditto,
		Surname:       surname,
		Name:          name,
		MiddleName:    middlename,
		Base:          base,
		Deleted:       deleted,
	}
}

// Конструктор создания объекта сотрудника компании в бд
func NewDtoCompanyEmployee(id int64, company_id int64, employee_type string, ditto string, surname string, name string, middlename string,
	base string, active bool) *DtoCompanyEmployee {
	return &DtoCompanyEmployee{
		ID:            id,
		Company_ID:    company_id,
		Employee_Type: employee_type,
		Ditto:         ditto,
		Surname:       surname,
		Name:          name,
		MiddleName:    middlename,
		Base:          base,
		Active:        active,
	}
}

func (employee *ViewApiCompanyEmployee) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(employee, errors, req)
}
