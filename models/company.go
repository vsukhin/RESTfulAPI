package models

import (
	"errors"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Структура для организации хранения компании
type ViewCompany struct {
	Primary          bool                     `json:"primary"`                                                  // Основной
	Company_Type_ID  int                      `json:"legalFormId" validate:"nonzero"`                           // Идентификатор типа компании
	FullName_Rus     string                   `json:"nameLongRus" validate:"min=1,max=255"`                     // Полное русское название
	FullName_Eng     string                   `json:"nameLongEng" db:"fullname_eng" validate:"min=1,max=255"`   // Полное английское название
	ShortName_Rus    string                   `json:"nameShortRus" db:"shortname_rus" validate:"min=1,max=255"` // Краткое русское название
	ShortName_Eng    string                   `json:"nameShortEng" db:"shortname_eng" validate:"min=1,max=255"` // Краткое английское название
	Resident         bool                     `json:"resident"`                                                 // Резидент
	CompanyCodes     []ViewApiCompanyCode     `json:"codes"`                                                    // Коды компании
	CompanyAddresses []ViewApiCompanyAddress  `json:"addresses"`                                                // Адреса компании
	CompanyBanks     []ViewApiCompanyBank     `json:"banks"`                                                    // Банки компании
	CompanyStaff     []ViewApiCompanyEmployee `json:"staff"`                                                    // Сотрудники компании
	VAT              byte                     `json:"vat" validate:"min=0,max=100"`                             // НДС компании                                                    // Удален
}

type ApiMetaCompany struct {
	Total int64 `json:"count"` // Общее число компаний
}

type ApiShortCompany struct {
	ID            int64  `json:"id" db:"id"`                     // Уникальный идентификатор компании
	ShortName_Rus string `json:"nameShortRus" db:"nameShortRus"` // Краткое русское название
	ShortName_Eng string `json:"nameShortEng" db:"nameShortEng"` // Краткое английское название
	Unit_ID       int64  `json:"unitId" db:"unitId"`             // Идентификатор объединения
	Locked        bool   `json:"lock" db:"lock"`                 // Неизменяемая
	IsPrimary     bool   `json:"primary" db:"primary"`           // Юридическое лицо по умолчанию
}

type ApiMiddleCompany struct {
	Primary          bool                     `json:"primary" db:"primary"`             // Основной
	Company_Type_ID  int                      `json:"legalFormId" db:"company_type_id"` // Идентификатор типа компании
	FullName_Rus     string                   `json:"nameLongRus" db:"fullname_rus"`    // Полное русское название
	FullName_Eng     string                   `json:"nameLongEng" db:"fullname_eng"`    // Полное английское название
	ShortName_Rus    string                   `json:"nameShortRus" db:"shortname_rus"`  // Краткое русское название
	ShortName_Eng    string                   `json:"nameShortEng" db:"shortname_eng"`  // Краткое английское название
	Resident         bool                     `json:"resident" db:"resident"`           // Резидент
	CompanyCodes     []ViewApiCompanyCode     `json:"codes,omitempty" db:"-"`           // Коды компании
	CompanyAddresses []ViewApiCompanyAddress  `json:"addresses,omitempty" db:"-"`       // Адреса компании
	CompanyBanks     []ViewApiCompanyBank     `json:"banks,omitempty" db:"-"`           // Банки компании
	CompanyStaff     []ViewApiCompanyEmployee `json:"staff,omitempty" db:"-"`           // Сотрудники компании
	VAT              byte                     `json:"vat" db:"vat"`                     // НДС компании
	Locked           bool                     `json:"lock" db:"locked"`                 // Неизменяемая
	Deleted          bool                     `json:"del" db:"del"`                     // Удален

}

type ApiLongCompany struct {
	ID               int64                    `json:"id" db:"id"`                       // Уникальный идентификатор компании
	Primary          bool                     `json:"primary" db:"primary"`             // Основной
	Company_Type_ID  int                      `json:"legalFormId" db:"company_type_id"` // Идентификатор типа компании
	FullName_Rus     string                   `json:"nameLongRus" db:"fullname_rus"`    // Полное русское название
	FullName_Eng     string                   `json:"nameLongEng" db:"fullname_eng"`    // Полное английское название
	ShortName_Rus    string                   `json:"nameShortRus" db:"shortname_rus"`  // Краткое русское название
	ShortName_Eng    string                   `json:"nameShortEng" db:"shortname_eng"`  // Краткое английское название
	Resident         bool                     `json:"resident" db:"resident"`           // Резидент
	CompanyCodes     []ViewApiCompanyCode     `json:"codes,omitempty" db:"-"`           // Коды компании
	CompanyAddresses []ViewApiCompanyAddress  `json:"addresses,omitempty" db:"-"`       // Адреса компании
	CompanyBanks     []ViewApiCompanyBank     `json:"banks,omitempty" db:"-"`           // Банки компании
	CompanyStaff     []ViewApiCompanyEmployee `json:"staff,omitempty" db:"-"`           // Сотрудники компании
	VAT              byte                     `json:"vat" db:"vat"`                     // НДС компании
	Locked           bool                     `json:"lock" db:"locked"`                 // Неизменяемая
	Deleted          bool                     `json:"del" db:"del"`                     // Удален
}

type CompanySearch struct {
	ID            int64  `query:"id" search:"id"`                                            // Уникальный идентификатор компании
	ShortName_Rus string `query:"nameShortRus" search:"shortname_rus" group:"shortname_rus"` // Краткое русское название
	ShortName_Eng string `query:"nameShortEng" search:"shortname_eng" group:"shortname_eng"` // Краткое английское название
	Unit_ID       int64  `query:"unitId" search:"unit_id"`                                   // Идентификатор объединения
	Locked        bool   `query:"lock" search:"locked"`                                      // Неизменяемая
	IsPrimary     bool   `query:"primary" search:"primary"`                                  // Юридическое лицо по умолчанию
}

type DtoCompany struct {
	ID               int64                `db:"id"`              // Уникальный идентификатор компании
	Unit_ID          int64                `db:"unit_id"`         // Идентификатор объединения
	ShortName_Rus    string               `db:"shortname_rus"`   // Краткое русское название
	ShortName_Eng    string               `db:"shortname_eng"`   // Краткое английское название
	FullName_Rus     string               `db:"fullname_rus"`    // Полное русское название
	FullName_Eng     string               `db:"fullname_eng"`    // Полное английское название
	Created          time.Time            `db:"created"`         // Время создания
	Primary          bool                 `db:"primary"`         // Основной
	Active           bool                 `db:"active"`          // Aктивен
	Company_Type_ID  int                  `db:"company_type_id"` // Идентификатор типа компании
	Resident         bool                 `db:"resident"`        // Резидент
	VAT              byte                 `db:"vat"`             // НДС компании
	CompanyCodes     []DtoCompanyCode     `db:"-"`               // Коды компании
	CompanyAddresses []DtoCompanyAddress  `db:"-"`               // Адреса компании
	CompanyBanks     []DtoCompanyBank     `db:"-"`               // Банки компании
	CompanyStaff     []DtoCompanyEmployee `db:"-"`               // Сотрудники компании
	Locked           bool                 `db:"locked"`          // Неизменяемая
}

// Конструктор создания объекта компании в api
func NewApiMetaCompany(total int64) *ApiMetaCompany {
	return &ApiMetaCompany{
		Total: total,
	}
}

func NewApiShortCompany(id int64, shortname_rus string, shortname_eng string, unitid int64, locked bool, isprimary bool) *ApiShortCompany {
	return &ApiShortCompany{
		ID:            id,
		ShortName_Rus: shortname_rus,
		ShortName_Eng: shortname_eng,
		Unit_ID:       unitid,
		Locked:        locked,
		IsPrimary:     isprimary,
	}
}

func NewApiMiddleCompany(primary bool, company_type_id int, fullname_rus string, fullname_eng string, shortname_rus string,
	shortname_eng string, resident bool, companycodes []ViewApiCompanyCode, companyaddresses []ViewApiCompanyAddress,
	companybanks []ViewApiCompanyBank, companystaff []ViewApiCompanyEmployee, vat byte, locked bool, deleted bool) *ApiMiddleCompany {
	return &ApiMiddleCompany{
		Primary:          primary,
		Company_Type_ID:  company_type_id,
		FullName_Rus:     fullname_rus,
		FullName_Eng:     fullname_eng,
		ShortName_Rus:    shortname_rus,
		ShortName_Eng:    shortname_eng,
		Resident:         resident,
		CompanyCodes:     companycodes,
		CompanyAddresses: companyaddresses,
		CompanyBanks:     companybanks,
		CompanyStaff:     companystaff,
		VAT:              vat,
		Locked:           locked,
		Deleted:          deleted,
	}
}

func NewApiLongCompany(id int64, primary bool, company_type_id int, fullname_rus string, fullname_eng string, shortname_rus string,
	shortname_eng string, resident bool, companycodes []ViewApiCompanyCode, companyaddresses []ViewApiCompanyAddress,
	companybanks []ViewApiCompanyBank, companystaff []ViewApiCompanyEmployee, vat byte, locked bool, deleted bool) *ApiLongCompany {
	return &ApiLongCompany{
		ID:               id,
		Primary:          primary,
		Company_Type_ID:  company_type_id,
		FullName_Rus:     fullname_rus,
		FullName_Eng:     fullname_eng,
		ShortName_Rus:    shortname_rus,
		ShortName_Eng:    shortname_eng,
		Resident:         resident,
		CompanyCodes:     companycodes,
		CompanyAddresses: companyaddresses,
		CompanyBanks:     companybanks,
		CompanyStaff:     companystaff,
		VAT:              vat,
		Locked:           locked,
		Deleted:          deleted,
	}
}

// Конструктор создания объекта компании в бд
func NewDtoCompany(id int64, unit_id int64, shortname_rus string, shortname_eng string, fullname_rus string, fullname_eng string,
	created time.Time, primary bool, active bool, company_type_id int, resident bool, vat byte, companycodes []DtoCompanyCode,
	companyaddresses []DtoCompanyAddress, companybanks []DtoCompanyBank, companystaff []DtoCompanyEmployee, locked bool) *DtoCompany {
	return &DtoCompany{
		ID:               id,
		Unit_ID:          unit_id,
		ShortName_Rus:    shortname_rus,
		ShortName_Eng:    shortname_eng,
		FullName_Rus:     fullname_rus,
		FullName_Eng:     fullname_eng,
		Created:          created,
		Primary:          primary,
		Active:           active,
		Company_Type_ID:  company_type_id,
		Resident:         resident,
		VAT:              vat,
		CompanyCodes:     companycodes,
		CompanyAddresses: companyaddresses,
		CompanyBanks:     companybanks,
		CompanyStaff:     companystaff,
		Locked:           locked,
	}
}

func (company *CompanySearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, company), nil
}

func (company *CompanySearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, company)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "unitId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "nameShortRus":
		fallthrough
	case "nameShortEng":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "lock":
		fallthrough
	case "primary":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (company *CompanySearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(company)
}

func (company *ViewCompany) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, companycode := range company.CompanyCodes {
		errors = Validate(&companycode, errors, req)
	}
	for _, companyaddress := range company.CompanyAddresses {
		errors = Validate(&companyaddress, errors, req)
	}
	for _, companybank := range company.CompanyBanks {
		errors = Validate(&companybank, errors, req)
	}
	for _, companyemployee := range company.CompanyStaff {
		errors = Validate(&companyemployee, errors, req)
	}
	return Validate(company, errors, req)
}
