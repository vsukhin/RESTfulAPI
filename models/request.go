package models

import (
	"time"
)

// Структура для организации хранения обращения
type DtoRequest struct {
	IP_Address  string    `db:"ip_address"`  // IP адрес обращения
	Method      string    `db:"method"`      // Вызываемый метод
	LastUpdated time.Time `db:"lastUpdated"` // Время последнего обновления
	Hits        int64     `db:"hits"`        // Количество обращений
}

// Конструктор создания объекта обращения в бд
func NewDtoRequest(ip_address string, method string, lastupdated time.Time, hits int64) *DtoRequest {
	return &DtoRequest{
		IP_Address:  ip_address,
		Method:      method,
		LastUpdated: lastupdated,
		Hits:        hits,
	}
}
