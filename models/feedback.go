package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
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
	ID           int64  `db:"id"`           // Уникальный идентификатор обращения
	User_ID      int64  `db:"user_id"`      // Идентификатор пользователя
	Name         string `db:"name"`         // Имя пользователя
	Email        string `db:"email"`        // Email пользователя
	Message      string `db:"message"`      // Сообщение
	AccessLog_ID int64  `db:"accesslog_id"` // Идентификатор лога
}

// Конструктор создания объекта обращения в бд
func NewDtoFeedback(id int64, user_id int64, name string, email string, message string, accesslog_id int64) *DtoFeedback {
	return &DtoFeedback{
		ID:           id,
		User_ID:      user_id,
		Name:         name,
		Email:        email,
		Message:      message,
		AccessLog_ID: accesslog_id,
	}
}

func (feedback *ViewFeedback) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(feedback, errors, req)
}
