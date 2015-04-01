package models

import (
	"time"
)

//Структура для организации хранения сервисов
type ViewFacility struct {
	ID int64 `json:"id"` // Уникальный идентификатор сервиса
}

type ViewFacilities []ViewFacility

type ApiShortFacility struct {
	ID          int64  `json:"id" db:"id"`                   // Уникальный идентификатор сервиса
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type ApiLongFacility struct {
	ID          int64  `json:"id" db:"id"`                   // Уникальный идентификатор сервиса
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
	Active      bool   `json:"active" db:"active"`           // Активный
}

type DtoFacility struct {
	ID          int64     `db:"id"`          // Уникальный идентификатор сервиса
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Активный
}

// Конструктор создания объекта сервиса в api
func NewApiShortFacility(id int64, name string, description string) *ApiShortFacility {
	return &ApiShortFacility{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

func NewApiLongFacility(id int64, name string, description string, active bool) *ApiLongFacility {
	return &ApiLongFacility{
		ID:          id,
		Name:        name,
		Description: description,
		Active:      active,
	}
}

// Конструктор создания объекта сервиса в бд
func NewDtoFacility(id int64, name string, description string, created time.Time, active bool) *DtoFacility {
	return &DtoFacility{
		ID:          id,
		Name:        name,
		Description: description,
		Created:     created,
		Active:      active,
	}
}
