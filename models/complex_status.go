package models

import (
	"time"
)

// Структура для организации хранения комлексного статуса
type ApiComplexStatus struct {
	ID          int    `json:"id" db:"id"`                   // Уникальный идентификатор комлексного статуса
	Final       bool   `json:"final" db:"final`              // Финальный
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type DtoComplexStatus struct {
	ID          int       `db:"id"`          // Уникальный идентификатор комлексного статуса
	Final       bool      `db:"final`        // Финальный
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта комлексного статуса в api
func NewApiComplexStatus(id int, final bool, name string, description string) *ApiComplexStatus {
	return &ApiComplexStatus{
		ID:          id,
		Final:       final,
		Name:        name,
		Description: description,
	}
}

// Конструктор создания объекта комлексного статуса в бд
func NewDtoComplexStatus(id int, final bool, name string, description string, created time.Time, active bool) *DtoComplexStatus {
	return &DtoComplexStatus{
		ID:          id,
		Final:       final,
		Name:        name,
		Description: description,
		Created:     created,
		Active:      active,
	}
}
