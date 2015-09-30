package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения подписки
type ViewSubscription struct {
	Email        string `json:"email" validate:"max=255,regexp=^.+@.+$"` // Email подписчика
	CaptchaValue string `json:"captchaValue" validate:"max=255"`         // Значение капчи
	CaptchaHash  string `json:"captchaHash" validate:"max=255"`          // Хэш капчи
	Language     string `json:"language" validate:"max=10"`              // Язык подписчика
}

type SubscriptionConfirm struct {
	Code string `json:"code" validate:"min=1,max=255"` // Код подтверждения подписки
}

type ApiShortSubscription struct {
	Email string `json:"emailTo"` // Уникальный email
}

type ApiMiddleSubscription struct {
	FromEmail string `json:"emailFrom"` // Email рассылки
	ToEmail   string `json:"emailTo"`   // Email подписчика
}

type ApiLongSubscription struct {
	Email     string `json:"email"`      // Уникальный email
	Confirmed bool   `json:"subscribe" ` // Подтвержденный email
	Valid     bool   `json:"correct" `   // Существующий email
}

type DtoSubscription struct {
	Email               string `db:"email"`               // Уникальный email
	Confirmed           bool   `db:"confirmed"`           // Подтвержден
	Language            string `db:"language"`            // Язык подписки
	Subscr_AccessLog_ID int64  `db:"subscr_accesslog_id"` // Идентификатор лога подписки
	Conf_AccessLog_ID   int64  `db:"conf_accesslog_id"`   // Идентификатор лога подтверждения
	Subscr_Code         string `db:"subscr_code"`         // Код подтверждения подписки
	Unsubscr_Code       string `db:"unsubscr_code"`       // Код подтверждения отписки
}

// Конструктор создания объекта подписки в api
func NewApiShortSubscription(email string) *ApiShortSubscription {
	return &ApiShortSubscription{
		Email: email,
	}
}

func NewApiMiddleSubscription(fromemail string, toemail string) *ApiMiddleSubscription {
	return &ApiMiddleSubscription{
		FromEmail: fromemail,
		ToEmail:   toemail,
	}
}

func NewApiLongSubscription(email string, confirmed bool, valid bool) *ApiLongSubscription {
	return &ApiLongSubscription{
		Email:     email,
		Confirmed: confirmed,
		Valid:     valid,
	}
}

// Конструктор создания объекта подписки в бд
func NewDtoSubscription(email string, confirmed bool, language string, subscr_accesslog_id int64, conf_accesslog_id int64,
	subscr_code string, unsubscr_code string) *DtoSubscription {
	return &DtoSubscription{
		Email:               email,
		Confirmed:           confirmed,
		Language:            language,
		Subscr_AccessLog_ID: subscr_accesslog_id,
		Conf_AccessLog_ID:   conf_accesslog_id,
		Subscr_Code:         subscr_code,
		Unsubscr_Code:       unsubscr_code,
	}
}

func (subscription *SubscriptionConfirm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(subscription, errors, req)
}

func (subscription *ViewSubscription) Validate(errors binding.Errors, req *http.Request) binding.Errors {

	return ValidateWithLanguage(subscription, errors, req, subscription.Language)
}
