package models

import (
	"time"
)

// Структура для хранения запроса/ответа поставщику
type ApiSupplierRequest struct {
	Supplier_ID   int64     `json:"supplierId" db:"supplier_id"`      // Идентификатор поставщика
	RequestDate   time.Time `json:"requestDate" db:"requestDate"`     // Дата и время отправки запроса
	Responded     bool      `json:"responded" db:"responded"`         // Ответил
	RespondedDate time.Time `json:"respondedDate" db:"respondedDate"` // Дата и время ответа
	EstimatedCost float64   `json:"estimatedCost" db:"estimatedCost"` // Расчётная стоимость заказа
	MyChoice      bool      `json:"myChoice" db:"myChoice"`           // Пользователь выбрал поставщика
}

type DtoSupplierRequest struct {
	Order_ID      int64     `db:"order_id"`      // Идентификатор заказа
	Supplier_ID   int64     `db:"supplier_id"`   // Идентификатор поставщика
	RequestDate   time.Time `db:"requestDate"`   // Дата и время отправки запроса
	Responded     bool      `db:"responded"`     // Ответил
	RespondedDate time.Time `db:"respondedDate"` // Дата и время ответа
	EstimatedCost float64   `db:"estimatedCost"` // Расчётная стоимость заказа
	MyChoice      bool      `db:"myChoice"`      // Пользователь выбрал поставщика
}

// Конструктор создания объекта запроса/ответа поставщику в api
func NewApiSupplierRequest(supplier_id int64, requestDate time.Time, responded bool,
	respondedDate time.Time, estimatedCost float64, myChoice bool) *ApiSupplierRequest {
	return &ApiSupplierRequest{
		Supplier_ID:   supplier_id,
		RequestDate:   requestDate,
		Responded:     responded,
		RespondedDate: respondedDate,
		EstimatedCost: estimatedCost,
		MyChoice:      myChoice,
	}
}

// Конструктор создания объекта запроса/ответа поставщику в бд
func NewDtoSupplierRequest(order_id int64, supplier_id int64, requestDate time.Time, responded bool,
	respondedDate time.Time, estimatedCost float64, myChoice bool) *DtoSupplierRequest {
	return &DtoSupplierRequest{
		Order_ID:      order_id,
		Supplier_ID:   supplier_id,
		RequestDate:   requestDate,
		Responded:     responded,
		RespondedDate: respondedDate,
		EstimatedCost: estimatedCost,
		MyChoice:      myChoice,
	}
}
