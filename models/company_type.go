package models

import (
	"time"
)

// Структура для организации хранения типа компании
type ApiCompanyType struct {
	ID            int    `json:"id" db:"id"`                      // Уникальный идентификатор типа компании
	FullName_rus  string `json:"nameFullRus" db:"fullname_rus"`   // Полное русское название
	FullName_eng  string `json:"nameFullEng" db:"fullname_eng"`   // Полное английское название
	ShortName_rus string `json:"nameShortRus" db:"shortname_rus"` // Короткое русское название
	ShortName_eng string `json:"nameShortEng" db:"shortname_eng"` // Короткое английское название
	Position      int    `json:"position" db:"position"`          // Позиция
}

type DtoCompanyType struct {
	ID            int       `db:"id"`            // Уникальный идентификатор типа компании
	FullName_rus  string    `db:"fullname_rus"`  // Полное русское название
	FullName_eng  string    `db:"fullname_eng"`  // Полное английское название
	ShortName_rus string    `db:"shortname_rus"` // Короткое русское название
	ShortName_eng string    `db:"shortname_eng"` // Короткое английское название
	Position      int       `db:"position"`      // Позиция
	Created       time.Time `db:"created"`       // Время создания
	Active        bool      `db:"active"`        // Aктивен
}

// Конструктор создания объекта типа компании в api
func NewApiCompanyType(id int, fullname_rus string, fullname_eng string, shortname_rus string, shortname_eng string,
	position int) *ApiCompanyType {
	return &ApiCompanyType{
		ID:            id,
		FullName_rus:  fullname_rus,
		FullName_eng:  fullname_eng,
		ShortName_rus: shortname_rus,
		ShortName_eng: shortname_eng,
		Position:      position,
	}
}

// Конструктор создания объекта типа компании в бд
func NewDtoCompanyType(id int, fullname_rus string, fullname_eng string, shortname_rus string, shortname_eng string,
	position int, created time.Time, active bool) *DtoCompanyType {
	return &DtoCompanyType{
		ID:            id,
		FullName_rus:  fullname_rus,
		FullName_eng:  fullname_eng,
		ShortName_rus: shortname_rus,
		ShortName_eng: shortname_eng,
		Position:      position,
		Created:       created,
		Active:        active,
	}
}
