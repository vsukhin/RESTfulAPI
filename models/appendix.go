package models

import (
	"time"
)

// Структура для организации хранения приложения договора
type ApiAppendix struct {
	SignedDate time.Time `json:"createdDate" db:"signed_date"` // Дата подписания
	Name       string    `json:"name" db:"name"`               // Название
	File_ID    int64     `json:"fileId" db:"file_id"`          // Идентификатор файла
}

type DtoAppendix struct {
	ID          int64     `db:"id"`          // Уникальный идентификатор приложения договора
	Contract_ID int64     `db:"contract_id"` // Идентификатор контракта
	Name        string    `db:"name"`        // Название
	File_ID     int64     `db:"file_id"`     // Идентификатор файла
	SignedDate  time.Time `db:"signed_date"` // Дата подписания
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта приложения договора в api
func NewApiAppendix(signeddate time.Time, name string, file_id int64) *ApiAppendix {
	return &ApiAppendix{
		SignedDate: signeddate,
		Name:       name,
		File_ID:    file_id,
	}
}

// Конструктор создания объекта приложения договора в бд
func NewDtoAppendix(id int64, contract_id int64, name string, file_id int64, signeddate time.Time,
	position int, created time.Time, active bool) *DtoAppendix {
	return &DtoAppendix{
		ID:          id,
		Contract_ID: contract_id,
		Name:        name,
		File_ID:     file_id,
		SignedDate:  signeddate,
		Created:     created,
		Active:      active,
	}
}
