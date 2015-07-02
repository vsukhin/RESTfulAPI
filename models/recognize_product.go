package models

import (
	"time"
)

// Структура для организации хранения позиции ввода данных
type ApiRecognizeProduct struct {
	ID          int    `json:"id" db:"id"`                   // Уникальный идентификатор позиции ввода данных
	Position    int    `json:"position" db:"position"`       // Позиция
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
	Increase    bool   `json:"orderIncrease" db:"increase`   // Увеличение
}

type DtoRecognizeProduct struct {
	ID          int       `db:"id"`          // Уникальный идентификатор позиции ввода данных
	Position    int       `db:"position"`    // Позиция
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Increase    bool      `db:"increase`     // Увеличение
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта позиции ввода данных в api
func NewApiRecognizeProduct(id int, position int, name string, description string, increase bool) *ApiRecognizeProduct {
	return &ApiRecognizeProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
		Increase:    increase,
	}
}

// Конструктор создания объекта позиции ввода данных в бд
func NewDtoRecognizeProduct(id int, position int, name string, description string, increase bool, created time.Time, active bool) *DtoRecognizeProduct {
	return &DtoRecognizeProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
		Increase:    increase,
		Created:     created,
		Active:      active,
	}
}
