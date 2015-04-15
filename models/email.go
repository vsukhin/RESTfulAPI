package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

//Структура для организации хранения email
type UpdateEmail struct {
	Email         string `json:"email" validate:"nonzero,min=1,max=255,regexp=^.+@.+$"` // Уникальный email
	Primary       bool   `json:"primary"`                                               // Основной
	Confirmed     bool   `json:"-"`                                                     // Подтвержден
	Subscription  bool   `json:"subscription"`                                          // Используется для рассылки
	Language      string `json:"language" validate:"nonzero,min=1,max=10"`              // Язык рассылки
	Classifier_ID int    `json:"contactClass" validate:"nonzero"`                       // Идентификатор классификатора
}

type ApiEmail struct {
	Email         string `json:"email" db:"email"`                // Уникальный email
	Primary       bool   `json:"primary" db:"primary"`            // Основной
	Confirmed     bool   `json:"confirmed" db:"confirmed"`        // Подтвержден
	Subscription  bool   `json:"subscription" db:"subscription"`  // Используется для рассылки
	Language      string `json:"language" db:"language"`          // Язык рассылки
	Classifier_ID int    `json:"contactClass" db:"classifier_id"` // Идентификатор классификатора
}

type UpdateEmails []UpdateEmail

type DtoEmail struct {
	Email         string    `db:"email"`        // Уникальный email
	UserID        int64     `db:"user_id"`      // Идентификатор владельца email
	Classifier_ID int       `db:classifier_id"` // Идентификатор классификатора
	Created       time.Time `db:"created"`      // Время создания email
	Primary       bool      `db:"primary"`      // Основной
	Confirmed     bool      `db:"confirmed"`    // Подтвержден
	Subscription  bool      `db:"subscription"` // Используется для рассылки
	Code          string    `db:"code"`         // Код подтверждения
	Language      string    `db:"language"`     // Язык рассылки
	Exists        bool      `db:"-"`            // Существующий
}

// Конструктор создания объекта email в api
func NewApiEmail(email string, primary bool, confirmed bool, subscription bool, language string, classifier_id int) *ApiEmail {
	return &ApiEmail{
		Email:         email,
		Primary:       primary,
		Confirmed:     confirmed,
		Subscription:  subscription,
		Language:      language,
		Classifier_ID: classifier_id,
	}
}

// Конструктор создания объекта email в бд
func NewDtoEmail(email string, userid int64, classifier_id int, created time.Time, primary bool, confirmed bool,
	subscription bool, code string, language string, exists bool) *DtoEmail {
	return &DtoEmail{
		Email:         email,
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

func (email UpdateEmail) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return ValidateWithLanguage(&email, errors, req, email.Language)
}
