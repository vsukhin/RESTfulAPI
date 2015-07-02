package models

// Структура для хранения позиции справочника прайс листа
type ApiProduct struct {
	Product_ID      int    // Идентификатор позиции
	Name            string // Название
	Table_Row_ID    int64  // Идентификатор строки
	Table_Column_ID int64  // Идентификатор колонки
}

// Конструктор создания объекта позиции справочника прайс листа в api
func NewApiProduct(product_id int, name string, table_row_id int64, table_column_id int64) *ApiProduct {
	return &ApiProduct{
		Product_ID:      product_id,
		Name:            name,
		Table_Row_ID:    table_row_id,
		Table_Column_ID: table_column_id,
	}
}
