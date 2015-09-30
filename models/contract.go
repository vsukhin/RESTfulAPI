package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для организации хранения компании
type ViewContract struct {
	Company_ID int64 `json:"organisationId" db:"company_id" validate:"nonzero"` // Идентификатор компании
	Confirmed  bool  `json:"iWantToSign" db:"confirmed"`                        // Принят
}

type ChangeContract []ViewContract

type ApiContract struct {
	Company_ID int64         `json:"organisationId" db:"company_id"`      // Идентификатор компании
	Confirmed  bool          `json:"iWantToSign" db:"confirmed"`          // Принят
	Signed     bool          `json:"signed" db:"signed"`                  // Подписан
	SignedDate time.Time     `json:"signedDate" db:"signed_date"`         // Время подписания
	Name       string        `json:"contractName" db:"name"`              // Название
	File_ID    int64         `json:"contractFileId" db:"file_id"`         // Идентификатор файла
	Appendices []ApiAppendix `json:"contractAdditional,omitempty" db:"-"` // Коды компании
}

type DtoContract struct {
	ID         int64         `db:"id"`          // Уникальный идентификатор договора
	Company_ID int64         `db:"company_id"`  // Идентификатор компании
	Name       string        `db:"name"`        // Название
	Confirmed  bool          `db:"confirmed"`   // Принят
	Signed     bool          `db:"signed"`      // Подписан
	SignedDate time.Time     `db:"signed_date"` // Время подписания
	File_ID    int64         `db:"file_id"`     // Идентификатор файла
	Created    time.Time     `db:"created"`     // Время создания
	Active     bool          `db:"active"`      // Aктивен
	Appendices []DtoAppendix `db:"-"`           // Коды компании
}

// Конструктор создания объекта компании в api
func NewApiContract(company_id int64, confirmed bool, signed bool, signeddate time.Time, name string, file_id int64,
	appendices []ApiAppendix) *ApiContract {
	return &ApiContract{
		Company_ID: company_id,
		Confirmed:  confirmed,
		Signed:     signed,
		SignedDate: signeddate,
		Name:       name,
		File_ID:    file_id,
		Appendices: appendices,
	}
}

// Конструктор создания объекта компании в бд
func NewDtoContract(id int64, company_id int64, name string, confirmed bool, signed bool, signeddate time.Time, file_id int64,
	created time.Time, active bool, appendices []DtoAppendix) *DtoContract {
	return &DtoContract{
		ID:         id,
		Company_ID: company_id,
		Name:       name,
		Confirmed:  confirmed,
		Signed:     signed,
		SignedDate: signeddate,
		File_ID:    file_id,
		Created:    created,
		Active:     active,
		Appendices: appendices,
	}
}

func (contract ViewContract) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(contract, errors, req)
}
