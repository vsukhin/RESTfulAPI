package models

import (
	"time"
)

const (
	ADDRESS_TYPE_LEGAL = 1
)

// Структура для организации хранения типа адреса
type ApiAddressType struct {
	ID       int    `json:"id" db:"id"`             // Уникальный идентификатор типа адреса
	Name     string `json:"name" db:"name"`         // Название
	Required bool   `json:"required" db:"required"` // Обязательность к заполнению
	Position int    `json:"position" db:"position"` // Позиция
}

type DtoAddressType struct {
	ID       int       `db:"id"`       // Уникальный идентификатор типа адреса
	Name     string    `db:"name"`     // Название
	Required bool      `db:"required"` // Обязательность к заполнению
	Position int       `db:"position"` // Позиция
	Created  time.Time `db:"created"`  // Время создания
	Active   bool      `db:"active"`   // Aктивен
}

// Конструктор создания объекта типа адреса в api
func NewApiAddressType(id int, name string, required bool, position int) *ApiAddressType {
	return &ApiAddressType{
		ID:       id,
		Name:     name,
		Required: required,
		Position: position,
	}
}

// Конструктор создания объекта типа адреса в бд
func NewDtoAddressType(id int, name string, required bool, position int, created time.Time, active bool) *DtoAddressType {
	return &DtoAddressType{
		ID:       id,
		Name:     name,
		Required: required,
		Position: position,
		Created:  created,
		Active:   active,
	}
}
