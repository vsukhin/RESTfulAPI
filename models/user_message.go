package models

import (
	"time"
)

// Структура для организации статуса сообщения
type DtoUserMessage struct {
	User_ID    int64     `db:"user_id"`    // Идентификатор пользователя
	Message_ID int64     `db:"message_id"` // Идентификатор сообщения
	Value      bool      `db:"value"`      // Значение
	Created    time.Time `db:"created"`    // Время создания
}

// Конструктор создания объекта статуса сообщения в бд
func NewDtoUserMessage(user_id int64, message_id int64, value bool, created time.Time) *DtoUserMessage {
	return &DtoUserMessage{
		User_ID:    user_id,
		Message_ID: message_id,
		Value:      value,
		Created:    created,
	}
}
