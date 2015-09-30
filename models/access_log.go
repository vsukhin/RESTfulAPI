package models

import (
	"time"
)

// Структура для хранения данных доступа
type DtoAccessLog struct {
	ID          int64     `db:"id"`          // Уникальный идентификатор доступа
	IP_Address  string    `db:"ip_address"`  // IP адрес доступа
	Reverse_DNS string    `db:"reverse_dns"` // Обратный DNS доступа
	User_Agent  string    `db:"user_agent"`  // User Agent доступа
	Referer     string    `db:"referer"`     // Referer доступа
	URL         string    `db:"url"`         // URL доступа
	Created     time.Time `db:"created"`     // Время доступа
}

// Конструктор создания объекта доступа в бд
func NewAccessLog(id int64, ip_address string, reverse_dns string, user_agent string, referer string, url string, created time.Time) *DtoAccessLog {
	return &DtoAccessLog{
		ID:          id,
		IP_Address:  ip_address,
		Reverse_DNS: reverse_dns,
		User_Agent:  user_agent,
		Referer:     referer,
		URL:         url,
		Created:     created,
	}
}
