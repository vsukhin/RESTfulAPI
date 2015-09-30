package models

const (
	TABLE_TYPE_DEFAULT         = 0
	TABLE_TYPE_READONLY        = 1
	TABLE_TYPE_HIDDEN          = 2
	TABLE_TYPE_HIDDEN_READONLY = 3
	TABLE_TYPE_PRICE           = 4
)

// Структура для организации хранения типов таблиц
type ApiTableType struct {
	ID   int    `json:"id" db:"id"`     // Уникальный идентификатор типа таблицы
	Name string `json:"name" db:"name"` // Назавание типа таблицы
}

type DtoTableType struct {
	ID   int    `db:"id"`   // Уникальный идентификатор типа таблицы
	Name string `db:"name"` // Назавание типа таблицы
}

// Конструктор создания объекта типа таблицы в api
func NewApiTableType(id int, name string) *ApiTableType {
	return &ApiTableType{
		ID:   id,
		Name: name,
	}
}

// Конструктор создания объекта типа таблицы в бд
func NewDtoTableType(id int, name string) *DtoTableType {
	return &DtoTableType{
		ID:   id,
		Name: name,
	}
}
