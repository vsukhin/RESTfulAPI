package models

import (
	"time"
)

// Структура для организации хранения операций
type DtoOperation struct {
	ID             int64     `db:"id"`             // Уникальный идентификатор операции
	Transaction_ID int64     `db:"transaction_id"` // Идентификатор транзакции
	Invoice_ID     int64     `db:"invoice_id"`     // Идентификатор счета
	Money          float64   `db:"money"`          // Деньги
	Type_ID        int       `db:"type_id"`        // Тип операции
	Created        time.Time `db:"created"`        // Время создания
}

// Конструктор создания объекта операции в бд
func NewDtoOperation(id int64, transaction_id int64, invoice_id int64,
	money float64, type_id int, created time.Time) *DtoOperation {
	return &DtoOperation{
		ID:             id,
		Transaction_ID: transaction_id,
		Invoice_ID:     invoice_id,
		Money:          money,
		Type_ID:        type_id,
		Created:        created,
	}
}
