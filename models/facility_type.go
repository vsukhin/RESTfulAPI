package models

import (
	"time"
)

// Структура для организации хранения типов сервиса
type ApiFacilityType struct {
	ID   int    `json:"id" db:"id"`     // Уникальный идентификатор типа сервиса
	Name string `json:"name" db:"name"` // Название
}

type DtoFacilityType struct {
	ID      int       `db:"id"`      // Уникальный идентификатор типа сервиса
	Name    string    `db:"name"`    // Название
	Active  bool      `db:"active"`  // Aктивен
	Created time.Time `db:"created"` // Время создания
}

// Конструктор создания объекта типа сервиса в api
func NewApiFacilityType(id int, name string) *ApiFacilityType {
	return &ApiFacilityType{
		ID:   id,
		Name: name,
	}
}

// Конструктор создания объекта типа сервиса в бд
func NewDtoFacilityType(id int, name string, active bool, created time.Time) *DtoFacilityType {
	return &DtoFacilityType{
		ID:      id,
		Name:    name,
		Active:  active,
		Created: created,
	}
}
