package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

//Структура для организации сессии
type ViewSession struct {
	Login        string `json:"login" validate:"nonzero,min=1,max=255,regexp=^((.+@.+)|([0-9]*))$"` // Логин пользователя
	Password     string `json:"password" validate:"nonzero,min=1,max=255"`                          // Пароль пользователя
	CaptchaValue string `json:"captchaValue" validate:"max=255"`                                    // Значение капчи
	CaptchaHash  string `json:"captchaHash" validate:"max=255"`                                     // Хэш капчи
	Language     string `json:"language" validate:"max=10"`                                         // Язык пользователя
}

type ApiSession struct {
	Timeout     time.Time `json:"timeout"`      // Таймаут сессии
	AccessToken string    `json:"access-token"` // Токен доступа сессии
}

type DtoSession struct {
	AccessToken  string     `db:"token"`        // Ключ сессии
	UserID       int64      `db:"user_id"`      // Идентификатор пользователя сессии
	Roles        []UserRole `db:"-"`            // Массив значений уровней доступа пользователя сессии
	LastActivity time.Time  `db:"lastActivity"` // Время последнего использования сессии
	Language     string     `db:"language"`     // Язык пользователя сессии
}

// Конструктор создания объекта сессии в api
func NewApiSession(timeout time.Time, accesstoken string) *ApiSession {
	return &ApiSession{
		Timeout:     timeout,
		AccessToken: accesstoken,
	}
}

// Конструктор создания объекта сессии в бд
func NewDtoSession(accesstoken string, userid int64, roles []UserRole, lastactivity time.Time, language string) *DtoSession {
	return &DtoSession{
		AccessToken:  accesstoken,
		UserID:       userid,
		Roles:        roles,
		LastActivity: lastactivity,
		Language:     language,
	}
}

func (session *ViewSession) Validate(errors binding.Errors, req *http.Request) binding.Errors {

	return ValidateWithLanguage(session, errors, req, session.Language)
}
