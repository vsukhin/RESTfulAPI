package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения банка компании
type ViewApiCompanyBank struct {
	Primary              bool   `json:"primary" db:"primary"`                                                    // Основной
	Bik                  string `json:"bik" db:"bik" validate:"min=1,max=25"`                                    // БИК
	Name                 string `json:"name" db:"name" validate:"min=1,max=255"`                                 // Наименование
	CheckingAccount      string `json:"checkingAccount" db:"checking_account" validate:"min=1,max=50"`           // Расчетный счет
	CorrespondingAccount string `json:"correspondentAccount" db:"corresponding_account" validate:"min=1,max=50"` // Корреспондентский счет
	Deleted              bool   `json:"del" db:"del"`                                                            // Активный
}

type DtoCompanyBank struct {
	ID                   int64  `db:"id"`                    // Уникальный идентификатор банка компании
	Company_ID           int64  `db:"company_id"`            // Идентификатор компании
	Primary              bool   `db:"primary"`               // Основной
	Bik                  string `db:"bik"`                   // БИК
	Name                 string `db:"name"`                  // Наименование
	CheckingAccount      string `db:"checking_account"`      // Расчетный счет
	CorrespondingAccount string `db:"corresponding_account"` // Корреспондентский счет
	Active               bool   `db:"active"`                // Активный
}

// Конструктор создания объекта банка компании в api
func NewViewApiCompanyBank(primary bool, bik string, name string, checkingaccount string, correspondingaccount string, deleted bool) *ViewApiCompanyBank {
	return &ViewApiCompanyBank{
		Primary:              primary,
		Bik:                  bik,
		Name:                 name,
		CheckingAccount:      checkingaccount,
		CorrespondingAccount: correspondingaccount,
		Deleted:              deleted,
	}
}

// Конструктор создания объекта банка компании в бд
func NewDtoCompanyBank(id int64, company_id int64, primary bool, bik string, name string, checkingaccount string, correspondingaccount string,
	active bool) *DtoCompanyBank {
	return &DtoCompanyBank{
		ID:                   id,
		Company_ID:           company_id,
		Primary:              primary,
		Bik:                  bik,
		Name:                 name,
		CheckingAccount:      checkingaccount,
		CorrespondingAccount: correspondingaccount,
		Active:               active,
	}
}

func (bank *ViewApiCompanyBank) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(bank, errors, req)
}
