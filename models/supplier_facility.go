package models

// Структура для организации хранения сервиса поставщика
type ApiSupplierFacility struct {
	Supplier_ID int64  `json:"id" db:"supplier_id"`        // Идентификатор поставщика
	Name        string `json:"name" db:"name"`             // Название
	Position    int64  `json:"position" db:"position"`     // Позиция
	Rating      int    `json:"rating" db:"rating"`         // Рейтинг
	Throughput  int    `json:"throughput" db:"throughput"` // Пропускная способность
}

type DtoSupplierFacility struct {
	Supplier_ID int64 `db:"supplier_id"` // Идентификатор поставщика
	Service_ID  int64 `db:"service_id"`  // Идентификатор сервиса
	Position    int64 `db:"position"`    // Позиция
	Rating      int   `db:"rating"`      // Рейтинг
	Throughput  int   `db:"throughput"`  // Пропускная способность
}

// Конструктор создания объекта сервиса поставщика в api
func NewApiSupplierFacility(supplier_id int64, name string, position int64, rating int, throughput int) *ApiSupplierFacility {
	return &ApiSupplierFacility{
		Supplier_ID: supplier_id,
		Name:        name,
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
