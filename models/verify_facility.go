package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения сервиса верификации данных заказа
type ViewVerifyFacility struct {
	EstimatedNumbersRecords int              `json:"estimatedNumbersRecords" validate:"min=0"` // Предполагаемое количество проверяемых записей
	TablesDataId            int64            `json:"tablesDataId" validate:"nonzero"`          // Уникальный идентификатор таблицы
	TablesDataDelete        bool             `json:"tablesDataDelete"`                         // Удалить указанную таблицу
	DataColumns             []ViewDataColumn `json:"tablesDataColumns"`                        // Массив проверяемых колонок
}

type ApiVerifyFacility struct {
	EstimatedNumbersRecords int              `json:"estimatedNumbersRecords" db:"estimatedNumbersRecords"` // Предполагаемое количество проверяемых записей
	TablesDataId            int64            `json:"tablesDataId" db:"tablesDataId"`                       // Уникальный идентификатор таблицы
	TablesDataDelete        bool             `json:"tablesDataDelete" db:"tablesDataDelete"`               // Удалить указанную таблицу
	DataColumns             []ApiDataColumn  `json:"tablesDataColumns,omitempty" db:"-"`                   // Массив проверяемых колонок
	Cost                    float64          `json:"cost" db:"cost"`                                       // Сумма заказа исходя из расчётных показателей заказа
	CostFactual             float64          `json:"costFactual" db:"costFactual"`                         // Текущая стоимость заказа
	ResultTables            []ApiResultTable `json:"resultTables,omitempty" db:"-"`                        // Таблицы результатов
	WorkTables              []ApiWorkTable   `json:"workTables,omitempty" db:"-"`                          // Рабочие таблицы
}

type DtoVerifyFacility struct {
	Order_ID                int64            `db:"order_id"`                // Идентификатор заказа
	EstimatedNumbersRecords int              `db:"estimatedNumbersRecords"` // Предполагаемое количество проверяемых записей
	TablesDataId            int64            `db:"tablesDataId"`            // Уникальный идентификатор таблицы
	TablesDataDelete        bool             `db:"tablesDataDelete"`        // Удалить указанную таблицу
	DataColumns             []DtoDataColumn  `db:"-"`                       // Массив проверяемых колонок
	Cost                    float64          `db:"cost"`                    // Сумма заказа исходя из расчётных показателей заказа
	CostFactual             float64          `db:"costFactual"`             // Текущая стоимость заказа
	ResultTables            []DtoResultTable `db:"-"`                       // Таблицы результатов
	WorkTables              []DtoWorkTable   `db:"-"`                       // Рабочие таблицы
}

// Конструктор создания объекта сервиса верификации данных заказа в api
func NewApiVerifyFacility(estimatedNumbersRecords int, tablesDataId int64, tablesDataDelete bool, dataColumns []ApiDataColumn,
	cost float64, costFactual float64, resultTables []ApiResultTable, workTables []ApiWorkTable) *ApiVerifyFacility {
	return &ApiVerifyFacility{
		EstimatedNumbersRecords: estimatedNumbersRecords,
		TablesDataId:            tablesDataId,
		TablesDataDelete:        tablesDataDelete,
		DataColumns:             dataColumns,
		Cost:                    cost,
		CostFactual:             costFactual,
		ResultTables:            resultTables,
		WorkTables:              workTables,
	}
}

// Конструктор создания объекта сервиса верификации данных заказа в бд
func NewDtoVerifyFacility(order_id int64, estimatedNumbersRecords int, tablesDataId int64, tablesDataDelete bool,
	dataColumns []DtoDataColumn, cost float64, costFactual float64, resultTables []DtoResultTable,
	workTables []DtoWorkTable) *DtoVerifyFacility {
	return &DtoVerifyFacility{
		Order_ID:                order_id,
		EstimatedNumbersRecords: estimatedNumbersRecords,
		TablesDataId:            tablesDataId,
		TablesDataDelete:        tablesDataDelete,
		DataColumns:             dataColumns,
		Cost:                    cost,
		CostFactual:             costFactual,
		ResultTables:            resultTables,
		WorkTables:              workTables,
	}
}

func (facility *ViewVerifyFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, column := range facility.DataColumns {
		errors = Validate(&column, errors, req)
	}
	return Validate(facility, errors, req)
}
