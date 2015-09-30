package models

import (
	"errors"
	"strconv"
	"strings"
)

// Структура для организации хранения сервиса поставщика
type ApiShortSupplierFacility struct {
	Supplier_ID int64  `json:"id" db:"id"`                 // Идентификатор поставщика
	Name        string `json:"name" db:"name"`             // Название
	Position    int64  `json:"position" db:"position"`     // Позиция
	Rating      int    `json:"rating" db:"rating"`         // Рейтинг
	Throughput  int    `json:"throughput" db:"throughput"` // Пропускная способность
}

type ApiLongSupplierFacility struct {
	Supplier_ID int64  `json:"id" db:"id"`                 // Идентификатор поставщика
	Name        string `json:"name" db:"name"`             // Название
	Service_ID  int64  `json:"serviceId" db:"serviceId"`   // Идентификатор услуги
	Position    int64  `json:"position" db:"position"`     // Позиция
	Rating      int    `json:"rating" db:"rating"`         // Рейтинг
	Throughput  int    `json:"throughput" db:"throughput"` // Пропускная способность
}

type SupplierFacilitySearch struct {
	Supplier_ID int64  `query:"id" search:"f.supplier_id"`                             // Идентификатор поставщика
	Name        string `query:"name" search:"u.name" group:"u.name"`                   // Название
	Service_ID  int64  `query:"serviceId" search:"f.service_id"`                       // Идентификатор услуги
	Position    int64  `query:"position" search:"f.position"`                          // Позиция
	Rating      int    `query:"rating" search:"f.rating" group:"f.rating"`             // Рейтинг
	Throughput  int    `query:"throughput" search:"f.throughput" group:"f.throughput"` // Пропускная способность
	Project_ID  int64  `query:"project" search:"project"`                              // Идентификатор проекта
	Order_ID    int64  `query:"order" search:"order"`                                  // Идентификатор заказа
}

type SupplierFacilitySort struct {
	Supplier_ID int64  `query:"id"`         // Идентификатор поставщика
	Name        string `query:"name"`       // Название
	Service_ID  int64  `query:"serviceId"`  // Идентификатор услуги
	Position    int64  `query:"position"`   // Позиция
	Rating      int    `query:"rating"`     // Рейтинг
	Throughput  int    `query:"throughput"` // Пропускная способность
}

type DtoSupplierFacility struct {
	Supplier_ID int64 `db:"supplier_id"` // Идентификатор поставщика
	Service_ID  int64 `db:"service_id"`  // Идентификатор сервиса
	Position    int64 `db:"position"`    // Позиция
	Rating      int   `db:"rating"`      // Рейтинг
	Throughput  int   `db:"throughput"`  // Пропускная способность
}

// Конструктор создания объекта сервиса поставщика в api
func NewApiShortSupplierFacility(supplier_id int64, name string, position int64, rating int, throughput int) *ApiShortSupplierFacility {
	return &ApiShortSupplierFacility{
		Supplier_ID: supplier_id,
		Name:        name,
		Position:    position,
		Rating:      rating,
		Throughput:  throughput,
	}
}

func NewApiLongSupplierFacility(supplier_id int64, name string, service_id int64, position int64, rating int, throughput int) *ApiLongSupplierFacility {
	return &ApiLongSupplierFacility{
		Supplier_ID: supplier_id,
		Name:        name,
		Service_ID:  service_id,
		Position:    position,
		Rating:      rating,
		Throughput:  throughput,
	}
}

// Конструктор создания объекта сервиса поставщика в бд
func NewDtoSupplierFacility(supplier_id int64, service_id int64, position int64, rating int, throughput int) *DtoSupplierFacility {
	return &DtoSupplierFacility{
		Supplier_ID: supplier_id,
		Service_ID:  service_id,
		Position:    position,
		Rating:      rating,
		Throughput:  throughput,
	}
}

func (supplierfacility *SupplierFacilitySearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, supplierfacility), nil
}

func (supplierfacility *SupplierFacilitySearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, supplierfacility)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "serviceId":
		fallthrough
	case "position":
		fallthrough
	case "project":
		fallthrough
	case "order":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "rating":
		fallthrough
	case "throughput":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (supplierfacility *SupplierFacilitySearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(supplierfacility)
}

func (supplierfacility *SupplierFacilitySort) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, supplierfacility), nil
}
