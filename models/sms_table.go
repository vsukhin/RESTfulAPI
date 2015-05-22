package models

// Структура для организации хранения таблицы sms услуги
type ApiSMSTable struct {
	ID           int64            `json:"id" db:"id"`                       // Идентификатор таблицы
	Name         string           `json:"name" db:"name"`                   // Название
	UnitID       int64            `json:"unitId" db:"unit_id"`              // Идентификатор объединения
	TypeID       int              `json:"type" db:"type_id"`                // Идентификатор типа
	MobilePhones []ApiTableColumn `json:"colsMobilePhone,omitempty" db:"-"` // Колонки мобильных телефонов
	Messages     []ApiTableColumn `json:"colsMessage,omitempty" db:"-"`     // Колонки сообщений
	SMSSenders   []ApiTableColumn `json:"colsMessageFrom,omitempty" db:"-"` // Колонки отправителей
	Birthdays    []ApiTableColumn `json:"colsBirthday,omitempty" db:"-"`    // Колонки дней рождений
}

// Конструктор создания объекта таблицы sms услуги в api
func NewApiSMSTable(id int64, name string, unitid int64, typeid int, mobilephones []ApiTableColumn, messages []ApiTableColumn,
	smssenders []ApiTableColumn, birthdays []ApiTableColumn) *ApiSMSTable {
	return &ApiSMSTable{
		ID:           id,
		Name:         name,
		UnitID:       unitid,
		TypeID:       typeid,
		MobilePhones: mobilephones,
		Messages:     messages,
		SMSSenders:   smssenders,
		Birthdays:    birthdays,
	}
}
