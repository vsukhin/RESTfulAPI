package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

const (
	SERVICE_TYPE_SMS       = "sms"
	SERVICE_TYPE_HLR       = "hlr"
	SERVICE_TYPE_RECOGNIZE = "recognize"
	SERVICE_TYPE_VERIFY    = "verification"
	SERVICE_TYPE_HEADER    = "header"
)

// Структура для организации хранения сервисов
type ViewFacility struct {
	ID int64 `json:"id" validate:"nonzero"` // Уникальный идентификатор сервиса
}

type ViewFacilities []ViewFacility

type ApiShortFacility struct {
	ID          int64  `json:"id" db:"id"`                   // Уникальный идентификатор сервиса
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
}

type ApiLongFacility struct {
	ID          int64  `json:"id" db:"id"`                   // Уникальный идентификатор сервиса
	Name        string `json:"name" db:"name"`               // Название
	Description string `json:"description" db:"description"` // Описание
	Active      bool   `json:"active" db:"active"`           // Активный
}

type ApiFullFacility struct {
	ID              int64  `json:"id" db:"id"`                            // Уникальный идентификатор сервиса
	Category_ID     int    `json:"category" db:"category_id"`             // Идентификатор класса
	Alias           string `json:"alias" db:"alias`                       // Уникальный псевдоним
	Name            string `json:"name" db:"name"`                        // Название
	Description     string `json:"description" db:"description"`          // Описание
	DescriptionSoon string `json:"soonDescription" db:"description_soon"` // Описание доступности
	Active          bool   `json:"active" db:"active"`                    // Активный
	PicNormal_ID    int64  `json:"picNormal" db:"picNormal_id"`           // Идентификатор картинки активного
	PicOver_ID      int64  `json:"picOver" db:"picOver_id"`               // Идентификатор картинки поверх
	PicSoon_ID      int64  `json:"picSoon" db:"picSoon_id"`               // Идентификатор картинки вскоре
	PicDisable_ID   int64  `json:"picDisable" db:"picDisable_id"`         // Идентификатор картинки отключенного
}

type DtoFacility struct {
	ID              int64     `db:"id"`               // Уникальный идентификатор сервиса
	Name            string    `db:"name"`             // Название
	Description     string    `db:"description"`      // Описание
	Created         time.Time `db:"created"`          // Время создания
	Active          bool      `db:"active"`           // Активный
	Category_ID     int       `db:"category_id"`      // Идентификатор класса
	DescriptionSoon string    `db:"description_soon"` // Описание доступности
	PicNormal_ID    int64     `db:"picNormal_id"`     // Идентификатор картинки активного
	PicOver_ID      int64     `db:"picOver_id"`       // Идентификатор картинки поверх
	PicSoon_ID      int64     `db:"picSoon_id"`       // Идентификатор картинки вскоре
	PicDisable_ID   int64     `db:"picDisable_id"`    // Идентификатор картинки отключенного
	Alias           string    `db:"alias`             // Уникальный псевдоним
}

// Конструктор создания объекта сервиса в api
func NewApiShortFacility(id int64, name string, description string) *ApiShortFacility {
	return &ApiShortFacility{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

func NewApiLongFacility(id int64, name string, description string, active bool) *ApiLongFacility {
	return &ApiLongFacility{
		ID:          id,
		Name:        name,
		Description: description,
		Active:      active,
	}
}

func NewApiFullFacility(id int64, category_id int, alias string, name string, description string, descriptionsoon string, active bool,
	picnormal_id int64, picover_id int64, picsoon_id int64, picdisable_id int64) *ApiFullFacility {
	return &ApiFullFacility{
		ID:              id,
		Category_ID:     category_id,
		Alias:           alias,
		Name:            name,
		Description:     description,
		DescriptionSoon: descriptionsoon,
		Active:          active,
		PicNormal_ID:    picnormal_id,
		PicOver_ID:      picover_id,
		PicSoon_ID:      picsoon_id,
		PicDisable_ID:   picdisable_id,
	}
}

// Конструктор создания объекта сервиса в бд
func NewDtoFacility(id int64, name string, description string, created time.Time, active bool, category_id int,
	descriptionsoon string, picnormal_id int64, picover_id int64, picsoon_id int64, picdisable_id int64, alias string) *DtoFacility {
	return &DtoFacility{
		ID:              id,
		Name:            name,
		Description:     description,
		Created:         created,
		Active:          active,
		Category_ID:     category_id,
		DescriptionSoon: descriptionsoon,
		PicNormal_ID:    picnormal_id,
		PicOver_ID:      picover_id,
		PicSoon_ID:      picsoon_id,
		PicDisable_ID:   picdisable_id,
		Alias:           alias,
	}
}

func (facility ViewFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(&facility, errors, req)
}
