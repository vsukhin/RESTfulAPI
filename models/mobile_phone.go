package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

//Структура для организации хранения мобильных телефонов
type ViewApiMobilePhone struct {
	Phone         string `json:"phone" db:"phone" validate:"nonzero,min=1,max=25,regexp=^[0-9]*$"` // Уникальный номер
	Primary       bool   `json:"primary" db:"primary"`                                             // Основной
	Confirmed     bool   `json:"confirmed" db:"confirmed"`                                         // Подтвержден
	Subscription  bool   `json:"subscription" db:"subscription"`                                   // Используется для рассылки
	Language      string `json:"language" db:"language" validate:"nonzero,min=1,max=10"`           // Язык рассылки
	Classifier_ID int    `json:"contactClass" db:"classifier_id" validate:"nonzero"`               // Идентификатор классификатора
}

type UpdateMobilePhones []ViewApiMobilePhone

type DtoMobilePhone struct {
	Phone         string    `db:"phone"`        // Уникальный номер
	UserID        int64     `db:"user_id"`      // Идентификатор владельца MobilePhone
	Classifier_ID int       `db:classifier_id"` // Идентификатор классификатора
	Created       time.Time `db:"created"`      // Время создания MobilePhone
	Primary       bool      `db:"primary"`      // Основной
	Confirmed     bool      `db:"confirmed"`    // Подтвержден
	Subscription  bool      `db:"subscription"` // Используется для рассылки
	Code          string    `db:"code"`         // Код подтверждения
	Language      string    `db:"language"`     // Язык рассылки
	Exists        bool      `db:"-"`            // Существующий
}

// Конструктор создания объекта мобильного телефона в api
func NewViewApiMobilePhone(phone string, primary bool, confirmed bool, subscription bool,
	language string, classifier_id int) *ViewApiMobilePhone {
	return &ViewApiMobilePhone{
		Phone:         phone,
		Primary:       primary,
		Confirmed:     confirmed,
		Subscription:  subscription,
		Language:      language,
		Classifier_ID: classifier_id,
	}
}

// Конструктор создания объекта мобильного телефона в бд
func NewDtoMobilePhone(phone string, userid int64, classifier_id int, created time.Time, primary bool, confirmed bool,
	subscription bool, code string, language string, exists bool) *DtoMobilePhone {
	return &DtoMobilePhone{
		Phone:         phone,
		UserID:        userid,
		Classifier_ID: classifier_id,
		Created:       created,
		Primary:       primary,
		Confirmed:     confirmed,
		Subscription:  subscription,
		Code:          code,
		Language:      language,
		Exists:        exists,
	}
}

func (phone ViewApiMobilePhone) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return ValidateWithLanguage(&phone, errors, req, phone.Language)
}
