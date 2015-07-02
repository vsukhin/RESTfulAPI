package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
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
	Email              string    `db:"email"`              // Уникальный email
	Confirmed          bool      `db:"confirmed"`          // Подтвержден
	Language           string    `db:"language"`           // Язык подписки
	Subscr_Created     time.Time `db:"subscr_created"`     // Время подписки
	Subscr_IP_Address  string    `db:"subscr_ip_address"`  // IP адрес подписки
	Subscr_Reverse_DNS string    `db:"subscr_reverse_dns"` // Обратный DNS подписки
	Subscr_User_Agent  string    `db:"subscr_user_agent"`  // User Agent подписки
	Conf_Created       time.Time `db:"conf_created"`       // Время подтверждения
	Conf_IP_Address    string    `db:"conf_ip_address"`    // IP адрес подтверждения
	Conf_Reverse_DNS   string    `db:"conf_reverse_dns"`   // Обратный DNS подтверждения
	Conf_User_Agent    string    `db:"conf_user_agent"`    // User Agent подтверждения
	Subscr_Code        string    `db:"subscr_code"`        // Код подтверждения подписки
	Unsubscr_Code      string    `db:"unsubscr_code"`      // Код подтверждения отписки
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
func NewDtoSubscription(email string, confirmed bool, language string, subscr_created time.Time, subscr_ip_address string,
	subscr_reverse_dns string, subscr_user_agent string, conf_created time.Time, conf_ip_address string, conf_reverse_dns string,
	conf_user_agent string, subscr_code string, unsubscr_code string) *DtoSubscription {
	return &DtoSubscription{
		Email:              email,
		Confirmed:          confirmed,
		Language:           language,
		Subscr_Created:     subscr_created,
		Subscr_IP_Address:  subscr_ip_address,
		Subscr_Reverse_DNS: subscr_reverse_dns,
		Subscr_User_Agent:  subscr_user_agent,
		Conf_Created:       conf_created,
		Conf_IP_Address:    conf_ip_address,
		Conf_Reverse_DNS:   conf_reverse_dns,
		Conf_User_Agent:    conf_user_agent,
		Subscr_Code:        subscr_code,
		Unsubscr_Code:      unsubscr_code,
	}
}

func (subscription *SubscriptionConfirm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(subscription, errors, req)
}

func (subscription *ViewSubscription) Validate(errors binding.Errors, req *http.Request) binding.Errors {

	return ValidateWithLanguage(subscription, errors, req, subscription.Language)
}
