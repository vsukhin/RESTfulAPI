package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
)

// Структура для хранения примера анкеты
type ViewInputFile struct {
	File_ID int64 `json:"id" validate:"nonzero"` // Идентификатор файла
}

type ApiInputFile struct {
	File_ID int64  `json:"id" db:"file_id"` // Идентификатор файла
	Name    string `json:"name" db:"name"`  // Название
}

type DtoInputFile struct {
	Order_ID int64 `db:"order_id"` // Идентификатор заказа
	File_ID  int64 `db:"file_id"`  // Идентификатор файла
}

// Конструктор создания объекта примера анкеты в api
func NewApiInputFile(file_id int64, name string) *ApiInputFile {
	return &ApiInputFile{
		File_ID: file_id,
		Name:    name,
	}
}

// Конструктор создания объекта примера анкеты в бд
func NewDtoInputFile(order_id int64, file_id int64) *DtoInputFile {
	return &DtoInputFile{
		Order_ID: order_id,
		File_ID:  file_id,
	}
}

func (file *ViewInputFile) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(file, errors, req)
}
