package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для организации хранения устройства
type ViewLongDevice struct {
	OS     string `json:"nameOs" validate:"min=1,max=255"`       // Операционная система
	App    string `json:"nameApp" validate:"min=1,max=255"`      // Приложение
	Serial string `json:"serialNumber" validate:"min=1,max=255"` // Серийный номер
}

type ViewHashDevice struct {
	Hash string `json:"hash" validate:"min=1,max=255"` // Хэш
}

type ViewCodeDevice struct {
	Code string `json:"code" validate:"min=1,max=255"` // Код
}

type ViewTokenDevice struct {
	Token string `json:"token" validate:"min=1,max=255"` // Токен
}

type ApiDevice struct {
	Token string `json:"token" db:"token"` // Постоянный токен
	Code  string `json:"code" db:"code"`   // Kod
}

type DtoDevice struct {
	ID         int64     `db:"id"`         // Уникальный идентификатор устройства
	User_ID    int64     `db:"user_id"`    // Идентификатор пользователя
	OS         string    `db:"os"`         // Операционная система
	App        string    `db:"app"`        // Приложение
	Serial     string    `db:"serial"`     // Серийный номер
	Token      string    `db:"token"`      // Постоянный токен
	Code       string    `db:"code"`       // Код
	Hash       string    `db:"hash"`       // Хэш
	Valid_Till time.Time `db:"valid_till"` // Действует до
	Created    time.Time `db:"created"`    // Время создания
	Active     bool      `db:"active"`     // Aктивен
}

// Конструктор создания объекта устройства в api
func NewApiDevice(token string, code string) *ApiDevice {
	return &ApiDevice{
		Token: token,
		Code:  code,
	}
}

// Конструктор создания объекта устройства в бд
func NewDtoDevice(id int64, user_id int64, os string, app string, serial string, token string, code string,
	hash string, valid_till time.Time, created time.Time, active bool) *DtoDevice {
	return &DtoDevice{
		ID:         id,
		User_ID:    user_id,
		OS:         os,
		App:        app,
		Serial:     serial,
		Token:      token,
		Code:       code,
		Hash:       hash,
		Valid_Till: valid_till,
		Created:    created,
		Active:     active,
	}
}

func (device *ViewLongDevice) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(device, errors, req)
}

func (device *ViewHashDevice) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(device, errors, req)
}

func (device *ViewCodeDevice) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(device, errors, req)
}

func (device *ViewTokenDevice) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(device, errors, req)
}
