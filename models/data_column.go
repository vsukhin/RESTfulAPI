package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения колонки данных
type ViewDataColumn struct {
	Table_Column_ID int64 `json:"id" validate:"nonzero"` // Идентификатор колонки
}

type ApiDataColumn struct {
	Table_Column_ID int64  `json:"id" db:"table_column_id"`    // Идентификатор колонки
	Name            string `json:"name" db:"name"`             // Название
	Column_Type_ID  int    `json:"typeId" db:"column_type_id"` // Тип
	Position        int64  `json:"position" db:"position"`     // Позиция
}

type DtoDataColumn struct {
	Order_ID        int64 `db:"order_id"`        // Идентификатор заказа
	Table_Column_ID int64 `db:"table_column_id"` // Идентификатор колонки
}

// Конструктор создания объекта колонки данных в api
func NewApiDataColumn(table_column_id int64, name string, column_type_id int, position int64) *ApiDataColumn {
	return &ApiDataColumn{
		Table_Column_ID: table_column_id,
		Name:            name,
		Column_Type_ID:  column_type_id,
		Position:        position,
	}
}

// Конструктор создания объекта колонки данных в бд
func NewDtoDataColumn(order_id int64, table_column_id int64) *DtoDataColumn {
	return &DtoDataColumn{
		Order_ID:        order_id,
		Table_Column_ID: table_column_id,
	}
}

func (datacolumn ViewDataColumn) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(&datacolumn, errors, req)
}
