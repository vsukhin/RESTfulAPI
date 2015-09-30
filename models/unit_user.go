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

// Структура для хранения данных пользователя объединения
type SearchUnitUser struct {
	ID         int64  `query:"id" search:"id"`                                    // id пользователя
	Blocked    bool   `query:"blocked" search:"(not active)"`                     // Пользователь заблокирован
	Confirmed  bool   `query:"confirmed" search:"confirmed"`                      // Пользователь подтвержден
	LastLogin  string `query:"lastLoginAt" search:"lastLogin"`                    // Последний заход
	UnitAdmin  bool   `query:"unitAdmin" search:"unitAdmin"`                      // Администратор объединения
	Surname    string `query:"surname" search:"surname" group:"surname"`          // Фамилия пользователя
	Name       string `query:"name" search:"name" group:"name"`                   // Имя пользователя
	MiddleName string `query:"middleName" search:"middleName" group:"middleName"` // Отчество пользователя
}

type ViewUnitUser struct {
	UnitAdmin    bool              `json:"unitAdmin"`                        // Администратор объединения
	Surname      string            `json:"surname" validate:"max=255"`       // Фамилия пользователя
	Name         string            `json:"name" validate:"max=255"`          // Имя пользователя
	MiddleName   string            `json:"middleName" validate:"max=255"`    // Отчество пользователя
	WorkPhone    string            `json:"workPhone" validate:"max=25"`      // Рабочий телефон
	JobTitle     string            `json:"jobTitle" validate:"max=255"`      // Должность
	Language     string            `json:"language" validate:"min=1,max=10"` // Язык пользователя по умолчанию
	Emails       []ViewEmail       `json:"emails"`                           // Массив email
	MobilePhones []ViewMobilePhone `json:"mobilePhones"`                     // Массив мобильных телефонов
}

type ApiMetaUnitUser struct {
	Total             int64 `json:"count"`        // Всего
	NumOfNotConfirmed int64 `json:"notConfirmed"` // Количество неподтвержденных
	NumOfAdmins       int64 `json:"admins"`       // Количество администраторов
}

type ApiSearchUnitUser struct {
	ID         int64     `json:"id" db:"id"`                   // id пользователя
	Blocked    bool      `json:"blocked" db:"blocked"`         // Пользователь заблокирован
	Confirmed  bool      `json:"confirmed" db:"confirmed"`     // Пользователь подтвержден
	LastLogin  time.Time `json:"lastLoginAt" db:"lastLoginAt"` // Последний заход
	UnitAdmin  bool      `json:"unitAdmin" db:"unitAdmin"`     // Администратор объединения
	Surname    string    `json:"surname" db:"surname"`         // Фамилия пользователя
	Name       string    `json:"name" db:"name"`               // Имя пользователя
	MiddleName string    `json:"middleName" db:"middleName"`   // Отчество пользователя
}

type ApiShortUnitUser struct {
	UnitAdmin    bool                 `json:"unitAdmin" db:"unitAdmin"`      // Администратор объединения
	Active       bool                 `json:"active" db:"active"`            // Пользователь активен
	Confirmed    bool                 `json:"confirmed" db:"confirmed"`      // Пользователь подтвержден
	Surname      string               `json:"surname" db:"surname"`          // Фамилия пользователя
	Name         string               `json:"name" db:"name"`                // Имя пользователя
	MiddleName   string               `json:"middleName" db:"middleName"`    // Отчество пользователя
	WorkPhone    string               `json:"workPhone" db:"workPhone"`      // Рабочий телефон
	JobTitle     string               `json:"jobTitle" db:"jobTitle"`        // Должность
	Language     string               `json:"language" db:"language"`        // Язык пользователя по умолчанию
	Emails       []ViewApiEmail       `json:"emails,omitempty" db:"-"`       // Массив email
	MobilePhones []ViewApiMobilePhone `json:"mobilePhones,omitempty" db:"-"` // Массив мобильных телефонов
}

type ApiLongUnitUser struct {
	UnitAdmin    bool                 `json:"unitAdmin" db:"unitAdmin"`      // Администратор объединения
	Active       bool                 `json:"active" db:"active"`            // Пользователь активен
	Confirmed    bool                 `json:"confirmed" db:"confirmed"`      // Пользователь подтвержден
	Surname      string               `json:"surname" db:"surname"`          // Фамилия пользователя
	Name         string               `json:"name" db:"name"`                // Имя пользователя
	MiddleName   string               `json:"middleName" db:"middleName"`    // Отчество пользователя
	WorkPhone    string               `json:"workPhone" db:"workPhone"`      // Рабочий телефон
	JobTitle     string               `json:"jobTitle" db:"jobTitle"`        // Должность
	Language     string               `json:"language" db:"language"`        // Язык пользователя по умолчанию
	Roles        []UserRole           `json:"groups,omitempty" db:"-"`       // Массив значений уровней доступа пользователя
	Emails       []ViewApiEmail       `json:"emails,omitempty" db:"-"`       // Массив email
	MobilePhones []ViewApiMobilePhone `json:"mobilePhones,omitempty" db:"-"` // Массив мобильных телефонов
}

// Конструктор создания объекта пользователя объединения в api
func NewApiMetaUnitUser(total int64, numofnotconfirmed int64, numofadmins int64) *ApiMetaUnitUser {
	return &ApiMetaUnitUser{
		Total:             total,
		NumOfNotConfirmed: numofnotconfirmed,
		NumOfAdmins:       numofadmins,
	}
}

func NewApiSearchUnitUser(id int64, blocked bool, confirmed bool, lastlogin time.Time, unitadmin bool, surname string,
	name string, middlename string) *ApiSearchUnitUser {
	return &ApiSearchUnitUser{
		ID:         id,
		Blocked:    blocked,
		Confirmed:  confirmed,
		LastLogin:  lastlogin,
		UnitAdmin:  unitadmin,
		Surname:    surname,
		Name:       name,
		MiddleName: middlename,
	}
}

func NewApiShortUnitUser(unitadmin bool, active bool, confirmed bool, surname string, name string, middlename string,
	workphone string, jobtitle string, language string, emails []ViewApiEmail, mobilephones []ViewApiMobilePhone) *ApiShortUnitUser {
	return &ApiShortUnitUser{
		UnitAdmin:    unitadmin,
		Active:       active,
		Confirmed:    confirmed,
		Surname:      surname,
		Name:         name,
		MiddleName:   middlename,
		WorkPhone:    workphone,
		JobTitle:     jobtitle,
		Language:     language,
		Emails:       emails,
		MobilePhones: mobilephones,
	}
}

func NewApiLongUnitUser(unitadmin bool, active bool, confirmed bool, surname string, name string, middlename string,
	workphone string, jobtitle string, language string, roles []UserRole, emails []ViewApiEmail, mobilephones []ViewApiMobilePhone) *ApiLongUnitUser {
	return &ApiLongUnitUser{
		UnitAdmin:    unitadmin,
		Active:       active,
		Confirmed:    confirmed,
		Surname:      surname,
		Name:         name,
		MiddleName:   middlename,
		WorkPhone:    workphone,
		JobTitle:     jobtitle,
		Roles:        roles,
		Language:     language,
		Emails:       emails,
		MobilePhones: mobilephones,
	}
}

func (user *SearchUnitUser) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, user), nil
}

func (user *SearchUnitUser) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, user)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "surname":
		fallthrough
	case "name":
		fallthrough
	case "middleName":
		fallthrough
	case "lastLoginAt":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "blocked":
		fallthrough
	case "confirmed":
		fallthrough
	case "unitAdmin":
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

func (user *SearchUnitUser) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(user)
}

func (user *ViewUnitUser) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, email := range user.Emails {
		errors = ValidateWithLanguage(&email, errors, req, email.Language)
	}
	for _, mobilephone := range user.MobilePhones {
		errors = ValidateWithLanguage(&mobilephone, errors, req, mobilephone.Language)
	}
	return ValidateWithLanguage(user, errors, req, user.Language)
}
