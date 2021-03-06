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

// Структура для организации хранения сообщения
type ViewLongMessage struct {
	Content     string `json:"message" validate:"max=255"`    // Содержание
	Receiver_ID int64  `json:"receiverId" validate:"nonzero"` // Идентификатор получателя
}

type ViewShortMessage struct {
	Content string `json:"message" validate:"max=255"` // Содержание
}

type ApiMetaMessageTotal struct {
	Receiver_ID int64 `db:"receiver_id"` // Идентификатор получателя
	Total       int64 `db:"count"`       // Число сообщений
}

type ApiMetaMessageReceiver struct {
	Receiver_ID int64 `json:"receiverId"` // Идентификатор получателя
	NumOfAll    int64 `json:"count"`      // Общее число сообщений
	NumOfNew    int64 `json:"new"`        // Число новых сообщений
}

type ApiMetaMessage struct {
	NumOfAll  int64                    `json:"count"`               // Общее число сообщений
	NumOfNew  int64                    `json:"new"`                 // Число новых сообщений
	Receivers []ApiMetaMessageReceiver `json:"receivers,omitempty"` // Получатели
}

type ApiShortMessage struct {
	ID int64 `json:"id"` // Уникальный идентификатор
}

type ApiLongMessage struct {
	ID          int64     `json:"id" db:"id"`                 // Уникальный идентификатор
	Created     time.Time `json:"created" db:"created"`       // Время создания
	IsNew       bool      `json:"new" db:"new"`               // Новое
	IsCreator   bool      `json:"isMine" db:"isMine"`         // Прочитано
	User_ID     int64     `json:"userId" db:"userId"`         // Идентификатор автора
	Receiver_ID int64     `json:"receiverId" db:"receiverId"` // Идентификатор получателя
	Content     string    `json:"message" db:"message"`       // Содержание
}

type MessageSearch struct {
	ID          int64  `query:"id" search:"m.id"`                                                 // Уникальный идентификатор
	Created     string `query:"created" search:"m.created" group:"convert(m.created using utf8)"` // Время создания
	IsNew       bool   `query:"new" search:"new"`                                                 // Новое
	IsCreator   bool   `query:"isMine" search:"isMine"`                                           // Прочитано
	User_ID     int64  `query:"userId" search:"m.user_id"`                                        // Идентификатор автора
	Receiver_ID int64  `query:"receiverId" search:"m.receiver_id"`                                // Идентификатор получателя
	Content     string `query:"message" search:"m.content" group:"m.content"`                     // Содержание
}

type DtoMessage struct {
	ID          int64     `db:"id"`          // Уникальный идентификатор
	User_ID     int64     `db:"user_id"`     // Идентификатор автора
	Order_ID    int64     `db:"order_id"`    // Идентификатор обсуждаемого заказа
	Receiver_ID int64     `db:"receiver_id"` // Идентификатор получателя
	Content     string    `db:"content"`     // Содержание
	Created     time.Time `db:"created"`     // Время создания
}

// Конструктор создания объекта сообщения в api
func NewApiMetaMessageTotal(receiver_id int64, total int64) *ApiMetaMessageTotal {
	return &ApiMetaMessageTotal{
		Receiver_ID: receiver_id,
		Total:       total,
	}
}

func NewApiMetaMessageReceiver(receiver_id int64, numofall int64, numofnew int64) *ApiMetaMessageReceiver {
	return &ApiMetaMessageReceiver{
		Receiver_ID: receiver_id,
		NumOfAll:    numofall,
		NumOfNew:    numofnew,
	}
}

func NewApiMetaMessage(numofall int64, numofnew int64, receivers []ApiMetaMessageReceiver) *ApiMetaMessage {
	return &ApiMetaMessage{
		NumOfAll:  numofall,
		NumOfNew:  numofnew,
		Receivers: receivers,
	}
}

func NewApiShortMessage(id int64) *ApiShortMessage {
	return &ApiShortMessage{
		ID: id,
	}
}

func NewApiLongMessage(id int64, created time.Time, isnew bool, iscreator bool, user_id int64, receiver_id int64, content string) *ApiLongMessage {
	return &ApiLongMessage{
		ID:          id,
		Created:     created,
		IsNew:       isnew,
		IsCreator:   iscreator,
		User_ID:     user_id,
		Receiver_ID: receiver_id,
		Content:     content,
	}
}

// Конструктор создания объекта сообщения в бд
func NewDtoMessage(id int64, user_id int64, order_id int64, receiver_id int64, content string, created time.Time) *DtoMessage {
	return &DtoMessage{
		ID:          id,
		User_ID:     user_id,
		Order_ID:    order_id,
		Receiver_ID: receiver_id,
		Content:     content,
		Created:     created,
	}
}

func (message *MessageSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, message), nil
}

func (message *MessageSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, message)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "userId":
		fallthrough
	case "receiverId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "message":
		fallthrough
	case "created":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "new":
		fallthrough
	case "isMine":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (message *MessageSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(message)
}

func (message *ViewLongMessage) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(message, errors, req)
}

func (message *ViewShortMessage) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(message, errors, req)
}
