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

// Структура для организации хранения отправителя
type ViewSMSSender struct {
	Name string `json:"name" validate:"nonzero,min=1,max=255"` // Название
}

type ApiMetaSMSSender struct {
	Total        int64 `json:"count"`   // Общее число отправителей
	NumOfDeleted int64 `json:"deleted"` // Число удаленных отправителей
	NumOfNew     int64 `json:"new"`     // Число незарегистрированных отправителей
}

type ApiShortSMSSender struct {
	ID         int64  `json:"id" db:"id"`                 // Уникальный идентификатор отправителя
	Name       string `json:"name" db:"name"`             // Название
	Registered bool   `json:"registered" db:"registered"` // Зарегистрирован
}

type ApiMiddleSMSSender struct {
	ID         int64  `json:"id" db:"id"`                 // Уникальный идентификатор отправителя
	Name       string `json:"name" db:"name"`             // Название
	Registered bool   `json:"registered" db:"registered"` // Зарегистрирован
	Deleted    bool   `json:"del" db:"del"`               // Удален
}

type ApiLongSMSSender struct {
	ID         int64  `json:"id" db:"id"`                 // Уникальный идентификатор отправителя
	Name       string `json:"name" db:"name"`             // Название
	Registered bool   `json:"registered" db:"registered"` // Зарегистрирован
	Withdraw   bool   `json:"withdraw" db:"withdraw"`     // Спрятан
	Withdrawn  bool   `json:"withdrawn" db:"withdrawn"`   // Отозван
	Deleted    bool   `json:"del" db:"del"`               // Удален
}

type SMSSenderSearch struct {
	ID         int64  `query:"id" search:"id"`                 // Уникальный идентификатор отправителя
	Name       string `query:"name" search:"name"`             // Название
	Registered bool   `query:"registered" search:"registered"` // Зарегистрирован
}

type DtoSMSSender struct {
	ID         int64     `db:"id"`         // Уникальный идентификатор отправителя
	Unit_ID    int64     `db:"unit_id"`    // Идентификатор объединения
	Name       string    `db:"name"`       // Название
	Created    time.Time `db:"created"`    // Время создания
	Registered bool      `db:"registered"` // Зарегистрирован
	Withdraw   bool      `db:"withdraw"`   // Спрятан
	Withdrawn  bool      `db:"withdrawn"`  // Отозван
	Active     bool      `db:"active"`     // Aктивен
}

// Конструктор создания объекта отправителя в api
func NewApiMetaSMSSender(total int64, numofdeleted int64, numofnew int64) *ApiMetaSMSSender {
	return &ApiMetaSMSSender{
		Total:        total,
		NumOfDeleted: numofdeleted,
		NumOfNew:     numofnew,
	}
}

func NewApiShortSMSSender(id int64, name string, registered bool) *ApiShortSMSSender {
	return &ApiShortSMSSender{
		ID:         id,
		Name:       name,
		Registered: registered,
	}
}

func NewApiMiddleSMSSender(id int64, name string, registered bool, deleted bool) *ApiMiddleSMSSender {
	return &ApiMiddleSMSSender{
		ID:         id,
		Name:       name,
		Registered: registered,
		Deleted:    deleted,
	}
}

func NewApiLongSMSSender(id int64, name string, registered bool, withdraw bool, withdrawn bool, deleted bool) *ApiLongSMSSender {
	return &ApiLongSMSSender{
		ID:         id,
		Name:       name,
		Registered: registered,
		Withdraw:   withdraw,
		Withdrawn:  withdrawn,
		Deleted:    deleted,
	}
}

// Конструктор создания объекта отправителя в бд
func NewDtoSMSSender(id int64, unit_id int64, name string, created time.Time, registered bool,
	withdraw bool, withdrawn bool, active bool) *DtoSMSSender {
	return &DtoSMSSender{
		ID:         id,
		Unit_ID:    unit_id,
		Name:       name,
		Created:    created,
		Registered: registered,
		Withdraw:   withdraw,
		Withdrawn:  withdrawn,
		Active:     active,
	}
}

func (smssender *SMSSenderSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, smssender), nil
}

func (smssender *SMSSenderSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, smssender)
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
	case "name":
		if strings.Contains(invalue, "'") {
			errValue = errors.New("Wrong field value")
			break
		}
		outvalue = "'" + invalue + "'"
	case "registered":
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

func (smssender *SMSSenderSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(smssender)
}

func (smssender *ViewSMSSender) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(smssender, errors, req)
}
