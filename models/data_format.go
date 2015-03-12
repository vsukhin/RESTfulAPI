package models

//Структура для организации хранения форматов данных
type ApiDataFormat struct {
	ID          int    `json:"id" db:"id"`                   // Уникальный идентификатор формата данных
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type DtoDataFormat struct {
	ID          int    `db:"id"`          // Уникальный идентификатор формата данных
	Name        string `db:"name"`        // Название
	Description string `db:"description"` // Описание
}

// Конструктор создания объекта формата данных в api
func NewApiDataFormat(id int, name string, description string) *ApiDataFormat {
	return &ApiDataFormat{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

// Конструктор создания объекта формата данных в бд
func NewDtoDataFormat(id int, name string, description string) *DtoDataFormat {
	return &DtoDataFormat{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
