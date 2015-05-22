package models

// Структура для организации хранения таблицы услуги верификации данных
type ApiVerifyTable struct {
	ID           int64            `json:"id" db:"id"`                        // Идентификатор таблицы
	Name         string           `json:"name" db:"name"`                    // Название
	UnitID       int64            `json:"unitId" db:"unit_id"`               // Идентификатор объединения
	TypeID       int              `json:"type" db:"type_id"`                 // Идентификатор типа
	Verification []ApiTableColumn `json:"colsVerification,omitempty" db:"-"` // Колонки верификации данных
}

// Конструктор создания объекта таблицы услуги верификации данных в api
func NewApiVerifyTable(id int64, name string, unitid int64, typeid int, verification []ApiTableColumn) *ApiVerifyTable {
	return &ApiVerifyTable{
		ID:           id,
		Name:         name,
		UnitID:       unitid,
		TypeID:       typeid,
		Verification: verification,
	}
}
