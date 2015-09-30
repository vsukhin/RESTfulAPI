package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения hlr сервиса заказа
type ViewHLRFacility struct {
	EstimatedNumbersShipments int                              `json:"estimatedNumbersShipments" validate:"min=0"` // Предполагаемое количество отправок HLR запросов
	EstimatedOperators        []ViewApiMobileOperatorOperation `json:"estimatedOperators"`                         // Прогноз распределения отправлений
	DeliveryDataId            int64                            `json:"deliveryDataId"`                             // Уникальный идентификатор таблицы
	DeliveryDataDelete        bool                             `json:"deliveryDataDelete"`                         // Удалить указанную таблицу
	MessageToInColumnId       int64                            `json:"messageToInColumnId"`                        // Идентификатор колонки в таблице
}

type ApiHLRFacility struct {
	EstimatedNumbersShipments int                              `json:"estimatedNumbersShipments" db:"estimatedNumbersShipments"` // Предполагаемое количество отправок HLR запросов
	EstimatedOperators        []ViewApiMobileOperatorOperation `json:"estimatedOperators,omitempty" db:"-"`                      // Прогноз распределения отправлений
	DeliveryDataId            int64                            `json:"deliveryDataId" db:"deliveryDataId"`                       // Уникальный идентификатор таблицы
	DeliveryDataDelete        bool                             `json:"deliveryDataDelete" db:"deliveryDataDelete"`               // Удалить указанную таблицу
	MessageToInColumnId       int64                            `json:"messageToInColumnId" db:"messageToInColumnId"`             // Идентификатор колонки в таблице
	Cost                      float64                          `json:"cost" db:"cost"`                                           // Сумма заказа исходя из расчётных показателей заказа
	CostFactual               float64                          `json:"costFactual" db:"costFactual"`                             // Текущая стоимость заказа
	ResultTables              []ApiResultTable                 `json:"resultTables,omitempty" db:"-"`                            // Таблицы результатов
	WorkTables                []ApiWorkTable                   `json:"workTables,omitempty" db:"-"`                              // Рабочие таблицы
}

type DtoHLRFacility struct {
	Order_ID                  int64                        `db:"order_id"`                  // Идентификатор заказа
	EstimatedNumbersShipments int                          `db:"estimatedNumbersShipments"` // Предполагаемое количество отправок HLR запросов
	EstimatedOperators        []DtoMobileOperatorOperation `db:"-"`                         // Прогноз распределения отправлений
	DeliveryDataId            int64                        `db:"deliveryDataId"`            // Уникальный идентификатор таблицы
	DeliveryDataDelete        bool                         `db:"deliveryDataDelete"`        // Удалить указанную таблицу
	MessageToInColumnId       int64                        `db:"messageToInColumnId"`       // Идентификатор колонки в таблице
	Cost                      float64                      `db:"cost"`                      // Сумма заказа исходя из расчётных показателей заказа
	CostFactual               float64                      `db:"costFactual"`               // Текущая стоимость заказа
	ResultTables              []DtoResultTable             `db:"-"`                         // Таблицы результатов
	WorkTables                []DtoWorkTable               `db:"-"`                         // Рабочие таблицы
}

// Конструктор создания объекта hlr сервиса заказа в api
func NewApiHLRFacility(estimatedNumbersShipments int, estimatedOperators []ViewApiMobileOperatorOperation,
	deliveryDataId int64, deliveryDataDelete bool, messageToInColumnId int64, cost float64, costFactual float64,
	resultTables []ApiResultTable, workTables []ApiWorkTable) *ApiHLRFacility {
	return &ApiHLRFacility{
		EstimatedNumbersShipments: estimatedNumbersShipments,
		EstimatedOperators:        estimatedOperators,
		DeliveryDataId:            deliveryDataId,
		DeliveryDataDelete:        deliveryDataDelete,
		MessageToInColumnId:       messageToInColumnId,
		Cost:                      cost,
		CostFactual:               costFactual,
		ResultTables:              resultTables,
		WorkTables:                workTables,
	}
}

// Конструктор создания объекта hlr сервиса заказа в бд
func NewDtoHLRFacility(order_id int64, estimatedNumbersShipments int, estimatedOperators []DtoMobileOperatorOperation,
	deliveryDataId int64, deliveryDataDelete bool, messageToInColumnId int64, cost float64, costFactual float64,
	resultTables []DtoResultTable, workTables []DtoWorkTable) *DtoHLRFacility {
	return &DtoHLRFacility{
		Order_ID:                  order_id,
		EstimatedNumbersShipments: estimatedNumbersShipments,
		EstimatedOperators:        estimatedOperators,
		DeliveryDataId:            deliveryDataId,
		DeliveryDataDelete:        deliveryDataDelete,
		MessageToInColumnId:       messageToInColumnId,
		Cost:                      cost,
		CostFactual:               costFactual,
		ResultTables:              resultTables,
		WorkTables:                workTables,
	}
}

func (facility *ViewHLRFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, operator := range facility.EstimatedOperators {
		errors = Validate(&operator, errors, req)
	}
	return Validate(facility, errors, req)
}
