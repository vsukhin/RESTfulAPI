package models

// Структура для организации хранения таблицы hlr услуги
type ApiHLRTable struct {
	ID           int64            `json:"id" db:"id"`                       // Идентификатор таблицы
	Name         string           `json:"name" db:"name"`                   // Название
	UnitID       int64            `json:"unitId" db:"unit_id"`              // Идентификатор объединения
	TypeID       int              `json:"type" db:"type_id"`                // Идентификатор типа
	MobilePhones []ApiTableColumn `json:"colsMobilePhone,omitempty" db:"-"` // Колонки мобильных телефонов
}

// Конструктор создания объекта таблицы hlr услуги в api
func NewApiHLRTable(id int64, name string, unitid int64, typeid int, mobilephones []ApiTableColumn) *ApiHLRTable {
	return &ApiHLRTable{
		ID:           id,
		Name:         name,
		UnitID:       unitid,
		TypeID:       typeid,
		MobilePhones: mobilephones,
	}
}
