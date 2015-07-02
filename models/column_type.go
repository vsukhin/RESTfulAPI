package models

import (
	"time"
)

type Alignment byte

const (
	ALIGNMENT_LEFT Alignment = iota + 1
	ALIGNMNET_RIGHT
	ALIGNMENT_CENTER
)

const (
	COLUMN_TYPE_DEFAULT                  = 0
	COLUMN_TYPE_MOBILE_PHONE             = 1
	COLUMN_TYPE_SMS                      = 2
	COLUMN_TYPE_SMS_SENDER               = 3
	COLUMN_TYPE_BIRTHDAY                 = 4
	COLUMN_TYPE_SOURCE_ADDRESS           = 5
	COLUMN_TYPE_SOURCE_PHONE             = 6
	COLUMN_TYPE_SOURCE_FIO               = 7
	COLUMN_TYPE_SOURCE_EMAIL             = 8
	COLUMN_TYPE_SOURCE_DATE              = 9
	COLUMN_TYPE_SOURCE_AUTOMOBILE        = 10
	COLUMN_TYPE_PRICELIST_NAME           = 78
	COLUMN_TYPE_PRICELIST_PRICE          = 79
	COLUMN_TYPE_PRICELIST_DISCOUNT       = 80
	COLUMN_TYPE_PRICELIST_MOBILEOPERATOR = 81
	COLUMN_TYPE_PRICELIST_RANGE          = 82
	COLUMN_TYPE_PRICELIST_ID             = 83
)

// Структура для организации хранения типов колонок
type ApiColumnType struct {
	ID               int    `json:"id" db:"id"`                       // Уникальный идентификатор типа колонки
	Name             string `json:"name" db:"name"`                   // Название
	Description      string `json:"description" db:"description"`     // Описание
	Required         bool   `json:"notNull" db:"required"`            // Обязательность к заполнению
	Regexp           string `json:"regexp" db:"regexp"`               // Регулярное выражение для проверки
	HorAlignmentHead string `json:"alignmentHead" db:"alignmentHead"` // Горизонтальное выравнивание заголовка
	HorAlignmentBody string `json:"alignmentBody" db:"alignmentBody"` // Горизонтальное выравнивание содержимого
}

type DtoColumnType struct {
	ID               int       `db:"id"`               // Уникальный идентификатор типа колонки
	Name             string    `db:"name"`             // Название
	Description      string    `db:"description"`      // Описание
	Required         bool      `db:"required"`         // Обязательность к заполнению
	Regexp           string    `db:"regexp"`           // Регулярное выражение для проверки
	HorAlignmentHead Alignment `db:"horAlignmentHead"` // Горизонтальное выравнивание заголовка
	HorAlignmentBody Alignment `db:"horAlignmentBody"` // Горизонтальное выравнивание содержимого
	Created          time.Time `db:"created"`          // Время создания
	Active           bool      `db:"active"`           // Активная
}

// Конструктор создания объекта типа колонки в api
func NewApiColumnType(id int, name string, description string, required bool, regexp string,
	horalignmenthead string, horalignmentbody string) *ApiColumnType {
	return &ApiColumnType{
		ID:               id,
		Name:             name,
		Description:      description,
		Required:         required,
		Regexp:           regexp,
		HorAlignmentHead: horalignmenthead,
		HorAlignmentBody: horalignmentbody,
	}
}

// Конструктор создания объекта типа колонки в бд
func NewDtoColumnType(id int, name string, description string, required bool, regexp string,
	horalignmenthead Alignment, horalignmentbody Alignment, created time.Time, active bool) *DtoColumnType {
	return &DtoColumnType{
		ID:               id,
		Name:             name,
		Description:      description,
		Required:         required,
		Regexp:           regexp,
		HorAlignmentHead: horalignmenthead,
		HorAlignmentBody: horalignmentbody,
		Created:          created,
		Active:           active,
	}
}
