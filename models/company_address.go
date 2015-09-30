package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения адреса компании
type ViewApiCompanyAddress struct {
	Primary         bool   `json:"primary" db:"primary"`                         // Основной
	Ditto           int    `json:"ditto" db:"ditto"`                             // Идентификатор повторного типа адреса
	Address_Type_ID int    `json:"type" db:"address_type_id" validate:"nonzero"` // Идентификатор типа адреса
	Full            string `json:"allInOne" db:"full"`                           // Полный адрес
	Zip             string `json:"zipCode" db:"zip" validate:"max=50"`           // Индекс
	Country         string `json:"country" db:"country" validate:"max=100"`      // Страна
	Region          string `json:"region" db:"region" validate:"max=255"`        // Регион
	City            string `json:"city" db:"city" validate:"max=255"`            // Город
	Street          string `json:"street" db:"street" validate:"max=255"`        // Улица
	Building        string `json:"build" db:"building" validate:"max=255"`       // Дом, корпус, строение, комната
	Postbox         string `json:"poBox" db:"postbox" validate:"max=50"`         // Номер почтового ящика
	Company         string `json:"company" db:"company" validate:"max=255"`      // Компания
	Comments        string `json:"comments" db:"comments"`                       // Комментарии
}

type DtoCompanyAddress struct {
	ID              int64  `db:"id"`              // Уникальный идентификатор адреса компании
	Company_ID      int64  `db:"company_id"`      // Идентификатор компании
	Address_Type_ID int    `db:"address_type_id"` // Идентификатор типа адреса
	Ditto           int    `db:"ditto"`           // Идентификатор повторного типа адреса
	Primary         bool   `db:"primary"`         // Основной
	Zip             string `db:"zip"`             // Индекс
	Country         string `db:"country"`         // Страна
	Region          string `db:"region"`          // Регион
	City            string `db:"city"`            // Город
	Street          string `db:"street"`          // Улица
	Building        string `db:"building"`        // Дом, корпус, строение, комната
	Postbox         string `db:"postbox"`         // Номер почтового ящика
	Company         string `db:"company"`         // Компания
	Comments        string `db:"comments"`        // Комментарии
	Full            string `db:"full"`            //  Полный адрес
}

// Конструктор создания объекта адреса компании в api
func NewViewApiCompanyAddress(primary bool, ditto int, address_type_id int, full string, zip string, country string, region string,
	city string, street string, building string, postbox string, company string, comments string) *ViewApiCompanyAddress {
	return &ViewApiCompanyAddress{
		Primary:         primary,
		Ditto:           ditto,
		Address_Type_ID: address_type_id,
		Full:            full,
		Zip:             zip,
		Country:         country,
		Region:          region,
		City:            city,
		Street:          street,
		Building:        building,
		Postbox:         postbox,
		Company:         company,
		Comments:        comments,
	}
}

// Конструктор создания объекта адреса компании в бд
func NewDtoCompanyAddress(id int64, company_id int64, address_type_id int, ditto int, primary bool, zip string, country string, region string,
	city string, street string, building string, postbox string, company string, comments string, full string) *DtoCompanyAddress {
	return &DtoCompanyAddress{
		ID:              id,
		Company_ID:      company_id,
		Address_Type_ID: address_type_id,
		Ditto:           ditto,
		Primary:         primary,
		Zip:             zip,
		Country:         country,
		Region:          region,
		City:            city,
		Street:          street,
		Building:        building,
		Postbox:         postbox,
		Company:         company,
		Comments:        comments,
		Full:            full,
	}
}

func (address *ViewApiCompanyAddress) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(address, errors, req)
}
