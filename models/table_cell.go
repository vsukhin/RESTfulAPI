package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

//Структура для организации хранения ячеек таблиц
type ViewTableCell struct {
	Value string `json:"data" validate:"max=255"` // Значение ячейки
}

type DtoTableCell struct {
	Table_Column_ID int64  `db:"table_column_id"` // Идентификатор колонки таблицы
	Value           string `db:"value"`           // Значение ячейки
	Checked         bool   `db:"checked"`         // Выполнялась проверка
	Valid           bool   `db:"valid"`           // Подходит под regexp
}

type ApiMetaTableCell struct {
	Checked        bool  `json:"verified" db:"checked"`    // Выполнялась проверка
	Valid          bool  `json:"correct" db:"valid"`       // Подходит под regexp
	Column_Type_ID int64 `json:"type" db:"column_type_id"` // Идентификатор типа колонки
}

type ApiShortTableCell struct {
	Value string `json:"data" db:"value"`    // Значение ячейки
	Valid bool   `json:"correct" db:"valid"` // Подходит под regexp
}

type ApiLongTableCell struct {
	Table_Column_ID int64  `json:"columnId" db:"table_column_id"` // Идентификатор колонки таблицы
	Value           string `json:"data" db:"value"`               // Значение ячейки
	Column_Type_ID  int64  `json:"type" db:"column_type_id"`      // Идентификатор типа колонки
	Valid           bool   `json:"correct" db:"valid"`            // Подходит под regexp
}

// Конструктор создания объекта ячейки таблицы в api
func NewApiMetaTableCell(checked bool, valid bool, column_type_id int64) *ApiMetaTableCell {
	return &ApiMetaTableCell{
		Checked:        checked,
		Valid:          valid,
		Column_Type_ID: column_type_id,
	}
}

func NewApiShortTableCell(value string, valid bool) *ApiShortTableCell {
	return &ApiShortTableCell{
		Value: value,
		Valid: valid,
	}
}

func NewApiLongTableCell(table_column_id int64, value string, column_type_id int64, valid bool) *ApiLongTableCell {
	return &ApiLongTableCell{
		Table_Column_ID: table_column_id,
		Value:           value,
		Column_Type_ID:  column_type_id,
		Valid:           valid,
	}
}

// Конструктор создания объекта ячейки таблицы в бд
func NewDtoTableCell(table_column_id int64, value string, checked bool, valid bool) *DtoTableCell {
	return &DtoTableCell{
		Table_Column_ID: table_column_id,
		Value:           value,
		Checked:         checked,
		Valid:           valid,
	}
}

func (cell *ViewTableCell) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(cell, errors, req)
}
