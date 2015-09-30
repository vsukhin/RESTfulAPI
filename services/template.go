package services

import (
	"bytes"
	"net/http"
)

const (
	TEMPLATE_LAYOUT                = "layout.tpl.html"                   // Макет электронных писем
	TEMPLATE_EMAIL_CONFIRMATION    = "email_added_confirmation.tpl.html" // Подтверждение e-mail при добавлении
	TEMPLATE_PASSWORD_REGISTRATION = "password_registration.tpl.html"    // Установка пароля при регистрации
	TEMPLATE_PASSWORD_RECOVERY     = "password_recovery.tpl.html"        // Сброс пароля зарегистрированного пользователя
	TEMPLATE_SUBSCRIPTION          = "subscription.tpl.html"             // Письмо подписки на новости
	TEMPLATE_FEEDBACK              = "sayhello.tpl.html"                 // Письмо обратной связи
	TEMPLATE_CONFIRMATION          = "confirmation.tpl.html"             // Подтверждение успешности регистрации
	TEMPLATE_MATCHING              = "matching.tpl.html"                 // Письмо запроса акта сверки
	TEMPLATE_INVOICE               = "invoice.tpl.html"                  // Счет-фактура
	TEMPLATE_DIRECTORY_EMAILS      = "/mailers"
)

type TemplateRepository interface {
	GenerateText(object interface{}, name string, directory string, layout string) (buf *bytes.Buffer, err error)
	GenerateHTML(name string, w http.ResponseWriter, object interface{}) (err error)
}

type TemplateService struct {
}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}
