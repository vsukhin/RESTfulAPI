package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для организации хранения сервиса ввода анкет заказа
type ViewRecognizeFacility struct {
	EstimatedNumbersForm         int                   `json:"estimatedNumbersForm" validate:"min=0"` // Предполагаемое количество анкет
	EstimatedCalculationOnFields bool                  `json:"estimatedCalculationOnFields"`          // Расчёт предварительной стоимости на основе полей
	EstimatedFields              []ViewApiInputField   `json:"estimatedFields"`                       // Прогноз количества вводимых типов полей
	InputProducts                []ViewApiInputProduct `json:"priceIncrease"`                         // Позиции прайс листа
	RequiredFields               string                `json:"requiredFields"`                        // Перечисление названий полей анкеты которые не должны быть пустыми
	LoadDefectiveForms           bool                  `json:"loadDefectiveForms"`                    // Что делать с бракованными анкетами
	CommentsForSupplier          string                `json:"commentsForSupplier"`                   // Комментарии
	EstimatedFormFiles           []ApiInputFile        `json:"estimatedFormFiles"`                    // Примеры загруженных анкет
}

type ApiRecognizeFacility struct {
	EstimatedNumbersForm         int                   `json:"estimatedNumbersForm" db:"estimatedNumbersForm"`                 // Предполагаемое количество анкет
	EstimatedCalculationOnFields bool                  `json:"estimatedCalculationOnFields" db:"estimatedCalculationOnFields"` // Расчёт предварительной стоимости на основе полей
	EstimatedFields              []ViewApiInputField   `json:"estimatedFields,omitempty" db:"-"`                               // Прогноз количества вводимых типов полей
	InputProducts                []ViewApiInputProduct `json:"priceIncrease,omitempty" db:"-"`                                 // Позиции прайс листа
	RequiredFields               string                `json:"requiredFields" db:"requiredFields"`                             // Перечисление названий полей анкеты которые не должны быть пустыми
	LoadDefectiveForms           bool                  `json:"loadDefectiveForms" db:"loadDefectiveForms"`                     // Что делать с бракованными анкетами
	CommentsForSupplier          string                `json:"commentsForSupplier" db:"commentsForSupplier"`                   // Комментарии
	EstimatedFormFiles           []ApiInputFile        `json:"estimatedFormFiles,omitempty" db:"-"`                            // Примеры загруженных анкет
	RequestsSend                 bool                  `json:"requestsSend" db:"requestsSend"`                                 // Можно высылать
	RequestsCancel               bool                  `json:"requestsCancel" db:"requestsCancel"`                             // Пользователь отменил заказ или выбрал поставщика
	SupplierRequests             []ApiSupplierRequest  `json:"supplierRequests,omitempty" db:"-"`                              // Запросы/ответы поставщикам услуг
	Cost                         float64               `json:"cost" db:"cost"`                                                 // Сумма заказа исходя из расчётных показателей заказа
	CostFactual                  float64               `json:"costFactual" db:"costFactual"`                                   // Текущая стоимость заказа
	Ftp                          ApiInputFtp           `json:"ftp" db:"-"`                                                     // Реквизиты ftp доступа
	ResultTables                 []ApiResultTable      `json:"resultTables,omitempty" db:"-"`                                  // Таблицы результатов
	WorkTables                   []ApiWorkTable        `json:"workTables,omitempty" db:"-"`                                    // Рабочие таблицы
}

type DtoRecognizeFacility struct {
	Order_ID                     int64                `db:"order_id"`                     // Идентификатор заказа
	EstimatedNumbersForm         int                  `db:"estimatedNumbersForm"`         // Предполагаемое количество анкет
	EstimatedCalculationOnFields bool                 `db:"estimatedCalculationOnFields"` // Расчёт предварительной стоимости на основе полей
	EstimatedFields              []DtoInputField      `db:"-"`                            // Прогноз количества вводимых типов полей
	InputProducts                []DtoInputProduct    `db:"-"`                            // Позиции прайс листа
	RequiredFields               string               `db:"requiredFields"`               // Перечисление названий полей анкеты которые не должны быть пустыми
	LoadDefectiveForms           bool                 `db:"loadDefectiveForms"`           // Что делать с бракованными анкетами
	CommentsForSupplier          string               `db:"commentsForSupplier"`          // Комментарии
	EstimatedFormFiles           []DtoInputFile       `db:"-"`                            // Примеры загруженных анкет
	RequestsSend                 bool                 `db:"requestsSend"`                 // Можно высылать
	RequestsCancel               bool                 `db:"requestsCancel"`               // Пользователь отменил заказ или выбрал поставщика
	SupplierRequests             []DtoSupplierRequest `db:"-"`                            // Запросы/ответы поставщикам услуг
	Cost                         float64              `db:"cost"`                         // Сумма заказа исходя из расчётных показателей заказа
	CostFactual                  float64              `db:"costFactual"`                  // Текущая стоимость заказа
	Ftp                          DtoInputFtp          `db:"-"`                            // Реквизиты ftp доступа
	ResultTables                 []DtoResultTable     `db:"-"`                            // Таблицы результатов
	WorkTables                   []DtoWorkTable       `db:"-"`                            // Рабочие таблицы
}

// Конструктор создания объекта сервиса ввода анкет заказа в api
func NewApiRecognizeFacility(estimatedNumbersForm int, estimatedCalculationOnFields bool, estimatedFields []ViewApiInputField,
	inputProducts []ViewApiInputProduct, requiredFields string, loadDefectiveForms bool, commentsForSupplier string,
	estimatedFormFiles []ApiInputFile, requestsSend bool, requestsCancel bool, supplierRequests []ApiSupplierRequest,
	cost float64, costFactual float64, ftp ApiInputFtp, resultTables []ApiResultTable, workTables []ApiWorkTable) *ApiRecognizeFacility {
	return &ApiRecognizeFacility{
		EstimatedNumbersForm:         estimatedNumbersForm,
		EstimatedCalculationOnFields: estimatedCalculationOnFields,
		EstimatedFields:              estimatedFields,
		InputProducts:                inputProducts,
		RequiredFields:               requiredFields,
		LoadDefectiveForms:           loadDefectiveForms,
		CommentsForSupplier:          commentsForSupplier,
		EstimatedFormFiles:           estimatedFormFiles,
		RequestsSend:                 requestsSend,
		RequestsCancel:               requestsCancel,
		SupplierRequests:             supplierRequests,
		Cost:                         cost,
		CostFactual:                  costFactual,
		Ftp:                          ftp,
		ResultTables:                 resultTables,
		WorkTables:                   workTables,
	}
}

// Конструктор создания объекта сервиса ввода анкет заказа в бд
func NewDtoRecognizeFacility(order_id int64, estimatedNumbersForm int, estimatedCalculationOnFields bool, estimatedFields []DtoInputField,
	inputProducts []DtoInputProduct, requiredFields string, loadDefectiveForms bool, commentsForSupplier string,
	estimatedFormFiles []DtoInputFile, requestsSend bool, requestsCancel bool, supplierRequests []DtoSupplierRequest,
	cost float64, costFactual float64, ftp DtoInputFtp, resultTables []DtoResultTable, workTables []DtoWorkTable) *DtoRecognizeFacility {
	return &DtoRecognizeFacility{
		Order_ID:                     order_id,
		EstimatedNumbersForm:         estimatedNumbersForm,
		EstimatedCalculationOnFields: estimatedCalculationOnFields,
		EstimatedFields:              estimatedFields,
		InputProducts:                inputProducts,
		RequiredFields:               requiredFields,
		LoadDefectiveForms:           loadDefectiveForms,
		CommentsForSupplier:          commentsForSupplier,
		EstimatedFormFiles:           estimatedFormFiles,
		RequestsSend:                 requestsSend,
		RequestsCancel:               requestsCancel,
		SupplierRequests:             supplierRequests,
		Cost:                         cost,
		CostFactual:                  costFactual,
		Ftp:                          ftp,
		ResultTables:                 resultTables,
		WorkTables:                   workTables,
	}
}

func (facility *ViewRecognizeFacility) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	for _, field := range facility.EstimatedFields {
		errors = Validate(&field, errors, req)
	}
	for _, product := range facility.InputProducts {
		errors = Validate(&product, errors, req)
	}
	for _, file := range facility.EstimatedFormFiles {
		errors = Validate(&file, errors, req)
	}
	return Validate(facility, errors, req)
}
