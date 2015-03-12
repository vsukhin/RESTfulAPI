package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

//Структура для организации хранения колонок таблиц
type ViewApiTableColumn struct {
	Name     string `json:"name" validate:"nonzero,min=1,max=255"` // Название колонки таблицы
	TypeID   int64  `json:"typeId"`                                // Идентификатор типа
	Position int64  `json:"position"`                              // Позиция
}

type ViewApiOrderTableColumn struct {
	ID       int64 `json:"id" db:"id" validate:"nonzero"` // Уникальный идентификатор колонки таблицы
	Position int64 `json:"position" db:"position"`        // Позиция
}

type ApiTableColumn struct {
	ID       int64  `json:"id" db:"id"`                 // Уникальный идентификатор временной колонки таблицы
	Name     string `json:"name" db:"name"`             // Название
	TypeID   int64  `json:"typeId" db:"column_type_id"` // Идентификатор типа
	Position int64  `json:"position" db:"position"`     // Позиция
}

type ViewApiOrderTableColumns []ViewApiOrderTableColumn

type DtoTableColumn struct {
	ID                int64     `db:"id"`                // Уникальный идентификатор колонки таблицы
	Name              string    `db:"name"`              // Название
	Column_Type_ID    int64     `db:"column_type_id"`    // Идентификатор типа колонки
	Customer_Table_ID int64     `db:"customer_table_id"` // Идентификатор пользовательской таблицы
	Position          int64     `db:"position"`          // Позиция
	Created           time.Time `db:"created"`           // Время создания
	Prebuilt          bool      `db:"prebuilt"`          // Предопределенная
	Active            bool      `db:"active"`            // Активный
	Edition           int64     `db:"edition"`           // Версия редакции
	Original_ID       int64     `db:"original_id"`       // Оригинальная колонка
}

// Конструктор создания объекта коклонки таблицы в api
func NewViewApiTableColumn(name string, typeid int64, position int64) *ViewApiTableColumn {
	return &ViewApiTableColumn{
		Name:     name,
		TypeID:   typeid,
		Position: position,
	}
}

func NewViewApiOrderTableColumn(id int64, position int64) *ViewApiOrderTableColumn {
	return &ViewApiOrderTableColumn{
		ID:       id,
		Position: position,
	}
}

func NewApiTableColumn(id int64, name string, typeid int64, position int64) *ApiTableColumn {
	return &ApiTableColumn{
		ID:       id,
		Name:     name,
		TypeID:   typeid,
		Position: position,
	}
}

// Конструктор создания объекта колонки таблицы в бд
func NewDtoTableColumn(id int64, name string, column_type_id int64, customer_table_id int64, position int64,
	created time.Time, prebuilt bool, active bool, edition int64, original_id int64) *DtoTableColumn {
	return &DtoTableColumn{
		ID:                id,
		Name:              name,
		Column_Type_ID:    column_type_id,
		Customer_Table_ID: customer_table_id,
		Position:          position,
		Created:           created,
		Prebuilt:          prebuilt,
		Active:            active,
		Edition:           edition,
		Original_ID:       original_id,
	}
}

func (tablecolumn *ViewApiTableColumn) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(tablecolumn, errors, req)
}

func (tablecolumn ViewApiOrderTableColumn) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(&tablecolumn, errors, req)
}

func (tablecolumns ViewApiOrderTableColumns) GetIDs() []int64 {
	ids := new([]int64)
	for _, tablecolumn := range tablecolumns {
		*ids = append(*ids, tablecolumn.ID)
	}

	return *ids
}
