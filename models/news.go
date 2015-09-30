package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Структура для организации хранения новости
type ApiRSSFeed struct {
	XMLName  struct{}      `xml:"rss"`          // RSS поток
	Version  string        `xml:"version,attr"` // Версия
	Channels *[]ApiChannel `xml:"channel"`      // Каналы
}

type ApiChannel struct {
	Title       string     `xml:"title"`       // Название канала
	Description string     `xml:"description"` // Описание
	URL         string     `xml:"link"`        // Прямая ссылка
	Language    string     `xml:"language"`    // Язык
	News        *[]DtoNews `xml:"item"`        // Новости
}

type ApiNews struct {
	ID          int64     `json:"id" db:"id"`                     // Уникальный идентификатор новости
	Created     time.Time `json:"date" db:"date"`                 // Время создания
	Title       string    `json:"title" db:"title"`               // Заголовок
	Description string    `json:"messageShort" db:"messageShort"` // Содержание
}

type NewsSearch struct {
	ID          int64     `query:"id" search:"id"`                                            // Уникальный идентификатор компании
	Created     time.Time `query:"date" search:"created" group:"convert(created using utf8)"` // Время создания
	Title       string    `query:"title" search:"title" group:"title"`                        // Заголовок
	Description string    `query:"messageShort" search:"description" group:"description"`     // Содержание
}

type DtoNews struct {
	ID          int64     `db:"id" xml:"-"`                    // Уникальный идентификатор новости
	Title       string    `db:"title" xml:"title"`             // Заголовок
	Description string    `db:"description" xml:"description"` // Содержание
	Created     time.Time `db:"created" xml:"-"`               // Время создания
	Language    string    `db:"language" xml:"-"`              // Язык
	Active      bool      `db:"active" xml:"-"`                // Активная
	URL         string    `db:"-" xml:"link"`                  // Прямая ссылка
}

// Конструктор создания объекта новости в api
func NewApiRSSFeed(version string, channels *[]ApiChannel) *ApiRSSFeed {
	return &ApiRSSFeed{
		Version:  version,
		Channels: channels,
	}
}

func NewApiChannel(title string, description string, url string, language string, news *[]DtoNews) *ApiChannel {
	return &ApiChannel{
		Title:       title,
		Description: description,
		URL:         url,
		Language:    language,
		News:        news,
	}
}

func NewApiNews(id int64, created time.Time, title string, description string) *ApiNews {
	return &ApiNews{
		ID:          id,
		Created:     created,
		Title:       title,
		Description: description,
	}
}

// Конструктор создания объекта новости в бд
func NewDtoNews(id int64, title string, description string, created time.Time, language string, active bool, url string) *DtoNews {
	return &DtoNews{
		ID:          id,
		Title:       title,
		Description: description,
		Created:     created,
		Language:    language,
		Active:      active,
		URL:         url,
	}
}

func (news *NewsSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, news), nil
}

func (news *NewsSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, news)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "date":
		fallthrough
	case "title":
		fallthrough
	case "messageShort":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (news *NewsSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(news)
}
