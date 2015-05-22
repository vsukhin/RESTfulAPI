package models

import (
	"time"
)

// Структура для организации хранения капчи
type ApiCaptcha struct {
	Hash  string `json:"captchaHash"`  // Уникальный hash капчи
	Image string `json:"captchaImage"` // Картинка капчи
}

type DtoCaptcha struct {
	Hash    string    `db:"hash"`    // Уникальный hash капчи
	Image   []byte    `db:"image"`   // Картинка капчи
	Value   string    `db:"value"`   // Значение капчи
	Created time.Time `db:"created"` // Время создания капчи
	InUse   bool      `db:"inUse"`   // Показывалась
}

// Конструктор создания объекта капчи в api
func NewApiCaptcha(hash string, image string) *ApiCaptcha {
	return &ApiCaptcha{
		Hash:  hash,
		Image: image,
	}
}

// Конструктор создания объекта капчи в бд
func NewDtoCaptcha(hash string, image []byte, value string, created time.Time, inuse bool) *DtoCaptcha {
	return &DtoCaptcha{
		Hash:    hash,
		Image:   image,
		Value:   value,
		Created: created,
		InUse:   inuse,
	}
}
