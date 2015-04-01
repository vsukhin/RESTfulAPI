package models

import (
	"time"
)

//Структура для организации хранения шага импорта
type ApiImportStep struct {
	Step       byte `json:"step" `   // Номер шага
	Ready      bool `json:"ready" `  // Готова
	Percentage byte `json:"percent"` // Процент готовности
}

type DtoImportStep struct {
	Customer_Table_ID int64     `db:"customer_table_id" ` // Уникальный идентификатор пользовательской таблицы
	Step              byte      `db:"step" `              // Номер шага
	Ready             bool      `db:"ready" `             // Готова
	Percentage        byte      `db:"percentage"`         // Процент готовности
	Started           time.Time `db:"started"`            // Время начала выполнения
	Completed         time.Time `db:"completed"`          // Время окончания выполнения
}

// Конструктор создания объекта шага импорта в api
func NewApiImportStep(step byte, ready bool, percentage byte) *ApiImportStep {
	return &ApiImportStep{
		Step:       step,
		Ready:      ready,
		Percentage: percentage,
	}
}

// Конструктор создания объекта шага импорта в бд
func NewDtoImportStep(customer_table_id int64, step byte, ready bool, percentage byte, started time.Time, completed time.Time) *DtoImportStep {
	return &DtoImportStep{
		Customer_Table_ID: customer_table_id,
		Step:              step,
		Ready:             ready,
		Percentage:        percentage,
		Started:           started,
		Completed:         completed,
	}
}
