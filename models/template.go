package models

import (
	"application/config"
)

//Структура для организации хранения данных шаблона
type DtoTemplate struct {
	Email    string // Email для отправки
	Language string // Язык шаблона
	Host     string // URL ссылки
	Code     string // Код подтверждения
}

// Конструктор создания объекта шаблона
func NewDtoTemplate(email string, language string, host string, code string) *DtoTemplate {
	return &DtoTemplate{
		Email:    email,
		Language: language,
		Host:     host,
		Code:     code,
	}
}

func (template *DtoTemplate) GetResource() (Resource config.Resource) {
	return config.Localization[template.Language]
}
