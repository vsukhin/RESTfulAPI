package models

// Структура для организации хранения транзакций
type DtoTransaction struct {
	ID             int64 `db:"id"`             // Уникальный идентификатор транзакции
	Source_ID      int64 `db:"source_id"`      // С какого аккаунта
	Destination_ID int64 `db:"destination_id"` // На какой аккаунт
	Type_ID        int   `db:"type_id"`        // Тип транзакции
}

// Конструктор создания объекта транзакции в бд
func NewDtoTransaction(id int64, source_id int64, destination_id int64, type_id int) *DtoTransaction {
	return &DtoTransaction{
		ID:             id,
		Source_ID:      source_id,
		Destination_ID: destination_id,
		Type_ID:        type_id,
	}
}
