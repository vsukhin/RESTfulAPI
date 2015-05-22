package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

type DeliveryType byte

const (
	TYPE_DELIVERY_ONCE DeliveryType = iota + 1
	TYPE_DELIVERY_SCHEDULED
	TYPE_DELIVERY_EVENTTRIGGERED
)

const (
	TYPE_DELIVERY_ONCE_VALUE           = "onetime"
	TYPE_DELIVERY_SCHEDULED_VALUE      = "periodic"
	TYPE_DELIVERY_EVENTTRIGGERED_VALUE = "eventdispatch"
)

// Структура для организации хранения sms сервиса заказа
type ViewSMSFacility struct {
	EstimatedNumbersShipments   int                              `json:"estimatedNumbersShipments" validate:"min=0"`   // Предполагаемое количество отправок SMS
	EstimatedMessageInCyrillic  bool                             `json:"estimatedMessageInCyrillic"`                   // Содержит кириллицу
	EstimatedNumberCharacters   int                              `json:"estimatedNumberCharacters" validate:"min=0"`   // Количество символов в одном сообщении
	EstimatedNumberSmsInMessage int                              `json:"estimatedNumberSmsInMessage" validate:"min=0"` // Количество SMS используемых для отправки одного сообщения
	EstimatedOperators          []ViewApiMobileOperatorOperation `json:"estimatedOperators"`                           // Прогноз распределения отправлений
	DeliveryType                string                           `json:"deliveryType" validate:"nonzero"`              // Константа периодичности оказания услуги
	DeliveryTime                bool                             `json:"deliveryTime"`                                 // Рассылка в период времени
	DeliveryTimeStart           time.Time                        `json:"deliveryTimeStart"`                            // Дата и время начала рассылки
	DeliveryTimeEnd             time.Time                        `json:"deliveryTimeEnd"`                              // Дата и время прекращения рассылки
	DeliveryBaseTime            time.Time                        `json:"deliveryBaseTime"`                             // Время для рассылок
	DeliveryDataId              int64                            `json:"deliveryDataId" validate:"nonzero"`            // Уникальный идентификатор таблицы
	DeliveryDataDelete          bool                             `json:"deliveryDataDelete"`                           // Удалить указанную таблицу
	MessageFromId               int64                            `json:"messageFromId"`                                // Идентификатор отправителя
	MessageFromInColumnId       int64                            `json:"messageFromInColumnId"`                        // Идентификатор отправителя сообщений
	MessageToInColumnId         int64                            `json:"messageToInColumnId"`                          // Идентификатор колонки в таблице
	MessageBody                 string                           `json:"messageBody"`                                  // Текст рассылки
	MessageBodyInColumnId       int64                            `json:"messageBodyInColumnId"`                        // Идентификатор рассылаемого сообщения
	TimeCorrection              bool                             `json:"timeCorrection"`                               // Время рассылки корректируется
}

type ApiSMSFacility struct {
	EstimatedNumbersShipments   int                              `json:"estimatedNumbersShipments" db:"estimatedNumbersShipments"`     // Предполагаемое количество отправок SMS
	EstimatedMessageInCyrillic  bool                             `json:"estimatedMessageInCyrillic" db:"estimatedMessageInCyrillic"`   // Содержит кириллицу
	EstimatedNumberCharacters   int                              `json:"estimatedNumberCharacters" db:"estimatedNumberCharacters"`     // Количество символов в одном сообщении
	EstimatedNumberSmsInMessage int                              `json:"estimatedNumberSmsInMessage" db:"estimatedNumberSmsInMessage"` // Количество SMS используемых для отправки одного сообщения
	EstimatedOperators          []ViewApiMobileOperatorOperation `json:"estimatedOperators,omitempty" db:"-"`                          // Прогноз распределения отправлений
	DeliveryType                string                           `json:"deliveryType" db:"deliveryType"`                               // Константа периодичности оказания услуги
	DeliveryTime                bool                             `json:"deliveryTime" db:"deliveryTime"`                               // Рассылка в период времени
	DeliveryTimeStart           time.Time                        `json:"deliveryTimeStart" db:"deliveryTimeStart"`                     // Дата и время начала рассылки
	DeliveryTimeEnd             time.Time                        `json:"deliveryTimeEnd" db:"deliveryTimeEnd"`                         // Дата и время прекращения рассылки
	DeliveryBaseTime            time.Time                        `json:"deliveryBaseTime" db:"deliveryBaseTime"`                       // Время для рассылок
	DeliveryDataId              int64                            `json:"deliveryDataId" db:"deliveryDataId"`                           // Уникальный идентификатор таблицы
	DeliveryDataDelete          bool                             `json:"deliveryDataDelete" db:"deliveryDataDelete"`                   // Удалить указанную таблицу
	MessageFromId               int64                            `json:"messageFromId" db:"messageFromId"`                             // Идентификатор отправителя
	MessageFromInColumnId       int64                            `json:"messageFromInColumnId" db:"messageFromInColumnId"`             // Идентификатор отправителя сообщений
	MessageToInColumnId         int64                            `json:"messageToInColumnId" db:"messageToInColumnId"`                 // Идентификатор колонки в таблице
	MessageBody                 string                           `json:"messageBody" db:"messageBody"`                                 // Текст рассылки
	MessageBodyInColumnId       int64                            `json:"messageBodyInColumnId" db:"messageBodyInColumnId"`             // Идентификатор рассылаемого сообщения
	TimeCorrection              bool                             `json:"timeCorrection" db:"timeCorrection"`                           // Время рассылки корректируется
	Cost                        float64                          `json:"cost" db:"cost"`                                               // Сумма заказа исходя из расчётных показателей заказа
	CostFactual                 float64                          `json:"costFactual" db:"costFactual"`                                 // Текущая стоимость заказа
	ResultTables                []ApiResultTable                 `json:"resultTables,omitempty" db:"-"`                                // Таблицы результатов
	WorkTables                  []ApiWorkTable                   `json:"workTables,omitempty" db:"-"`                                  // Рабочие таблицы
}

type DtoSMSFacility struct {
	Order_ID                    int64                        `db:"order_id"`                    // Идентификатор заказа
	EstimatedNumbersShipments   int                          `db:"estimatedNumbersShipments"`   // Предполагаемое количество отправок SMS
	EstimatedMessageInCyrillic  bool                         `db:"estimatedMessageInCyrillic"`  // Содержит кириллицу
	EstimatedNumberCharacters   int                          `db:"estimatedNumberCharacters"`   // Количество символов в одном сообщении
	EstimatedNumberSmsInMessage int                          `db:"estimatedNumberSmsInMessage"` // Количество SMS используемых для отправки одного сообщения
	EstimatedOperators          []DtoMobileOperatorOperation `db:"-"`                           // Прогноз распределения отправлений
	DeliveryType                DeliveryType                 `db:"deliveryType"`                // Константа периодичности оказания услуги
	DeliveryTime                bool                         `db:"deliveryTime"`                // Рассылка в период времени
	DeliveryTimeStart           time.Time                    `db:"deliveryTimeStart"`           // Дата и время начала рассылки
	DeliveryTimeEnd             time.Time                    `db:"deliveryTimeEnd"`             // Дата и время прекращения рассылки
	DeliveryBaseTime            time.Time                    `db:"deliveryBaseTime"`            // Время для рассылок
	DeliveryDataId              int64                        `db:"deliveryDataId"`              // Уникальный идентификатор таблицы
	DeliveryDataDelete          bool                         `db:"deliveryDataDelete"`          // Удалить указанную таблицу
	MessageFromId               int64                        `db:"messageFromId"`               // Идентификатор отправителя
	MessageFromInColumnId       int64                        `db:"messageFromInColumnId"`       // Идентификатор отправителя сообщений
	MessageToInColumnId         int64                        `db:"messageToInColumnId"`         // Идентификатор колонки в таблице
	MessageBody                 string                       `db:"messageBody"`                 // Текст рассылки
	MessageBodyInColumnId       int64                        `db:"messageBodyInColumnId"`       // Идентификатор рассылаемого сообщения
	TimeCorrection              bool                         `db:"timeCorrection"`              // Время рассылки корректируется
	Cost                        float64                      `db:"cost"`                        // Сумма заказа исходя из расчётных показателей заказа
	CostFactual                 float64                      `db:"costFactual"`                 // Текущая стоимость заказа
	ResultTables                []DtoResultTable             `db:"-"`                           // Таблицы результатов
	WorkTables                  []DtoWorkTable               `db:"-"`                           // Рабочие таблицы
}

// Конструктор создания объекта sms сервиса заказа в api
func NewApiSMSFacility(estimatedNumbersShipments int, estimatedMessageInCyrillic bool,
	estimatedNumberCharacters int, estimatedNumberSmsInMessage int, estimatedOperators []ViewApiMobileOperatorOperation,
	deliveryType string, deliveryTime bool, deliveryTimeStart time.Time, deliveryTimeEnd time.Time, deliveryBaseTime time.Time,
	deliveryDataId int64, deliveryDataDelete bool, messageFromId int64, messageFromInColumnId int64, messageToInColumnId int64,
	messageBody string, messageBodyInColumnId int64, timeCorrection bool, cost float64, costFactual float64,
	resultTables []ApiResultTable, workTables []ApiWorkTable) *ApiSMSFacility {
	return &ApiSMSFacility{
		EstimatedNumbersShipments:   estimatedNumbersShipments,
		EstimatedMessageInCyrillic:  estimatedMessageInCyrillic,
		EstimatedNumberCharacters:   estimatedNumberCharacters,
		EstimatedNumberSmsInMessage: estimatedNumberSmsInMessage,
		EstimatedOperators:          estimatedOperators,
		DeliveryType:                deliveryType,
		DeliveryTime:                deliveryTime,
		DeliveryTimeStart:           deliveryTimeStart,
		DeliveryTimeEnd:             deliveryTimeEnd,
		DeliveryBaseTime:            deliveryBaseTime,
		DeliveryDataId:              deliveryDataId,
		DeliveryDataDelete:          deliveryDataDelete,
		MessageFromId:               messageFromId,
		MessageFromInColumnId:       messageFromInColumnId,
		MessageToInColumnId:         messageToInColumnId,
		MessageBody:                 messageBody,
		MessageBodyInColumnId:       messageBodyInColumnId,
		TimeCorrection:              timeCorrection,
		Cost:                        cost,
		CostFactual:                 costFactual,
		ResultTables:                resultTables,
		WorkTables:                  workTables,
	}
}

// Конструктор создания объекта sms сервиса заказа в бд
func NewDtoSMSFacility(order_id int64, estimatedNumbersShipments int, estimatedMessageInCyrillic bool,
	estimatedNumberCharacters int, estimatedNumberSmsInMessage int, estimatedOperators []DtoMobileOperatorOperation,
	deliveryType DeliveryType, deliveryTime bool, deliveryTimeStart time.Time, deliveryTimeEnd time.Time, deliveryBaseTime time.Time,
	deliveryDataId int64, deliveryDataDelete bool, messageFromId int64, messageFromInColumnId int64, messageToInColumnId int64,
	messageBody string, messageBodyInColumnId int64, timeCorrection bool, cost float64, costFactual float64,
	resultTables []DtoResultTable, workTables []DtoWorkTable) *DtoSMSFacility {
	return &DtoSMSFacility{
		Order_ID:                    order_id,
		EstimatedNumbersShipments:   estimatedNumbersShipments,
		EstimatedMessageInCyrillic:  estimatedMessageInCyrillic,
		EstimatedNumberCharacters:   estimatedNumberCharacters,
		EstimatedNumberSmsInMessage: estimatedNumberSmsInMessage,
		EstimatedOperators:          estimatedOperators,
		DeliveryType:                deliveryType,
		DeliveryTime:                deliveryTime,
		DeliveryTimeStart:           deliveryTimeStart,
		DeliveryTimeEnd:             deliveryTimeEnd,
		DeliveryBaseTime:            deliveryBaseTime,
		DeliveryDataId:              deliveryDataId,
		DeliveryDataDelete:          deliveryDataDelete,
		MessageFromId:               messageFromId,
		MessageFromInColumnId:       messageFromInColumnId,
		MessageToInColumnId:         messageToInColumnId,
		MessageBody:                 messageBody,
		MessageBodyInColumnId:       messageBodyInColumnId,
		TimeCorrection:              timeCorrection,
		Cost:                        cost,
		CostFactual:                 costFactual,
		ResultTables:                resultTables,
		WorkTables:                  workTables,
	}
}

func (facility *ViewSMSFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, operator := range facility.EstimatedOperators {
		errors = Validate(&operator, errors, req)
	}
	return Validate(facility, errors, req)
}
