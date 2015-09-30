package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Alignment byte

const (
	ALIGNMENT_LEFT Alignment = iota + 1
	ALIGNMNET_RIGHT
	ALIGNMENT_CENTER
)

const (
	COLUMN_TYPE_DEFAULT                               = 0
	COLUMN_TYPE_MOBILE_PHONE                          = 1
	COLUMN_TYPE_SMS                                   = 2
	COLUMN_TYPE_SMS_SENDER                            = 3
	COLUMN_TYPE_BIRTHDAY                              = 4
	COLUMN_TYPE_SOURCE_ADDRESS                        = 5
	COLUMN_TYPE_SOURCE_PHONE                          = 6
	COLUMN_TYPE_SOURCE_FIO                            = 7
	COLUMN_TYPE_SOURCE_EMAIL                          = 8
	COLUMN_TYPE_SOURCE_DATE                           = 9
	COLUMN_TYPE_SOURCE_AUTOMOBILE                     = 10
	COLUMN_TYPE_ANSWER_POSTADDRESS_RESULT             = 11
	COLUMN_TYPE_ANSWER_PHONE_RESULT                   = 12
	COLUMN_TYPE_ANSWER_FULLNAME_RESULT                = 13
	COLUMN_TYPE_ANSWER_EMAIL_RESULT                   = 14
	COLUMN_TYPE_ANSWER_DATE_RESULT                    = 15
	COLUMN_TYPE_ANSWER_VEHICLE_RESULT                 = 16
	COLUMN_TYPE_ANSWER_POSTADDRESS_POSTALCODE         = 17
	COLUMN_TYPE_ANSWER_POSTADDRESS_COUNTRY            = 18
	COLUMN_TYPE_ANSWER_POSTADDRESS_REGIONTYPE         = 19
	COLUMN_TYPE_ANSWER_POSTADDRESS_REGIONTYPEFULL     = 20
	COLUMN_TYPE_ANSWER_POSTADDRESS_REGION             = 21
	COLUMN_TYPE_ANSWER_POSTADDRESS_AREATYPE           = 22
	COLUMN_TYPE_ANSWER_POSTADDRESS_AREATYPEFULL       = 23
	COLUMN_TYPE_ANSWER_POSTADDRESS_AREA               = 24
	COLUMN_TYPE_ANSWER_POSTADDRESS_CITYTYPE           = 25
	COLUMN_TYPE_ANSWER_POSTADDRESS_CITYTYPEFULL       = 26
	COLUMN_TYPE_ANSWER_POSTADDRESS_CITY               = 27
	COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENTTYPE     = 28
	COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENTTYPEFULL = 29
	COLUMN_TYPE_ANSWER_POSTADDRESS_SETTLEMENT         = 30
	COLUMN_TYPE_ANSWER_POSTADDRESS_STREETTYPE         = 31
	COLUMN_TYPE_ANSWER_POSTADDRESS_STREETTYPEFULL     = 32
	COLUMN_TYPE_ANSWER_POSTADDRESS_STREET             = 33
	COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSETYPE          = 34
	COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSETYPEFULL      = 35
	COLUMN_TYPE_ANSWER_POSTADDRESS_HOUSE              = 36
	COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCKTYPE          = 37
	COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCKTYPEFULL      = 38
	COLUMN_TYPE_ANSWER_POSTADDRESS_BLOCK              = 39
	COLUMN_TYPE_ANSWER_POSTADDRESS_FLATTYPE           = 40
	COLUMN_TYPE_ANSWER_POSTADDRESS_FLAT               = 41
	COLUMN_TYPE_ANSWER_POSTADDRESS_FLATAREA           = 42
	COLUMN_TYPE_ANSWER_POSTADDRESS_SQUAREMETERPRICE   = 43
	COLUMN_TYPE_ANSWER_POSTADDRESS_FLATPRICE          = 44
	COLUMN_TYPE_ANSWER_POSTADDRESS_POSTALBOX          = 45
	COLUMN_TYPE_ANSWER_POSTADDRESS_FIASID             = 46
	COLUMN_TYPE_ANSWER_POSTADDRESS_KLADRID            = 47
	COLUMN_TYPE_ANSWER_POSTADDRESS_OKATO              = 48
	COLUMN_TYPE_ANSWER_POSTADDRESS_OKTMO              = 49
	COLUMN_TYPE_ANSWER_POSTADDRESS_TAXOFFICE          = 50
	COLUMN_TYPE_ANSWER_POSTADDRESS_TAXOFFICELEGAL     = 51
	COLUMN_TYPE_ANSWER_POSTADDRESS_TIMEZONE           = 52
	COLUMN_TYPE_ANSWER_POSTADDRESS_GEOLAT             = 53
	COLUMN_TYPE_ANSWER_POSTADDRESS_GEOLON             = 54
	COLUMN_TYPE_ANSWER_POSTADDRESS_QCGEO              = 55
	COLUMN_TYPE_ANSWER_POSTADDRESS_QCCOMPLETE         = 56
	COLUMN_TYPE_ANSWER_POSTADDRESS_QCHOUSE            = 57
	COLUMN_TYPE_ANSWER_POSTADDRESS_QUALITYCODE        = 58
	COLUMN_TYPE_ANSWER_POSTADDRESS_UNPARSEDPARTS      = 59
	COLUMN_TYPE_ANSWER_PHONE_TYPE                     = 60
	COLUMN_TYPE_ANSWER_PHONE_COUNTRYCODE              = 61
	COLUMN_TYPE_ANSWER_PHONE_CITYCODE                 = 62
	COLUMN_TYPE_ANSWER_PHONE_NUMBER                   = 63
	COLUMN_TYPE_ANSWER_PHONE_EXTENSION                = 64
	COLUMN_TYPE_ANSWER_PHONE_PROVIDER                 = 65
	COLUMN_TYPE_ANSWER_PHONE_REGION                   = 66
	COLUMN_TYPE_ANSWER_PHONE_TIMEZONE                 = 67
	COLUMN_TYPE_ANSWER_PHONE_QCCONFLICT               = 68
	COLUMN_TYPE_ANSWER_PHONE_QUALITYCODE              = 69
	COLUMN_TYPE_ANSWER_FULLNAME_SURNAME               = 70
	COLUMN_TYPE_ANSWER_FULLNAME_NAME                  = 71
	COLUMN_TYPE_ANSWER_FULLNAME_PATRONYMIC            = 72
	COLUMN_TYPE_ANSWER_FULLNAME_GENDER                = 73
	COLUMN_TYPE_ANSWER_FULLNAME_QUALITYCODE           = 74
	COLUMN_TYPE_ANSWER_VEHICLE_BRAND                  = 75
	COLUMN_TYPE_ANSWER_VEHICLE_MODEL                  = 76
	COLUMN_TYPE_ANSWER_VEHICLE_QUALITYCODE            = 77
	COLUMN_TYPE_PRICELIST_NAME                        = 78
	COLUMN_TYPE_PRICELIST_PRICE                       = 79
	COLUMN_TYPE_PRICELIST_DISCOUNT                    = 80
	COLUMN_TYPE_PRICELIST_MOBILEOPERATOR              = 81
	COLUMN_TYPE_PRICELIST_RANGE                       = 82
	COLUMN_TYPE_PRICELIST_ID                          = 83
	COLUMN_TYPE_SOURCE_PASSPORT                       = 84
	COLUMN_TYPE_ANSWER_PASSPORT_CODE                  = 85
	COLUMN_TYPE_ANSWER_PASSPORT_NUMBER                = 86
	COLUMN_TYPE_PRICELIST_FEE_ONCE                    = 87
	COLUMN_TYPE_PRICELIST_FEE_MONTHLY                 = 88
	COLUMN_TYPE_ANSWER_PASSPORT_ISSUEDATE             = 89
	COLUMN_TYPE_PASSPORT                              = 90
	COLUMN_TYPE_URL                                   = 91
	COLUMN_TYPE_ANSWER_PASSPORT_UNITCODE              = 92
	COLUMN_TYPE_ANSWER_PASSPORT_BIRTHPLACE            = 93
	COLUMN_TYPE_ANSWER_PASSPORT_ISSUEDBY              = 94
	COLUMN_TYPE_ANSWER_PASSPORT_QUALITYCODE           = 95
)

// Структура для организации хранения типов колонок
type ApiColumnType struct {
	ID               int    `json:"id" db:"id"`                       // Уникальный идентификатор типа колонки
	Name             string `json:"name" db:"name"`                   // Название
	Position         int64  `json:"position" db:"position"`           // Позиция сортировки для UI
	Description      string `json:"description" db:"description"`     // Описание
	Required         bool   `json:"notNull" db:"notNull"`             // Обязательность к заполнению
	Regexp           string `json:"regexp" db:"regexp"`               // Регулярное выражение для проверки
	HorAlignmentHead string `json:"alignmentHead" db:"alignmentHead"` // Горизонтальное выравнивание заголовка
	HorAlignmentBody string `json:"alignmentBody" db:"alignmentBody"` // Горизонтальное выравнивание содержимого
}

type ColumnTypeSearch struct {
	ID               int    `query:"id" search:"id"`                                       // Уникальный идентификатор типа колонки
	Name             string `query:"name" search:"name" group:"name"`                      // Название
	Position         int64  `query:"position" search:"position"`                           // Позиция сортировки для UI
	Description      string `query:"description" search:"description" group:"description"` // Описание
	Required         bool   `query:"notNull" search:"required"`                            // Обязательность к заполнению
	Regexp           string `query:"regexp" search:"regexp"`                               // Регулярное выражение для проверки
	HorAlignmentHead string `query:"alignmentHead" search:"horAlignmentHead"`              // Горизонтальное выравнивание заголовка
	HorAlignmentBody string `query:"alignmentBody" search:"horAlignmentBody"`              // Горизонтальное выравнивание содержимого
	Private          bool   `query:"private" search:"private"`                             // Если =true (1) - тип приватный, выгружается только при явном указании, Если =false (0) - тип публичный, выгружается по умолчанию
}

type DtoColumnType struct {
	ID               int       `db:"id"`               // Уникальный идентификатор типа колонки
	Name             string    `db:"name"`             // Название
	Position         int64     `db:"position"`         // Позиция сортировки для UI
	Description      string    `db:"description"`      // Описание
	Required         bool      `db:"required"`         // Обязательность к заполнению
	Regexp           string    `db:"regexp"`           // Регулярное выражение для проверки
	HorAlignmentHead Alignment `db:"horAlignmentHead"` // Горизонтальное выравнивание заголовка
	HorAlignmentBody Alignment `db:"horAlignmentBody"` // Горизонтальное выравнивание содержимого
	Created          time.Time `db:"created"`          // Время создания
	Active           bool      `db:"active"`           // Активная
	Private          bool      `db:"private"`          // Если =true (1) - тип приватный, выгружается только при явном указании, Если =false (0) - тип публичный, выгружается по умолчанию
}

// Конструктор создания объекта типа колонки в api
func NewApiColumnType(id int, name string, position int64, description string, required bool, regexp string,
	horalignmenthead string, horalignmentbody string) *ApiColumnType {
	return &ApiColumnType{
		ID:               id,
		Name:             name,
		Position:         position,
		Description:      description,
		Required:         required,
		Regexp:           regexp,
		HorAlignmentHead: horalignmenthead,
		HorAlignmentBody: horalignmentbody,
	}
}

// Конструктор создания объекта типа колонки в бд
func NewDtoColumnType(id int, name string, position int64, description string, required bool, regexp string,
	horalignmenthead Alignment, horalignmentbody Alignment, created time.Time, active bool, private bool) *DtoColumnType {
	return &DtoColumnType{
		ID:               id,
		Name:             name,
		Position:         position,
		Description:      description,
		Required:         required,
		Regexp:           regexp,
		HorAlignmentHead: horalignmenthead,
		HorAlignmentBody: horalignmentbody,
		Created:          created,
		Active:           active,
		Private:          private,
	}
}

func (columntype *ColumnTypeSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, columntype), nil
}

func (columntype *ColumnTypeSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, columntype)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "position":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		fallthrough
	case "description":
		fallthrough
	case "regexp":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
		fallthrough
	case "alignmentHead":
		fallthrough
	case "alignmentBody":
		if strings.ToLower(invalue) == "left" {
			outvalue = "1"
		} else if strings.ToLower(invalue) == "center" {
			outvalue = "2"
		} else if strings.ToLower(invalue) == "right" {
			outvalue = "3"
		} else {
			errValue = errors.New("Wrong value")
		}
	case "notNull":
		fallthrough
	case "private":
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

func (columntype *ColumnTypeSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(columntype)
}
