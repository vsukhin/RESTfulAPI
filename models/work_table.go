package models

import (
	"time"
)

// Структура для хранения рабочей таблицы
type ApiWorkTable struct {
	Customer_Table_ID int64     `json:"tableId" db:"customer_table_id"` // Идентификатор таблицы
	Created           time.Time `json:"created" db:"created"`           // Дата и время создания
	TypeID            int       `json:"type" db:"type_id"`              // Тип
}

type DtoWorkTable struct {
	Order_ID          int64 `db:"order_id"`          // Идентификатор заказа
	Customer_Table_ID int64 `db:"customer_table_id"` // Идентификатор таблицы
}

// Конструктор создания объекта рабочей таблицы в api
func NewApiWorkTable(customer_table_id int64, created time.Time, typeid int) *ApiWorkTable {
	return &ApiWorkTable{
		Customer_Table_ID: customer_table_id,
		Created:           created,
		TypeID:            typeid,
	}
}

// Конструктор создания объекта рабочей таблицы в бд
func NewDtoWorkTable(order_id int64, customer_table_id int64) *DtoWorkTable {
	return &DtoWorkTable{
		Order_ID:          order_id,
		Customer_Table_ID: customer_table_id,
	}
}
