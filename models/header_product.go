package models

import (
	"time"
)

// Структура для организации хранения позиции имени отправителя
type ApiHeaderProduct struct {
	ID          int    `json:"id" db:"id"`                         // Уникальный идентификатор позиции имени отправителя
	Position    int    `json:"position" db:"position"`             // Позиция
	Name        string `json:"name" db:"name"`                     // Название
	Description string `json:"description" db:"description"`       // Описание
	Increase    bool   `json:"orderIncrease" db:"increase"`        // Увеличение
	FeeOnce     bool   `json:"subscription" db:"fee_once"`         // Разовая оплата
	FeeMonthly  bool   `json:"subscriptionMonth" db:"fee_monthly"` // Ежемесячная оплата
}

type DtoHeaderProduct struct {
	ID          int       `db:"id"`          // Уникальный идентификатор позиции имени отправителя
	Position    int       `db:"position"`    // Позиция
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Increase    bool      `db:"increase"`    // Увеличение
	FeeOnce     bool      `db:"fee_once"`    // Разовая оплата
	FeeMonthly  bool      `db:"fee_monthly"` // Ежемесячная оплата
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта позиции имени отправителя в api
func NewApiHeaderProduct(id int, position int, name string, description string, increase, feeonce, feemonthly bool) *ApiHeaderProduct {
	return &ApiHeaderProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
		Increase:    increase,
		FeeOnce:     feeonce,
		FeeMonthly:  feemonthly,
	}
}

// Конструктор создания объекта имени отправителя в бд
func NewDtoHeaderProduct(id int, position int, name string, description string, increase, feeonce, feemonthly bool,
	created time.Time, active bool) *DtoHeaderProduct {
	return &DtoHeaderProduct{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
		Increase:    increase,
		FeeOnce:     feeonce,
		FeeMonthly:  feemonthly,
		Created:     created,
		Active:      active,
	}
}
