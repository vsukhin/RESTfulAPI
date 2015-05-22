package models

import (
	"time"
)

// Структура для организации хранения события
type ApiEvent struct {
	ID       int    `json:"id" db:"id"`             // Уникальный идентификатор события
	Name     string `json:"name" db:"name"`         // Название
	Position int    `json:"position" db:"position"` // Позиция
}

type DtoEvent struct {
	ID       int       `db:"id"`       // Уникальный идентификатор события
	Name     string    `db:"name"`     // Название
	Created  time.Time `db:"created"`  // Время создания
	Position int       `db:"position"` // Позиция
	Active   bool      `db:"active"`   // Aктивен
}

// Конструктор создания объекта события в api
func NewApiEvent(id int, name string, position int) *ApiEvent {
	return &ApiEvent{
		ID:       id,
		Name:     name,
		Position: position,
	}
}

// Конструктор создания объекта события в бд
func NewDtoEvent(id int, name string, created time.Time, position int, active bool) *DtoEvent {
	return &DtoEvent{
		ID:       id,
		Name:     name,
		Created:  created,
		Position: position,
		Active:   active,
	}
}
