package models

const (
	OPERATION_TYPE_RECEIVE  = 1
	OPERATION_TYPE_WITHDRAW = 2
)

// Структура для организации хранения типов операций
type DtoOperationType struct {
	ID          int    `db:"id"`          // Уникальный идентификатор типа операции
	Name        string `db:"name"`        // Название
	Description string `db:"description"` // Описание
}

// Конструктор создания объекта типа операции в бд
func NewDtoOperationType(id int, name string, description string) *DtoOperationType {
	return &DtoOperationType{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
