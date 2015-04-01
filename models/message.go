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

//Структура для организации хранения сообщения
type ViewMessage struct {
	Content string `json:"message" validate:"max=255"` // Содержание
}

type ApiMetaMessage struct {
	NumOfAll int64 `json:"count"` // Общее число сообщений
	NumOfNew int64 `json:"new"`   // Число новых сообщений
}

type ApiShortMessage struct {
	ID int64 `json:"id"` // Уникальный идентификатор
}

type ApiLongMessage struct {
	ID        int64     `json:"id" db:"id"`           // Уникальный идентификатор
	Created   time.Time `json:"created" db:"created"` // Время создания
	IsNew     bool      `json:"new" db:"new"`         // Новое
	IsCreator bool      `json:"isMine" db:"isMine"`   // Прочитано
	User_ID   int64     `json:"userId" db:"userId"`   // Идентификатор автора
	Content   string    `json:"message" db:"message"` // Содержание
}

type MessageSearch struct {
	ID        int64  `query:"id" search:"m.id"`           // Уникальный идентификатор
	Created   string `query:"created" search:"m.created"` // Время создания
	IsNew     bool   `query:"new" search:"new"`           // Новое
	IsCreator bool   `query:"isMine" search:"isMine"`     // Прочитано
	User_ID   int64  `query:"userId" search:"m.user_id"`  // Идентификатор автора
	Content   string `query:"message" search:"m.content"` // Содержание
}

type DtoMessage struct {
	ID       int64     `db:"id"`       // Уникальный идентификатор
	User_ID  int64     `db:"user_id"`  // Идентификатор автора
	Order_ID int64     `db:"order_id"` // Идентификатор обсуждаемого заказа
	Content  string    `db:"content"`  // Содержание
	Created  time.Time `db:"created"`  // Время создания

}

// Конструктор создания объекта сообщения в api
func NewApiMetaMessage(numofall int64, numofnew int64) *ApiMetaMessage {
	return &ApiMetaMessage{
		NumOfAll: numofall,
		NumOfNew: numofnew,
	}
}

func NewApiShortMessage(id int64) *ApiShortMessage {
	return &ApiShortMessage{
		ID: id,
	}
}

func NewApiLongMessage(id int64, created time.Time, isnew bool, iscreator bool, user_id int64, content string) *ApiLongMessage {
	return &ApiLongMessage{
		ID:        id,
		Created:   created,
		IsNew:     isnew,
		IsCreator: iscreator,
		User_ID:   user_id,
		Content:   content,
	}
}

// Конструктор создания объекта сообщения в бд
func NewDtoMessage(id int64, user_id int64, order_id int64, content string, created time.Time) *DtoMessage {
	return &DtoMessage{
		ID:       id,
		User_ID:  user_id,
		Order_ID: order_id,
		Content:  content,
		Created:  created,
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
			errValue = errors.New("Wrong field value")
			break
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
		errField = errors.New("Uknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (message *MessageSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllSearchTags(message)
}

func (message *ViewMessage) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(message, errors, req)
}
