package models

import (
	"time"
)

// Структура для организации хранения тарифного плана
type ApiTariffPlan struct {
	ID       int    `json:"id" db:"id"`             // Уникальный идентификатор тарифного плана
	Name     string `json:"name" db:"name"`         // Название
	Position int    `json:"position" db:"position"` // Позиция
	Public   bool   `json:"public" db:"public"`     // Публичность
}

type DtoTariffPlan struct {
	ID       int       `db:"id"`       // Уникальный идентификатор тарифного плана
	Name     string    `db:"name"`     // Название
	Position int       `db:"position"` // Позиция
	Public   bool      `db:"public"`   // Публичность
	Created  time.Time `db:"created"`  // Время создания
	Active   bool      `db:"active"`   // Aктивен
}

// Конструктор создания объекта тарифного плана в api
func NewApiTariffPlan(id int, name string, position int, public bool) *ApiTariffPlan {
	return &ApiTariffPlan{
		ID:       id,
		Name:     name,
		Position: position,
		Public:   public,
	}
}

// Конструктор создания объекта тарифного плана в бд
func NewDtoTariffPlan(id int, name string, position int, public bool, created time.Time, active bool) *DtoTariffPlan {
	return &DtoTariffPlan{
		ID:       id,
		Name:     name,
		Position: position,
		Public:   public,
		Created:  created,
		Active:   active,
	}
}
