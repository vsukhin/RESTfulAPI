package models

const (
	TRANSACTION_TYPE_REFILLING_ACCOUNT = 1
	TRANSACTION_TYPE_RETURNING_MONEY   = 2
	TRANSACTION_TYPE_SERVICE_FEE_MONTH = 3
	TRANSACTION_TYPE_SERVICE_FEE_SMS   = 4
	TRANSACTION_TYPE_SERVICE_FEE_HLR   = 5
)

// Структура для организации хранения типов транзакций
type DtoTransactionType struct {
	ID          int    `db:"id"`          // Уникальный идентификатор типа транзакции
	Name        string `db:"name"`        // Название
	Description string `db:"description"` // Описание
}

// Конструктор создания объекта типа транзакции в бд
func NewDtoTransactionType(id int, name string, description string) *DtoTransactionType {
	return &DtoTransactionType{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
