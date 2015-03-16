package models

import (
	"mime/multipart"
	"time"
)

//Структура для организации хранения файлов
type ViewFile struct {
	FileData *multipart.FileHeader `form:"filename"` // Содержание файла
}

type ApiFile struct {
	ID int64 `json:"id"` // Уникальный идентификатор файла
}

type ApiImage struct {
	ID int64 `json:"key"` // Идентификатор запрашиваемой картинки
}

type DtoFile struct {
	ID         int64     `db:"id"`         // Уникальный идентификатор файла
	Name       string    `db:"name"`       // Оригинальное имя файла
	Path       string    `db:"path"`       // Путь к файлу
	Created    time.Time `db:"created"`    // Время создания файла
	Permanent  bool      `db:"permanent"`  // Постоянный файл
	Ready      bool      `db:"ready"`      // Готовность
	Percentage byte      `db:"percentage"` // Процент готовности
	Object_ID  int64     `db:"object_id"`  // Идентификатор связанного объекта БД
	FileData   []byte    `db:"-"`          // Содержание файла
}

func NewApiFile(id int64) *ApiFile {
	return &ApiFile{
		ID: id,
	}
}

func NewApiImage(id int64) *ApiImage {
	return &ApiImage{
		ID: id,
	}
}

// Конструктор создания объекта файла в бд
func NewDtoFile(id int64, name string, path string, created time.Time, permanent bool, ready bool, percentage byte,
	object_id int64, filedata []byte) *DtoFile {
	return &DtoFile{
		ID:         id,
		Name:       name,
		Path:       path,
		Created:    created,
		Permanent:  permanent,
		Ready:      ready,
		Percentage: percentage,
		Object_ID:  object_id,
		FileData:   filedata,
	}
}
