package models

import (
	"time"
)

// Структура для организации хранения периода
type ApiPeriod struct {
	ID       int    `json:"id" db:"id"`             // Уникальный идентификатор периода
	Name     string `json:"name" db:"name"`         // Название
	Position int    `json:"position" db:"position"` // Позиция
}

type DtoPeriod struct {
	ID       int       `db:"id"`       // Уникальный идентификатор периода
	Name     string    `db:"name"`     // Название
	Created  time.Time `db:"created"`  // Время создания
	Position int       `db:"position"` // Позиция
	Active   bool      `db:"active"`   // Aктивен
}

// Конструктор создания объекта периода в api
func NewApiPeriod(id int, name string, position int) *ApiPeriod {
	return &ApiPeriod{
		ID:       id,
		Name:     name,
		Position: position,
	}
}

// Конструктор создания объекта периода в бд
func NewDtoPeriod(id int, name string, created time.Time, position int, active bool) *DtoPeriod {
	return &DtoPeriod{
		ID:       id,
		Name:     name,
		Created:  created,
		Position: position,
		Active:   active,
	}
}
