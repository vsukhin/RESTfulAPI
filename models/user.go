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

type UserRole int

const (
	USER_ROLE_DEVELOPER UserRole = iota + 1
	USER_ROLE_ADMINISTRATOR
	USER_ROLE_SUPPLIER
	USER_ROLE_CUSTOMER
)

// Структура для хранения данных пользователя
type ViewUser struct {
	Login        string `json:"email" validate:"nonzero,min=1,max=255,regexp=^.+@.+$"` // Логин пользователя
	CaptchaValue string `json:"captchaValue" validate:"max=255"`                       // Значение капчи
	CaptchaHash  string `json:"captchaHash" validate:"max=255"`                        // Хэш капчи
}

type ChangeUser struct {
	Surname    string `json:"surname" validate:"nonzero,min=1,max=255"` // Фамилия пользователя
	Name       string `json:"name" validate:"nonzero,min=1,max=255"`    // Логин пользователя
	MiddleName string `json:"middleName"  validate:"max=255"`           // Отчество пользователя
	WorkPhone  string `json:"workPhone" validate:"max=25"`              // Рабочий телефон
	JobTitle   string `json:"jobTitle" validate:"max=255"`              // Должность
	Language   string `json:"language" validate:"max=10"`               // Язык пользователя по умолчанию
}

type ViewApiUserFull struct {
	Creator_ID   int64                `json:"userId" db:"user_id"`                                    // Уникальный id создателя
	Unit_ID      int64                `json:"unitId" db:"unit_id"`                                    // Уникальный id объединения
	UnitAdmin    bool                 `json:"unitAdmin" db:"unitAdmin"`                               // Администратор объединения
	Active       bool                 `json:"active" db:"active"`                                     // Пользователь активен
	Confirmed    bool                 `json:"confirmed" db:"confirmed"`                               // Пользователь подтвержден
	Surname      string               `json:"surname" db:"surname" validate:"nonzero,min=1,max=255"`  // Фамилия пользователя
	Name         string               `json:"name" db:"name" validate:"nonzero,min=1,max=255"`        // Имя пользователя
	MiddleName   string               `json:"middleName" db:"middleName" validate:"max=255"`          // Отчество пользователя
	WorkPhone    string               `json:"workPhone" db:"workPhone" validate:"max=25"`             // Рабочий телефон
	JobTitle     string               `json:"jobTitle" db:"jobTitle" validate:"max=255"`              // Должность
	Language     string               `json:"language" db:"language" validate:"nonzero,min=1,max=10"` // Язык пользователя по умолчанию
	Roles        []UserRole           `json:"groups,omitempty" db:"-"`                                // Массив значений уровней доступа пользователя
	Emails       []ViewApiEmail       `json:"emails,omitempty" db:"-"`                                // Массив email
	MobilePhones []ViewApiMobilePhone `json:"mobilePhones,omitempty" db:"-"`                          // Массив мобильных телефонов
}

type ApiUserTiny struct {
	ID int64 `json:"id" db:"id"` // Уникальный id пользователя
}

type ApiUserShort struct {
	ID         int64     `json:"id" db:"id"`                   // Уникальный id пользователя
	Blocked    bool      `json:"blocked" db:"blocked"`         // Пользователь активен
	Confirmed  bool      `json:"confirmed" db:"confirmed"`     // Пользователь подтвержден
	LastLogin  time.Time `json:"lastLoginAt" db:"lastLoginAt"` // Время последнего логина
	Surname    string    `json:"surname" db:"surname"`         // Фамилия пользователя
	Name       string    `json:"name" db:"name"`               // Имя пользователя
	MiddleName string    `json:"middleName" db:"middleName"`   // Отчество пользователя
}

type ApiUserLong struct {
	ID         int64  `json:"id" db:"id"`                 // Уникальный id пользователя
	UnitID     int64  `json:"unitId" db:"unit_id"`        // Уникальный id объединения
	UnitAdmin  bool   `json:"unitAdmin" db:"unitAdmin"`   // Администратор объединения
	Active     bool   `json:"active" db:"active"`         // Пользователь активен
	Confirmed  bool   `json:"confirmed" db:"confirmed"`   // Пользователь подтвержден
	Surname    string `json:"surname" db:"surname"`       // Фамилия пользователя
	Name       string `json:"name" db:"name"`             // Имя пользователя
	MiddleName string `json:"middleName" db:"middleName"` // Отчество пользователя
	WorkPhone  string `json:"workPhone" db:"workPhone"`   // Рабочий телефон
	JobTitle   string `json:"jobTitle" db:"jobTitle"`     // Должность
	Language   string `json:"language" db:"language"`     // Язык пользователя по умолчанию
}

type ApiUserMeta struct {
	NumOfRows int64 `json:"rows"` // Число строк
}

type UserSearch struct {
	ID         int64  `query:"id" search:"id"`                 // id пользователя
	Blocked    bool   `query:"blocked" search:"(not active)"`  // Пользователь заблокирован
	Confirmed  bool   `query:"confirmed" search:"confirmed"`   // Пользователь подтвержден
	LastLogin  string `query:"lastLoginAt" search:"lastLogin"` // Последний заход
	Surname    string `query:"surname" search:"surname"`       // Фамилия пользователя
	Name       string `query:"name" search:"name"`             // Имя пользователя
	MiddleName string `query:"middleName" search:"middleName"` // Отчество пользователя
}

type DtoUser struct {
	ID           int64             `db:"id"`         // Уникальный id пользователя
	Creator_ID   int64             `db:"user_id"`    // Уникальный id создателя
	UnitID       int64             `db:"unit_id"`    // Уникальный id объединения
	Roles        []UserRole        `db:"-"`          // Массив значений уровней доступа пользователя
	UnitAdmin    bool              `db:"unitAdmin"`  // Администратор объединения
	Active       bool              `db:"active"`     // Пользователь активен
	Confirmed    bool              `db:"confirmed"`  // Пользователь подтвержден
	Created      time.Time         `db:"created"`    // Время создания пользователя
	LastLogin    time.Time         `db:"lastLogin"`  // Время последнего логина
	Password     string            `db:"password"`   // Хэш пароля
	Surname      string            `db:"surname"`    // Фамилия пользователя
	Name         string            `db:"name"`       // Имя пользователя
	MiddleName   string            `db:"middleName"` // Отчество пользователя
	WorkPhone    string            `db:"workPhone"`  // Рабочий телефон
	JobTitle     string            `db:"jobTitle"`   // Должность
	Code         string            `db:"code"`       // Koд подтверждения пользователя
	Language     string            `db:"language"`   // Язык пользователя по умолчанию
	Emails       *[]DtoEmail       `db:"-"`          // Массив email
	MobilePhones *[]DtoMobilePhone `db:"-"`          // Массив мобильных телефонов
}

// Конструктор создания объекта пользователя в api
func NewApiUserTiny(id int64) *ApiUserTiny {
	return &ApiUserTiny{
		ID: id,
	}
}

func NewApiUserShort(id int64, blocked bool, confirmed bool, created time.Time,
	lastlogin time.Time, surname string, name string, middlename string) *ApiUserShort {
	return &ApiUserShort{
		ID:         id,
		Blocked:    blocked,
		Confirmed:  confirmed,
		LastLogin:  lastlogin,
		Surname:    surname,
		Name:       name,
		MiddleName: middlename,
	}
}

func NewApiUserLong(id int64, unitid int64, unitadmin bool, active bool, confirmed bool, surname string,
	name string, middleName string, workphone string, jobtitle string, language string) *ApiUserLong {
	return &ApiUserLong{
		ID:         id,
		UnitID:     unitid,
		UnitAdmin:  unitadmin,
		Active:     active,
		Confirmed:  confirmed,
		Surname:    surname,
		Name:       name,
		MiddleName: middleName,
		WorkPhone:  workphone,
		JobTitle:   jobtitle,
		Language:   language,
	}
}

func NewViewApiUserFull(creator_id int64, unit_id int64, unitadmin bool, active bool, confirmed bool,
	surname string, name string, middlename string, workphone string, jobtitle string, language string,
	roles []UserRole, emails []ViewApiEmail, mobilephones []ViewApiMobilePhone) *ViewApiUserFull {
	return &ViewApiUserFull{
		Creator_ID:   creator_id,
		Unit_ID:      unit_id,
		UnitAdmin:    unitadmin,
		Active:       active,
		Confirmed:    confirmed,
		Surname:      surname,
		Name:         name,
		MiddleName:   middlename,
		WorkPhone:    workphone,
		JobTitle:     jobtitle,
		Language:     language,
		Roles:        roles,
		Emails:       emails,
		MobilePhones: mobilephones,
	}
}

func NewApiUserMeta(numofrows int64) *ApiUserMeta {
	return &ApiUserMeta{
		NumOfRows: numofrows,
	}
}

// Конструктор создания объекта пользователя в бд
func NewDtoUser(id int64, creator_id int64, unitid int64, roles []UserRole,
	unitadmin bool, active bool, confirmed bool, created time.Time, lastlogin time.Time,
	password string, surname string, name string, middlename string, workphone string, jobtitle string,
	code string, language string, emails *[]DtoEmail, mobilephones *[]DtoMobilePhone) *DtoUser {
	return &DtoUser{
		ID:           id,
		Creator_ID:   creator_id,
		UnitID:       unitid,
		Roles:        roles,
		UnitAdmin:    unitadmin,
		Active:       active,
		Confirmed:    confirmed,
		Created:      created,
		LastLogin:    lastlogin,
		Password:     password,
		Surname:      surname,
		Name:         name,
		MiddleName:   middlename,
		WorkPhone:    workphone,
		JobTitle:     jobtitle,
		Code:         code,
		Language:     language,
		Emails:       emails,
		MobilePhones: mobilephones,
	}
}

func (user *ViewUser) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(user, errors, req)
}

func (user *ViewApiUserFull) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, email := range user.Emails {
		errors = ValidateWithLanguage(&email, errors, req, email.Language)
	}
	for _, mobilephone := range user.MobilePhones {
		errors = ValidateWithLanguage(&mobilephone, errors, req, mobilephone.Language)
	}
	return ValidateWithLanguage(user, errors, req, user.Language)
}

func (user *ChangeUser) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return ValidateWithLanguage(user, errors, req, user.Language)
}

func (user *UserSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, user), nil
}

func (user *UserSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
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
			errValue = errors.New("Wrong field value")
			break
		}
		outvalue = "'" + invalue + "'"
	case "blocked":
		fallthrough
	case "confirmed":
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

func (user *UserSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(user)
}
