package models

//Структура для организации хранения статусов
type DtoStatus struct {
	ID          int    `db:"id"`          // Уникальный идентификатор статуса
	Name        string `db:"name"`        // Название
	Description string `db:"description"` // Описание
}

// Конструктор создания объекта статуса в бд
func NewDtoStatus(id int, name string, description string) *DtoStatus {
	return &DtoStatus{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
