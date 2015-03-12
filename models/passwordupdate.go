package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

//Структура для смены пароля
type PasswordUpdate struct {
	Code  string `json:"code" validate:"nonzero,min=1,max=255"`     // Код для смены пароля
	Value string `json:"password" validate:"nonzero,min=1,max=255"` // Значение пароля
}

// Структура для подтверждения пароля
type EmailConfirm struct {
	ConfirmationToken string `json:"confirmation_token" validate:"nonzero,min=1,max=255"` // Токен подтверждения смены пароля
}

func (password *PasswordUpdate) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(password, errors, req)
}

func (email *EmailConfirm) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(email, errors, req)
}
