package models

import (
	"time"
)

// Структура для организации хранения виртуальных директорий
type DtoVirtualDir struct {
	Token   string    `db:"token"`   // Уникальный токен для доступа
	Created time.Time `db:"created"` // Время создания
}

// Конструктор создания объекта виртуальной директории в бд
func NewDtoVirtualDir(token string, created time.Time) *DtoVirtualDir {
	return &DtoVirtualDir{
		Token:   token,
		Created: created,
	}
}
