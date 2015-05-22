package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

const (
	CODE_FIELD_MAX_LENGTH_VALUE = 255
)

// Структура для организации хранения кода компании
type ViewApiCompanyCode struct {
	Company_Class_ID int    `json:"organisationClassesId" db:"company_class_id" validate:"nonzero"` // Идентификатор класса компании
	Codes            string `json:"value" db:"codes"`                                               // Коды
}

type DtoCompanyCode struct {
	ID               int64  `db:"id"`               // Уникальный идентификатор кода компании
	Company_ID       int64  `db:"company_id"`       // Идентификатор компании
	Company_Class_ID int    `db:"company_class_id"` // Идентификатор класса компании
	Code             string `db:"code"`             // Код
}

// Конструктор создания объекта кода компании в api
func NewViewApiCompanyCode(company_class_id int, codes string) *ViewApiCompanyCode {
	return &ViewApiCompanyCode{
		Company_Class_ID: company_class_id,
		Codes:            codes,
	}
}

// Конструктор создания объекта кода компании в бд
func NewDtoCompanyCode(id int64, company_id int64, company_class_id int, code string) *DtoCompanyCode {
	return &DtoCompanyCode{
		ID:               id,
		Company_ID:       company_id,
		Company_Class_ID: company_class_id,
		Code:             code,
	}
}

func (code *ViewApiCompanyCode) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(code, errors, req)
}
