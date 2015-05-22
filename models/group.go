package models

// Структура для хранения группы доступа
type ApiGroup struct {
	ID   int    `json:"id" db:"id"`     // Идентификатор группы доступа
	Name string `json:"name" db:"name"` // Название группы доступа
}

type DtoGroup struct {
	ID        int    `db:"id"`      // Идентификатор группы доступа
	Name      string `db:"name"`    // Название группы доступа
	IsDefault bool   `db:"default"` // Является ли группой доступа по умолчанию
	IsActive  bool   `db:"active"`  // Является ли активной
}

// Конструктор создания объекта группы в api
func NewApiGroup(id int, name string) *ApiGroup {
	return &ApiGroup{
		ID:   id,
		Name: name,
	}
}

// Конструктор создания объекта группы в бд
func NewDtoGroup(id int, name string, isdefault bool, isactive bool) *DtoGroup {
	return &DtoGroup{
		ID:        id,
		Name:      name,
		IsDefault: isdefault,
		IsActive:  isactive,
	}
}
