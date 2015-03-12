package models

import (
	"time"
)

//Структура для организации хранения объединений
type DtoUnit struct {
	ID      int64     `db:"id"`      // Уникальный идентификатор объединения
	Created time.Time `db:"created"` // Время создания объединения
	Name    string    `db:"name"`    // Название объединения
}

// Конструктор создания объекта объединения в бд
func NewDtoUnit(id int64, created time.Time, name string) *DtoUnit {
	return &DtoUnit{
		ID:      id,
		Created: created,
		Name:    name,
	}
}
