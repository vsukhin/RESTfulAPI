package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

type DataFormat int

const (
	DATA_FORMAT_UKNOWN DataFormat = iota
	DATA_FORMAT_TXT
	DATA_FORMAT_CSV
	DATA_FORMAT_SSV
)

const (
	EXPORT_DATA_ALL     = "all"
	EXPORT_DATA_INVALID = "incorrect"
	EXPORT_DATA_VALID   = "correct"
)

// Структура для организации хранения импорта-экспорта
type ViewImportTable struct {
	File_ID   string `json:"fileId" validate:"nonzero,min=1,max=255"` // Уникальный идентификатор файла
	HasHeader bool   `json:"names"`                                   // Есть строка заголовка
}

type ViewImportColumn struct {
	ID       int64  `json:"id" db:"id" validate:"nonzero"`         // Уникальный идентификатор временной колонки таблицы
	Name     string `json:"name" validate:"nonzero,min=1,max=255"` // Название колонки таблицы
	Position int64  `json:"position" validate:"min=0"`             // Позиция
	TypeID   int    `json:"typeId"`                                // Идентификатор типа
	Use      bool   `json:"pass"`                                  // Импортируется
}

type ViewImportColumns []ViewImportColumn

type ViewExportTable struct {
	Data_Format_ID int    `json:"format" validate:"nonzero"`             // Уникальный идентификатор формата данных
	Type           string `json:"rows" validate:"nonzero,min=1,max=255"` // Тип экспортируемых данных
}

type ApiImportTable struct {
	ID int64 `json:"id"` // Уникальный идентификатор временной таблицы
}

type ApiImportStatus struct {
	Ready          bool            `json:"ready" `             // Готова
	Percentage     byte            `json:"percent"`            // Процент готовности
	Percentages    []ApiImportStep `json:"percents,omitempty"` // Процент готовности по шагам
	NumOfCols      int64           `json:"columns"`            // Число колонок
	NumOfRows      int64           `json:"rows"`               // Число строк
	NumOfWrongRows int64           `json:"errorRows"`          // Количество строк с неверным числом столбцов
}

type ApiImportColumn struct {
	ID       int64  `json:"id" db:"id"`             // Уникальный идентификатор колонки таблицы
	Name     string `json:"name" db:"name"`         // Название
	Position int64  `json:"position" db:"position"` // Позиция
}

type ApiMetaExportTable struct {
	Formats []ApiDataFormat `json:"formats,omitempty" ` // Список форматов для выгрузки данных
	URL     string          `json:"url" `               // URL для выгрузки данных
}

type ApiExportStatus struct {
	Ready      bool   `json:"ready" `  // Готова
	Percentage byte   `json:"percent"` // Процент готовности
	ExpiredAt  string `json:"timeout"` // Срок валидности
}

// Конструктор создания объекта импорта-экспорта
func NewApiImportTable(id int64) *ApiImportTable {
	return &ApiImportTable{
		ID: id,
	}
}

func NewApiImportStatus(ready bool, percentage byte, percentages []ApiImportStep,
	numofcols int64, numofrows int64, numofwrongrows int64) *ApiImportStatus {
	return &ApiImportStatus{
		Ready:          ready,
		Percentage:     percentage,
		Percentages:    percentages,
		NumOfCols:      numofcols,
		NumOfRows:      numofrows,
		NumOfWrongRows: numofwrongrows,
	}
}

func NewApiImportColumn(id int64, name string, position int64) *ApiImportColumn {
	return &ApiImportColumn{
		ID:       id,
		Name:     name,
		Position: position,
	}
}

func NewApiMetaExportTable(formats []ApiDataFormat, url string) *ApiMetaExportTable {
	return &ApiMetaExportTable{
		Formats: formats,
		URL:     url,
	}
}

func NewApiExportStatus(ready bool, percentage byte, expiredat string) *ApiExportStatus {
	return &ApiExportStatus{
		Ready:      ready,
		Percentage: percentage,
		ExpiredAt:  expiredat,
	}
}

func (importtable *ViewImportTable) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(importtable, errors, req)
}

func (importcolumn ViewImportColumn) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(&importcolumn, errors, req)
}

func (importcolumns ViewImportColumns) GetIDs() []int64 {
	ids := new([]int64)
	for _, importcolumn := range importcolumns {
		*ids = append(*ids, importcolumn.ID)
	}

	return *ids
}

func GetDataSeparator(dataformat DataFormat) (separator rune) {
	switch dataformat {
	case DATA_FORMAT_TXT:
		separator = '\t'
	case DATA_FORMAT_CSV:
		separator = ','
	case DATA_FORMAT_SSV:
		separator = ';'
	default:
		separator = ','
	}

	return separator
}
