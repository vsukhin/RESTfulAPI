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

const (
	MAX_STEP_NUMBER = 3
)

//Структура для организации хранения заказа
type ViewOrder struct {
	Name               string  `json:"name" validate:"nonzero,min=1,max=255"` // Название
	Step               byte    `json:"step" validate:"nonzero"`               // Шаг
	IsAssembled        bool    `json:"completed"`                             // Собран
	IsConfirmed        bool    `json:"moderatorConfirmed"`                    // Утвержден
	Facility_ID        int64   `json:"type" validate:"nonzero"`               // Идентификатор услуги
	Supplier_ID        int64   `json:"supplierId" validate:"nonzero"`         // Идентификатор поставщика
	IsNew              bool    `json:"new"`                                   // Новый
	IsOpen             bool    `json:"open"`                                  // Открыт
	IsCancelled        bool    `json:"cancel"`                                // Отказ
	Reason             string  `json:"cancelDescription" validate:"max=255"`  // Причина отказа
	Proposed_Price     float64 `json:"supplierCost"`                          // Предложенная цена
	IsNewCost          bool    `json:"supplierCostNew"`                       // Новая цена
	IsNewCostConfirmed bool    `json:"customerNewCostConfirmed"`              // Новая цена утверждена
	IsPaid             bool    `json:"paid"`                                  // Оплачен
	IsStarted          bool    `json:"moderatorBegin"`                        // Запущен
	Charged_Fee        float64 `json:"supplierFactualCost"`                   // Фактическая цена
	IsExecuted         bool    `json:"supplierClose"`                         // Выполнен
	IsDocumented       bool    `json:"moderatorDocumentsGotten"`              // Документы имеются
	IsClosed           bool    `json:"moderatorClose"`                        // Закрыт
	IsArchived         bool    `json:"archive"`                               // Архивирован
	IsDeleted          bool    `json:"del"`                                   // Удален
}

type ApiMetaOrder struct {
	Total         int64 `json:"count"`   // Общее число заказов
	NumOfNew      int64 `json:"new"`     // Число новых заказов
	NumOfOpen     int64 `json:"open"`    // Число заказов в работе
	NumOfClosed   int64 `json:"close"`   // Число закрытых заказов
	NumOfArchived int64 `json:"archive"` // Число архивных заказов
	NumOfAlert    int64 `json:"alert"`   // Число заказов с уведомлениями
}

type ApiShortOrder struct {
	ID          int64  `json:"id" db:"id"`                 // Уникальный идентификатор
	IsAssembled bool   `json:"completed" db:"completed"`   // Собран
	Facility_ID int64  `json:"type" db:"type"`             // Идентификатор услуги
	Supplier_ID int64  `json:"supplierId" db:"supplierId"` // Идентификатор поставщика
	IsNew       bool   `json:"new" db:"new"`               // Новый
	IsOpen      bool   `json:"open" db:"open"`             // Открыт
	Name        string `json:"name" db:"name"`             // Название
}

type ApiLongOrder struct {
	ID                 int64   `json:"id"`                       // Уникальный идентификатор
	Name               string  `json:"name"`                     // Название
	Step               byte    `json:"step"`                     // Шаг
	IsAssembled        bool    `json:"completed"`                // Собран
	IsConfirmed        bool    `json:"moderatorConfirmed"`       // Утвержден
	Facility_ID        int64   `json:"type"`                     // Идентификатор услуги
	Supplier_ID        int64   `json:"supplierId"`               // Идентификатор поставщика
	IsNew              bool    `json:"new"`                      // Новый
	IsOpen             bool    `json:"open"`                     // Открыт
	IsCancelled        bool    `json:"cancel"`                   // Отказ
	Reason             string  `json:"cancelDescription"`        // Причина отказа
	Proposed_Price     float64 `json:"supplierCost"`             // Предложенная цена
	IsNewCost          bool    `json:"supplierCostNew"`          // Новая цена
	IsNewCostConfirmed bool    `json:"customerNewCostConfirmed"` // Новая цена утверждена
	IsPaid             bool    `json:"paid"`                     // Оплачен
	IsStarted          bool    `json:"moderatorBegin"`           // Запущен
	Charged_Fee        float64 `json:"supplierFactualCost"`      // Фактическая цена
	IsExecuted         bool    `json:"supplierClose"`            // Выполнен
	IsDocumented       bool    `json:"moderatorDocumentsGotten"` // Документы имеются
	IsClosed           bool    `json:"moderatorClose"`           // Закрыт
	IsArchived         bool    `json:"archive"`                  // Архивирован
	IsDeleted          bool    `json:"del"`                      // Удален
}

type OrderSearch struct {
	ID          int64  `query:"id" search:"o.id"`                        // Уникальный идентификатор
	IsAssembled bool   `query:"completed" search:"coalesce(c.value, 0)"` // Собран
	Facility_ID int64  `query:"type" search:"o.service_id"`              // Идентификатор услуги
	Supplier_ID int64  `query:"supplierId" search:"o.supplier_id"`       // Идентификатор поставщика
	IsNew       bool   `query:"new" search:"coalesce(n.value, 0)"`       // Новый
	IsOpen      bool   `query:"open" search:"coalesce(p.value, 0)"`      // Открыт
	Name        string `query:"name" search:"o.name"`                    // Название
}

type DtoOrder struct {
	ID             int64     `db:"id"`             // Уникальный идентификатор
	Creator_ID     int64     `db:"user_id"`        // Идентификатор создателя
	Unit_ID        int64     `db:"unit_id"`        // Идентификатор объединения
	Supplier_ID    int64     `db:"supplier_id"`    // Идентификатор поставщика
	Facility_ID    int64     `db:"service_id"`     // Идентификатор услуги
	Name           string    `db:"name"`           // Название
	Step           byte      `db:"step"`           // Шаг
	Created        time.Time `db:"created"`        // Время создания
	Proposed_Price float64   `db:"proposed_price"` // Предложенная цена
	Charged_Fee    float64   `db:"charged_fee"`    // Фактическая цена
}

// Конструктор создания объекта заказа в api
func NewApiMetaOrder(total int64, numofnew int64, numofopen int64, numofclosed int64, numofarchived int64,
	numofalert int64) *ApiMetaOrder {
	return &ApiMetaOrder{
		Total:         total,
		NumOfNew:      numofnew,
		NumOfOpen:     numofopen,
		NumOfClosed:   numofclosed,
		NumOfArchived: numofarchived,
		NumOfAlert:    numofalert,
	}
}

func NewApiShortOrder(id int64, isassembled bool, facility_id int64, supplier_id int64, isnew bool, isopen bool,
	name string) *ApiShortOrder {
	return &ApiShortOrder{
		ID:          id,
		IsAssembled: isassembled,
		Facility_ID: facility_id,
		Supplier_ID: supplier_id,
		IsNew:       isnew,
		IsOpen:      isopen,
		Name:        name,
	}
}

func NewApiLongOrder(id int64, name string, step byte, isassembled bool, isconfirmed bool, facility_id int64,
	supplier_id int64, isnew bool, isopen bool, iscancelled bool, reason string, proposed_price float64,
	isnewcost bool, isnewcostconfirmed bool, ispaid bool, isstarted bool, charged_fee float64, isexecuted bool,
	isdocumented bool, isclosed bool, isarchived bool, isdeleted bool) *ApiLongOrder {
	return &ApiLongOrder{
		ID:                 id,
		Name:               name,
		Step:               step,
		IsAssembled:        isassembled,
		IsConfirmed:        isconfirmed,
		Facility_ID:        facility_id,
		Supplier_ID:        supplier_id,
		IsNew:              isnew,
		IsOpen:             isopen,
		IsCancelled:        iscancelled,
		Reason:             reason,
		Proposed_Price:     proposed_price,
		IsNewCost:          isnewcost,
		IsNewCostConfirmed: isnewcostconfirmed,
		IsPaid:             ispaid,
		IsStarted:          isstarted,
		Charged_Fee:        charged_fee,
		IsExecuted:         isexecuted,
		IsDocumented:       isdocumented,
		IsClosed:           isclosed,
		IsArchived:         isarchived,
		IsDeleted:          isdeleted,
	}
}

// Конструктор создания объекта заказа в бд
func NewDtoOrder(id int64, creator_id int64, supplier_id int64, facility_id int64, name string, step byte, created time.Time,
	proposed_price float64, charged_fee float64) *DtoOrder {
	return &DtoOrder{
		ID:             id,
		Creator_ID:     creator_id,
		Supplier_ID:    supplier_id,
		Facility_ID:    facility_id,
		Name:           name,
		Step:           step,
		Created:        created,
		Proposed_Price: proposed_price,
		Charged_Fee:    charged_fee,
	}
}

func (order *OrderSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, order), nil
}

func (order *OrderSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, order)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "type":
		fallthrough
	case "supplierId":
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
	case "completed":
		fallthrough
	case "new":
		fallthrough
	case "open":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Uknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (order *OrderSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(order)
}

func (order *ViewOrder) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(order, errors, req)
}
