package models

import (
	"time"
)

// Структура для организации хранения мобильного оператора
type ApiMobileOperator struct {
	ID        int    `json:"id" db:"id"`               // Уникальный идентификатор мобильного оператора
	ShortName string `json:"nameShort" db:"shortname"` // Короткое название
	LongName  string `json:"nameLong" db:"longname"`   // Длинное название
	Position  int    `json:"position" db:"position"`   // Позиция
}

type DtoMobileOperator struct {
	ID        int       `db:"id"`        // Уникальный идентификатор мобильного оператора
	ShortName string    `db:"shortname"` // Короткое название
	LongName  string    `db:"longname"`  // Длинное название
	Created   time.Time `db:"created"`   // Время создания
	Position  int       `db:"position"`  // Позиция
	Active    bool      `db:"active"`    // Aктивен
}

// Конструктор создания объекта мобильного оператора в api
func NewApiMobileOperator(id int, shortname string, longname string, position int) *ApiMobileOperator {
	return &ApiMobileOperator{
		ID:        id,
		ShortName: shortname,
		LongName:  longname,
		Position:  position,
	}
}

// Конструктор создания объекта мобильного оператора в бд
func NewDtoMobileOperator(id int, shortname string, longname string, created time.Time, position int, active bool) *DtoMobileOperator {
	return &DtoMobileOperator{
		ID:        id,
		ShortName: shortname,
		LongName:  longname,
		Created:   created,
		Position:  position,
		Active:    active,
	}
}
