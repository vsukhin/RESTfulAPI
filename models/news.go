package models

import (
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
