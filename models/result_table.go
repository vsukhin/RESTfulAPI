package models

import (
	"time"
)

// Структура для хранения таблицы результатов
type ApiResultTable struct {
	Customer_Table_ID int64     `json:"tableId" db:"customer_table_id"` // Идентификатор таблицы
	Created           time.Time `json:"created" db:"created"`           // Дата и время создания
	TypeID            int       `json:"type" db:"type_id"`              // Тип
}

type DtoResultTable struct {
	Order_ID          int64 `db:"order_id"`          // Идентификатор заказа
	Customer_Table_ID int64 `db:"customer_table_id"` // Идентификатор таблицы
}

// Конструктор создания объекта таблицы результатов в api
func NewApiResultTable(customer_table_id int64, created time.Time, typeid int) *ApiResultTable {
	return &ApiResultTable{
		Customer_Table_ID: customer_table_id,
		Created:           created,
		TypeID:            typeid,
	}
}

// Конструктор создания объекта таблицы результатов в бд
func NewDtoResultTable(order_id int64, customer_table_id int64) *DtoResultTable {
	return &DtoResultTable{
		Order_ID:          order_id,
		Customer_Table_ID: customer_table_id,
	}
}
