package models

// Структура для хранения счетa заказа
type DtoOrderInvoice struct {
	ID         int64 `db:"id"`         // Уникальный идентификатор счета заказа
	Order_ID   int64 `db:"order_id"`   // Идентификатор заказа
	Invoice_ID int64 `db:"invoice_id"` // Идентификатор счета
}

// Конструктор создания объекта счета заказа в бд
func NewDtoOrderInvoice(id int64, order_id int64, invoice_id int64) *DtoOrderInvoice {
	return &DtoOrderInvoice{
		ID:         id,
		Order_ID:   order_id,
		Invoice_ID: invoice_id,
	}
}
