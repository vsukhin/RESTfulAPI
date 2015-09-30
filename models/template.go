package models

import (
	"application/config"
	"time"
)

// Структура для организации хранения данных шаблона
type DtoTemplate struct {
	Email      string    // Email для отправки
	Language   string    // Язык шаблона
	Host       string    // URL ссылки
	Created    time.Time // Дата и время создания
	IP_Address string    // IP адрес

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

type DtoHTMLTemplate struct {
	Content  string // Содержание
	Language string // Язык шаблона
}

type DtoCompanyTemplate struct {
	Name       string // Название
	INN        string // ИНН
	KPP        string // КПП
	Address    string // Адрес
	CEO        string // Генеральный директор
	Accountant string // Главный бухгалтер
}

type DtoInvoiceTemplate struct {
	Invoice  ApiFullInvoice     // Счет
	Bank     ViewApiCompanyBank // Банк для платежа
	Seller   DtoCompanyTemplate // Поставщик
	Buyer    DtoCompanyTemplate // Покупатель
	Contract DtoContract        // Договор
}

// Конструктор создания объекта шаблона
func NewDtoTemplate(email string, language string, host string, created time.Time, ip_address string) *DtoTemplate {
	return &DtoTemplate{
		Email:      email,
		Language:   language,
		Host:       host,
		Created:    created,
		IP_Address: ip_address,
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

func NewDtoHTMLTemplate(content string, language string) *DtoHTMLTemplate {
	return &DtoHTMLTemplate{
		Content:  content,
		Language: language,
	}
}

func NewDtoCompanyTemplate(name, inn, kpp, address, ceo, accountant string) *DtoCompanyTemplate {
	return &DtoCompanyTemplate{
		Name:       name,
		INN:        inn,
		KPP:        kpp,
		Address:    address,
		CEO:        ceo,
		Accountant: accountant,
	}
}

func NewDtoInvoiceTemplate(invoice ApiFullInvoice, bank ViewApiCompanyBank, seller DtoCompanyTemplate, buyer DtoCompanyTemplate,
	contract DtoContract) *DtoInvoiceTemplate {
	return &DtoInvoiceTemplate{
		Invoice:  invoice,
		Bank:     bank,
		Seller:   seller,
		Buyer:    buyer,
		Contract: contract,
	}
}
func (template *DtoTemplate) GetResource() (Resource config.Resource) {
	return config.Localization[template.Language]
}

func (template *DtoHTMLTemplate) GetResource() (Resource config.Resource) {
	return config.Localization[template.Language]
}
