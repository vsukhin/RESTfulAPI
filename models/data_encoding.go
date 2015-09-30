package models

// Структура для организации хранения кодировок данных
type ApiDataEncoding struct {
	ID   int    `json:"id" db:"id"`     // Уникальный идентификатор кодировки данных
	Name string `json:"name" db:"name"` // Название
}

type DtoDataEncoding struct {
	ID   int    `db:"id"`   // Уникальный идентификатор кодировки данных
	Name string `db:"name"` // Название
}

// Конструктор создания объекта кодировки данных в api
func NewApiDataEncoding(id int, name string) *ApiDataEncoding {
	return &ApiDataEncoding{
		ID:   id,
		Name: name,
	}
}

// Конструктор создания объекта кодировки данных в бд
func NewDtoDataEncoding(id int, name string) *DtoDataEncoding {
	return &DtoDataEncoding{
		ID:   id,
		Name: name,
	}
}
