package models

import (
	"time"
)

type OrderStatus int

const (
	ORDER_STATUS_COMPLETED OrderStatus = iota + 1
	ORDER_STATUS_NEW
	ORDER_STATUS_OPEN
	ORDER_STATUS_MODERATOR_CONFIRMED
	ORDER_STATUS_CANCEL
	ORDER_STATUS_SUPPLIER_COST_NEW
	ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED
	ORDER_STATUS_PAID
	ORDER_STATUS_MODERATOR_BEGIN
	ORDER_STATUS_SUPPLIER_CLOSE
	ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN
	ORDER_STATUS_MODERATOR_CLOSE
	ORDER_STATUS_ARCHIVE
	ORDER_STATUS_DEL
)

//Структура для организации статуса заказа
type DtoOrderStatus struct {
	Order_ID  int64       `db:"order_id"`  // Идентификатор заказа
	Status_ID OrderStatus `db:"status_id"` // Идентификатор статуса
	Value     bool        `db:"value"`     // Значение
	Comments  string      `db:"comments"`  // Комментарий
	Created   time.Time   `db:"created"`   // Время создания
}

// Конструктор создания объекта статуса сообщений в бд
func NewDtoOrderStatus(order_id int64, status_id OrderStatus, value bool, comments string, created time.Time) *DtoOrderStatus {
	return &DtoOrderStatus{
		Order_ID:  order_id,
		Status_ID: status_id,
		Value:     value,
		Comments:  comments,
		Created:   created,
	}
}
