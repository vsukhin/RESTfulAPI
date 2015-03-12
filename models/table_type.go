package models

const (
	TABLE_TYPE_DEFAULT = "none"
	TABLE_TYPE_PRICE   = "price"
)

//Структура для организации хранения типов таблиц
type DtoTableType struct {
	ID   int64  `db:"id"`   // Уникальный идентификатор типа таблицы
	Name string `db:"name"` // Назавание типа таблицы
}

// Конструктор создания объекта типа таблицы в бд
func NewDtoTableType(id int64, name string) *DtoTableType {
	return &DtoTableType{
		ID:   id,
		Name: name,
	}
}
