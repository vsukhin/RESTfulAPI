package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Структура для организации хранения класса компании
type ApiCompanyClass struct {
	ID        int    `json:"id" db:"id"`                     // Уникальный идентификатор класса компании
	FullName  string `json:"nameFull" db:"nameFull"`         // Полное название
	ShortName string `json:"nameShort" db:"nameShort"`       // Короткое название
	Format    string `json:"format" db:"format"`             // Регулярное выражение для проверки
	Required  bool   `json:"required" db:"required"`         // Обязательность к заполнению
	Visible   bool   `json:"outward" db:"outward"`           // Всегда видимый
	Multiple  bool   `json:"multiplicity" db:"multiplicity"` // Множественный
	Position  int    `json:"position" db:"position"`         // Позиция
	Deleted   bool   `json:"del" db:"del"`                   // Удален
}

type CompanyClassSearch struct {
	ID        int    `query:"id" search:"id"`                                 // Уникальный идентификатор классификатора компании
	FullName  string `query:"nameFull" search:"fullname" group:"fullname"`    // Полное название
	ShortName string `query:"nameShort" search:"shortname" group:"shortname"` // Короткое название
	Format    string `query:"format" search:"format"`                         // Регулярное выражение для проверки
	Required  bool   `query:"required" search:"required"`                     // Обязательность к заполнению
	Visible   bool   `query:"outward" search:"visible"`                       // Всегда видимый
	Multiple  bool   `query:"multiplicity" search:"multiple"`                 // Множественный
	Position  int    `query:"position" search:"position"`                     // Позиция
	Deleted   bool   `query:"del" search:"(not active)"`                      // Удален
}

type DtoCompanyClass struct {
	ID        int       `db:"id"`        // Уникальный идентификатор класса компании
	FullName  string    `db:"fullname"`  // Полное название
	ShortName string    `db:"shortname"` // Короткое название
	Format    string    `db:"format"`    // Регулярное выражение для проверки
	Required  bool      `db:"required"`  // Обязательность к заполнению
	Visible   bool      `db:"visible"`   // Всегда видимый
	Multiple  bool      `db:"multiple"`  // Множественный
	Position  int       `db:"position"`  // Позиция
	Created   time.Time `db:"created"`   // Время создания
	Active    bool      `db:"active"`    // Aктивен
}

// Конструктор создания объекта класса компании в api
func NewApiCompanyClass(id int, fullname string, shortname string, format string, required bool, visible bool, multiple bool,
	position int, deleted bool) *ApiCompanyClass {
	return &ApiCompanyClass{
		ID:        id,
		FullName:  fullname,
		ShortName: shortname,
		Format:    format,
		Required:  required,
		Visible:   visible,
		Multiple:  multiple,
		Position:  position,
		Deleted:   deleted,
	}
}

// Конструктор создания объекта класса компании в бд
func NewDtoCompanyClass(id int, fullname string, shortname string, format string, required bool, visible bool, multiple bool,
	position int, created time.Time, active bool) *DtoCompanyClass {
	return &DtoCompanyClass{
		ID:        id,
		FullName:  fullname,
		ShortName: shortname,
		Format:    format,
		Required:  required,
		Visible:   visible,
		Multiple:  multiple,
		Position:  position,
		Created:   created,
		Active:    active,
	}
}

func (companyclass *CompanyClassSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, companyclass), nil
}

func (companyclass *CompanyClassSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, companyclass)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "position":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "nameFull":
		fallthrough
	case "nameShort":
		fallthrough
	case "format":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "required":
		fallthrough
	case "outward":
		fallthrough
	case "multiplicity":
		fallthrough
	case "del":
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

func (companyclass *CompanyClassSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(companyclass)
}
