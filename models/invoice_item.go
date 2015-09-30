package models

const (
	INVOICE_ITEM_TYPE_ROUBLE  = "руб"
	INVOICE_ITEM_NAME_DEFAULT = "Оплата по договору"
)

// Структура для организации позиции счета
type ApiInvoiceItem struct {
	ID      int64   `json:"id" db:"id"`             // Уникальный идентификатор позиции счета
	Name    string  `json:"name" db:"name"`         // Название
	Measure string  `json:"measure" db:"measure"`   // Единица измерения
	Amount  float64 `json:"amount" db:"amount"`     // Количество
	Price   float64 `json:"priceAmount" db:"price"` // Цена
	Total   float64 `json:"priceTotal" db:"total"`  // Стоимость
}

type DtoInvoiceItem struct {
	ID         int64   `db:"id"`         // Уникальный идентификатор позиции счета
	Invoice_ID int64   `db:"invoice_id"` // Идентификатор счета
	Name       string  `db:"name"`       // Название
	Measure    string  `db:"measure"`    // Единица измерения
	Amount     float64 `db:"amount"`     // Количество
	Price      float64 `db:"price"`      // Цена
	Total      float64 `db:"total"`      // Стоимость

}

// Конструктор создания объекта позиции счета в api
func NewApiInvoiceItem(id int64, name string, measure string, amount float64, price float64, total float64) *ApiInvoiceItem {
	return &ApiInvoiceItem{
		ID:      id,
		Name:    name,
		Measure: measure,
		Amount:  amount,
		Price:   price,
		Total:   total,
	}
}

// Конструктор создания объекта позиции счета в бд
func NewDtoInvoiceItem(id int64, invoice_id int64, name string, measure string, amount float64, price float64, total float64) *DtoInvoiceItem {
	return &DtoInvoiceItem{
		ID:         id,
		Invoice_ID: invoice_id,
		Name:       name,
		Measure:    measure,
		Amount:     amount,
		Price:      price,
		Total:      total,
	}
}

func (item *ApiInvoiceItem) GetNumber(index int) (number int) {
	return index + 1
}
