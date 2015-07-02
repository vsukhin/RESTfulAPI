package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для организации хранения обращения
type ViewFeedback struct {
	Name         string `json:"name" validate:"min=1,max=255"`           // Имя пользователя
	Email        string `json:"email" validate:"max=255,regexp=^.+@.+$"` // Email пользователя
	Message      string `json:"message" validate:"nonzero"`              // Сообщение
	CaptchaValue string `json:"captchaValue" validate:"max=255"`         // Значение капчи
	CaptchaHash  string `json:"captchaHash" validate:"max=255"`          // Хэш капчи
}

type DtoFeedback struct {
	ID          int64     `db:"id"`          // Уникальный идентификатор обращения
	User_ID     int64     `db:"user_id"`     // Идентификатор пользователя
	Name        string    `db:"name"`        // Имя пользователя
	Email       string    `db:"email"`       // Email пользователя
	Message     string    `db:"message"`     // Сообщение
	Created     time.Time `db:"created"`     // Время обращения
	IP_Address  string    `db:"ip_address"`  // IP адрес обращения
	Reverse_DNS string    `db:"reverse_dns"` // Обратный DNS обращения
	User_Agent  string    `db:"user_agent"`  // User Agent обращения
}

// Конструктор создания объекта обращения в бд
func NewDtoFeedback(id int64, user_id int64, name string, email string, message string, created time.Time, ip_address string,
	reverse_dns string, user_agent string) *DtoFeedback {
	return &DtoFeedback{
		ID:          id,
		User_ID:     user_id,
		Name:        name,
		Email:       email,
		Message:     message,
		Created:     created,
		IP_Address:  ip_address,
		Reverse_DNS: reverse_dns,
		User_Agent:  user_agent,
	}
}

func (feedback *ViewFeedback) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(feedback, errors, req)
}
