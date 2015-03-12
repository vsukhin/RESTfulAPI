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

// Структура для хранения данных пользователя
type ViewUser struct {
	Login        string `json:"email" validate:"nonzero,min=1,max=255,regexp=^.+@.+$"` // Логин пользователя
	CaptchaValue string `json:"captchaValue" validate:"max=255"`                       // Значение капчи
	CaptchaHash  string `json:"captchaHash" validate:"max=255"`                        // Хэш капчи
}

type ApiUserTiny struct {
	ID int64 `json:"id" db:"id"` // Уникальный id пользователя
}

type ApiUserShort struct {
	ID        int64     `json:"id" db:"id"`                   // Уникальный id пользователя
	Login     string    `json:"login" db:"login"`             // Логин пользователя
	Blocked   bool      `json:"blocked" db:"blocked"`         // Пользователь активен
	Confirmed bool      `json:"confirmed" db:"confirmed"`     // Пользователь подтвержден
	LastLogin time.Time `json:"lastLoginAt" db:"lastLoginAt"` // Время последнего логина
	Name      string    `json:"name" db:"name"`               // Имя пользователя
}

type ApiUserLong struct {
	ID        int64  `json:"id" db:"id"`               // Уникальный id пользователя
	UnitID    int64  `json:"unitId" db:"unit_id"`      // Уникальный id объединения
	UnitAdmin bool   `json:"unitAdmin" db:"unitAdmin"` // Администратор объединения
	Active    bool   `json:"active" db:"active"`       // Пользователь активен
	Confirmed bool   `json:"confirmed" db:"confirmed"` // Пользователь подтвержден
	Name      string `json:"name" db:"name"`           // Имя пользователя
	Language  string `json:"language" db:"language"`   // Язык пользователя по умолчанию
}

type ViewApiUserFull struct {
	Creator_ID int64          `json:"userId" db:"user_id"`                                    // Уникальный id создателя
	Unit_ID    int64          `json:"unitId" db:"unit_id"`                                    // Уникальный id объединения
	UnitAdmin  bool           `json:"unitAdmin" db:"unitAdmin"`                               // Администратор объединения
	Active     bool           `json:"active" db:"active"`                                     // Пользователь активен
	Confirmed  bool           `json:"confirmed" db:"confirmed"`                               // Пользователь подтвержден
	Name       string         `json:"name" db:"name" validate:"nonzero,min=1,max=255"`        // Имя пользователя
	Language   string         `json:"language" db:"language" validate:"nonzero,min=1,max=10"` // Язык пользователя по умолчанию
	Roles      []UserRole     `json:"groups,omitempty" db:"-"`                                // Массив значений уровней доступа пользователя
	Emails     []ViewApiEmail `json:"emails,omitempty" db:"-"`                                // Массив email
}

type ApiUserMeta struct {
	NumOfRows int64 `json:"rows"` // Число строк
}

type ChangeUser struct {
	Name     string `json:"name" validate:"nonzero,min=1,max=255"` // Логин пользователя
	Language string `json:"language" validate:"max=10"`            // Язык пользователя по умолчанию
}

type UserSearch struct {
	ID        int64  `query:"id" search:"u.id"`                 // id пользователя
	Login     string `query:"login" search:"e.email"`           // Логин пользователя
	Blocked   bool   `query:"blocked" search:"u.active"`        // Пользователь заблокирован
	Confirmed bool   `query:"confirmed" search:"u.confirmed"`   // Пользователь подтвержден
	LastLogin string `query:"lastLoginAt" search:"u.lastLogin"` // Последний заход
	Name      string `query:"name" search:"u.name"`             // Имя пользователя
}

type DtoUser struct {
	ID         int64       `db:"id"`        // Уникальный id пользователя
	Creator_ID int64       `db:"user_id"`   // Уникальный id создателя
	UnitID     int64       `db:"unit_id"`   // Уникальный id объединения
	Roles      []UserRole  `db:"-"`         // Массив значений уровней доступа пользователя
	UnitAdmin  bool        `db:"unitAdmin"` // Администратор объединения
	Active     bool        `db:"active"`    // Пользователь активен
	Confirmed  bool        `db:"confirmed"` // Пользователь подтвержден
	Created    time.Time   `db:"created"`   // Время создания пользователя
	LastLogin  time.Time   `db:"lastLogin"` // Время последнего логина
	Password   string      `db:"password"`  // Хэш пароля
	Name       string      `db:"name"`      // Имя пользователя
	Code       string      `db:"code"`      // Koд подтверждения пользователя
	Language   string      `db:"language"`  // Язык пользователя по умолчанию
	Emails     *[]DtoEmail `db:"-"`         // Массив email
}

// Конструктор создания объекта пользователя в api
func NewApiUserTiny(id int64) *ApiUserTiny {
	return &ApiUserTiny{
		ID: id,
	}
}

func NewApiUserShort(id int64, login string, blocked bool, confirmed bool, created time.Time,
	lastlogin time.Time, name string) *ApiUserShort {
	return &ApiUserShort{
		ID:        id,
		Login:     login,
		Blocked:   blocked,
		Confirmed: confirmed,
		LastLogin: lastlogin,
		Name:      name,
	}
}

func NewApiUserLong(id int64, unitid int64, unitadmin bool, active bool, confirmed bool,
	name string, language string) *ApiUserLong {
	return &ApiUserLong{
		ID:        id,
		UnitID:    unitid,
		UnitAdmin: unitadmin,
		Active:    active,
		Confirmed: confirmed,
		Name:      name,
		Language:  language,
	}
}

func NewViewApiUserFull(creator_id int64, unit_id int64, unitadmin bool, active bool, confirmed bool,
	name string, language string, roles []UserRole, emails []ViewApiEmail) *ViewApiUserFull {
	return &ViewApiUserFull{
		Creator_ID: creator_id,
		Unit_ID:    unit_id,
		UnitAdmin:  unitadmin,
		Active:     active,
		Confirmed:  confirmed,
		Name:       name,
		Language:   language,
		Roles:      roles,
		Emails:     emails,
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
	password string, name string, code string, language string, emails *[]DtoEmail) *DtoUser {
	return &DtoUser{
		ID:         id,
		Creator_ID: creator_id,
		UnitID:     unitid,
		Roles:      roles,
		UnitAdmin:  unitadmin,
		Active:     active,
		Confirmed:  confirmed,
		Created:    created,
		LastLogin:  lastlogin,
		Password:   password,
		Name:       name,
		Code:       code,
		Language:   language,
		Emails:     emails,
	}
}

func (user *ViewUser) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(user, errors, req)
}

func (user *ViewApiUserFull) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, apiEmail := range user.Emails {
		errors = ValidateWithLanguage(&apiEmail, errors, req, apiEmail.Language)
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
	case "login":
		invalue = strings.ToLower(invalue)
		fallthrough
	case "name":
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
		if infield == "blocked" {
			val = !val
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Uknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (user *UserSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(user)
}
