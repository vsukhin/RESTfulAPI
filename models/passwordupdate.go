package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для смены пароля
type PasswordUpdate struct {
	Code  string `json:"code" validate:"min=1,max=255"`     // Код для смены пароля
	Value string `json:"password" validate:"min=1,max=255"` // Значение пароля
}

type ChangePassword struct {
	OldPassword     string `json:"passwordOld" validate:"min=1,max=255"`     // Старое значение пароля
	NewPassword     string `json:"passwordNew" validate:"min=1,max=255"`     // Новое значение пароля
	ConfirmPassword string `json:"passwordConfirm" validate:"min=1,max=255"` // Подтверждение нового значения пароля
}

// Структура для подтверждения пароля
type EmailConfirm struct {
	ConfirmationToken string `json:"confirmation_token" validate:"min=1,max=255"` // Токен подтверждения смены пароля
}

func (password *PasswordUpdate) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(password, errors, req)
}

func (password *ChangePassword) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(password, errors, req)
}

func (email *EmailConfirm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(email, errors, req)
}
