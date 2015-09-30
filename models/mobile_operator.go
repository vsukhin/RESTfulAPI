package models

import (
	"time"
)

type BillingModel byte

const (
	BILLING_MODEL_RANGE BillingModel = iota
	BILLING_MODEL_CUMULATIVE_PERHEADER
)

const (
	MOBILE_OPERATOR_UUID_UNKNOWN = "43915643-c1c1-477e-93bb-af5855f07bf9"
	MOBILE_OPERATOR_UUID_MEGAFON = "dd18d6a6-3d4c-11e5-b16b-00185130a65b"
)

// Структура для организации хранения мобильного оператора
type ApiMobileOperator struct {
	ID        int    `json:"id" db:"id"`               // Уникальный идентификатор мобильного оператора
	ShortName string `json:"nameShort" db:"shortname"` // Короткое название
	LongName  string `json:"nameLong" db:"longname"`   // Длинное название
	Position  int    `json:"position" db:"position"`   // Позиция
}

type DtoMobileOperator struct {
	ID              int          `db:"id"`                // Уникальный идентификатор мобильного оператора
	ShortName       string       `db:"shortname"`         // Короткое название
	LongName        string       `db:"longname"`          // Длинное название
	Created         time.Time    `db:"created"`           // Время создания
	Position        int          `db:"position"`          // Позиция
	Active          bool         `db:"active"`            // Aктивен
	UUID            string       `db:"uuid"`              // UUID объединения
	IsDefault       bool         `db:"default"`           // Использовать по умолчанию
	SMSBillingModel BillingModel `db:"sms_billing_model"` // Модель биллинга SMS
	HLRBillingModel BillingModel `db:"hlr_billing_model"` // Модель биллинга HLR
}

// Конструктор создания объекта мобильного оператора в api
func NewApiMobileOperator(id int, shortname string, longname string, position int) *ApiMobileOperator {
	return &ApiMobileOperator{
		ID:        id,
		ShortName: shortname,
		LongName:  longname,
		Position:  position,
	}
}

// Конструктор создания объекта мобильного оператора в бд
func NewDtoMobileOperator(id int, shortname string, longname string, created time.Time, position int, active bool,
	uuid string, isdefault bool, smsbillingmodel, hlrbillingmodel BillingModel) *DtoMobileOperator {
	return &DtoMobileOperator{
		ID:              id,
		ShortName:       shortname,
		LongName:        longname,
		Created:         created,
		Position:        position,
		Active:          active,
		UUID:            uuid,
		IsDefault:       isdefault,
		SMSBillingModel: smsbillingmodel,
		HLRBillingModel: hlrbillingmodel,
	}
}
