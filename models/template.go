package models

import (
	"application/config"
)

// Структура для организации хранения данных шаблона
type DtoTemplate struct {
	Email    string // Email для отправки
	Language string // Язык шаблона
	Host     string // URL ссылки
}

type DtoCodeTemplate struct {
	*DtoTemplate
	Code string // Код подтверждения
}

type DtoDualCodeTemplate struct {
	*DtoTemplate
	Subscr_Code   string // Код подписки
	Unsubscr_Code string // Код отписки
}

// Конструктор создания объекта шаблона
func NewDtoTemplate(email string, language string, host string) *DtoTemplate {
	return &DtoTemplate{
		Email:    email,
		Language: language,
		Host:     host,
	}
}

func NewDtoCodeTemplate(dtotemplate *DtoTemplate, code string) *DtoCodeTemplate {
	return &DtoCodeTemplate{
		DtoTemplate: dtotemplate,
		Code:        code,
	}
}

func NewDtoDualCodeTemplate(dtotemplate *DtoTemplate, subscr_code string, unsubscr_code string) *DtoDualCodeTemplate {
	return &DtoDualCodeTemplate{
		DtoTemplate:   dtotemplate,
		Subscr_Code:   subscr_code,
		Unsubscr_Code: unsubscr_code,
	}
}

func (template *DtoTemplate) GetResource() (Resource config.Resource) {
	return config.Localization[template.Language]
}
