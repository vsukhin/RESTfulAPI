package models

import (
	"time"
)

// Структура для организации хранения позиции верификации данных
type ApiVerifyProduct struct {
	ID          int    `json:"id" db:"id"`                   // Уникальный идентификатор позиции верификации данных
	Position    int    `json:"position" db:"position"`       // Позиция
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type DtoVerifyProduct struct {
	ID          int       `db:"id"`          // Уникальный идентификатор позиции верификации данных
	Position    int       `db:"position"`    // Позиция
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта позиции верификации данных в api
func NewApiVerifyProduct(id int, position int, name string, description string) *ApiVerifyProduct {
	return &ApiVerifyProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
	}
}

// Конструктор создания объекта верификации данных в бд
func NewDtoVerifyProduct(id int, position int, name string, description string, created time.Time, active bool) *DtoVerifyProduct {
	return &DtoVerifyProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
		Created:     created,
		Active:      active,
	}
}
