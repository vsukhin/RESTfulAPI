package models

import (
	"mime/multipart"
	"time"
)

// Структура для организации хранения файлов
type ViewFile struct {
	FileData *multipart.FileHeader `form:"filename"` // Содержание файла
}

type ApiFile struct {
	ID int64 `json:"id"` // Уникальный идентификатор файла
}

type ApiFileHead struct {
	ID          int64  `json:"id"`          // Уникальный идентификатор файла
	Name        string `json:"name"`        // Оригинальное имя файла
	Extension   string `json:"extension"`   // Расширение файла
	IsPicture   bool   `json:"isPicture"`   // Файл является картинкой
	Width       int    `json:"width"`       // Ширина картинки в пикселях
	Height      int    `json:"height"`      // Высота картинки в пикселях
	Size        int    `json:"size"`        // Размер файла в байтах
	Modtime     string `json:"modtime"`     // Дата и время модификации файла в формате unixtime
	ContentType string `json:"contentType"` // Формат контента
	Codepage    string `json:"codepage"`    // Кодировка контента
	Hash        string `json:"hash"`        // Контрольная сумма файла рассчитанная алгоритмом sha512
}

type ApiFileObject struct {
	ApiFileHead
	Data string `json:"data"` // Тело файла в формате base64
}

type ApiImage struct {
	ID int64 `json:"key"` // Идентификатор запрашиваемой картинки
}

type DtoFile struct {
	ID                      int64     `db:"id"`                      // Уникальный идентификатор файла
	Name                    string    `db:"name"`                    // Оригинальное имя файла
	Path                    string    `db:"path"`                    // Путь к файлу
	Created                 time.Time `db:"created"`                 // Время создания файла
	Permanent               bool      `db:"permanent"`               // Постоянный файл
	Export_Ready            bool      `db:"export_ready"`            // Готовность экспорта
	Export_Percentage       byte      `db:"export_percentage"`       // Процент готовности экспорта
	Export_Object_ID        int64     `db:"export_object_id"`        // Идентификатор связанного объекта БД для экспорта
	Export_Error            bool      `db:"export_error"`            // Ошибка экспорта
	Export_ErrorDescription string    `db:"export_errordescription"` // Описание ошибки экспорта
	FileData                []byte    `db:"-"`                       // Содержание файла
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

// Конструктор создания объекта файла в api
func NewApiFileHead(id int64, name string, extension string, ispicture bool, width int, height int, size int, modtime string,
	contenttype string, codepage string, hash string) *ApiFileHead {
	return &ApiFileHead{
		ID:          id,
		Name:        name,
		Extension:   extension,
		IsPicture:   ispicture,
		Width:       width,
		Height:      height,
		Size:        size,
		Modtime:     modtime,
		ContentType: contenttype,
		Codepage:    codepage,
		Hash:        hash,
	}
}

func NewApiFileObject(apifilehead ApiFileHead, data string) *ApiFileObject {
	return &ApiFileObject{
		ApiFileHead: apifilehead,
		Data:        data,
	}
}

// Конструктор создания объекта файла в бд
func NewDtoFile(id int64, name string, path string, created time.Time, permanent bool, export_ready bool, export_percentage byte,
	export_object_id int64, export_error bool, export_errordescription string, filedata []byte) *DtoFile {
	return &DtoFile{
		ID:                      id,
		Name:                    name,
		Path:                    path,
		Created:                 created,
		Permanent:               permanent,
		Export_Ready:            export_ready,
		Export_Percentage:       export_percentage,
		Export_Object_ID:        export_object_id,
		Export_Error:            export_error,
		Export_ErrorDescription: export_errordescription,
		FileData:                filedata,
	}
}
