package models

import (
	"time"
)

const (
	DOCUMENT_TYPE_CHARTER              = 10
	DOCUMENT_TYPE_EXTRACTINCORPORATION = 11
	DOCUMENT_TYPE_MATCHING             = 2
)

// Структура для организации хранения типа документа
type ApiDocumentType struct {
	ID          int    `json:"id" db:"id"`                   // Уникальный идентификатор типа документа
	Position    int    `json:"position" db:"position"`       // Позиция
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type DtoDocumentType struct {
	ID          int       `db:"id"`          // Уникальный идентификатор типа документа
	Name        string    `db:"name"`        // Название
	Description string    `db:"description"` // Описание
	Position    int       `db:"position"`    // Позиция
	Created     time.Time `db:"created"`     // Время создания
	Active      bool      `db:"active"`      // Aктивен
}

// Конструктор создания объекта типа документа в api
func NewApiDocumentType(id int, position int, name string, description string) *ApiDocumentType {
	return &ApiDocumentType{
		ID:          id,
		Position:    position,
		Name:        name,
		Description: description,
	}
}

// Конструктор создания объекта типа доукмента в бд
func NewDtoDocumentType(id int, name string, description string, position int, created time.Time, active bool) *DtoDocumentType {
	return &DtoDocumentType{
		ID:          id,
		Name:        name,
		Description: description,
		Position:    position,
		Created:     created,
		Active:      active,
	}
}
