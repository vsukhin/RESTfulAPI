package models

import (
	"errors"
	"fmt"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	MAX_STEP_NUMBER = 4
)

// Структура для организации хранения заказа
type ViewShortOrder struct {
	Name string `json:"name" validate:"max=255"` // Название
}

type ViewMiddleOrder struct {
	Name               string `json:"name" validate:"max=255"`  // Название
	Step               byte   `json:"step" validate:"min=0"`    // Шаг
	IsAssembled        bool   `json:"completed"`                // Собран
	Facility_ID        int64  `json:"type"`                     // Идентификатор услуги
	Supplier_ID        int64  `json:"supplierId"`               // Идентификатор поставщика
	IsNewCostConfirmed bool   `json:"customerNewCostConfirmed"` // Новая цена утверждена
}

type ViewLongOrder struct {
	Name               string  `json:"name" validate:"max=255"`                   // Название
	Step               byte    `json:"step" validate:"min=0"`                     // Шаг
	IsAssembled        bool    `json:"completed"`                                 // Собран
	IsConfirmed        bool    `json:"moderatorConfirmed"`                        // Утвержден
	Facility_ID        int64   `json:"type"`                                      // Идентификатор услуги
	Supplier_ID        int64   `json:"supplierId"`                                // Идентификатор поставщика
	IsNew              bool    `json:"new"`                                       // Новый
	IsOpen             bool    `json:"open"`                                      // Открыт
	IsCancelled        bool    `json:"cancel"`                                    // Отказ
	Reason             string  `json:"cancelDescription" validate:"max=255"`      // Причина отказа
	Execution_Forecast int     `json:"supplierForecastWorkDays" validate:"min=0"` // Прогноз исполнения
	Proposed_Price     float64 `json:"supplierCost" validate:"min=0"`             // Предложенная цена
	IsNewCost          bool    `json:"supplierCostNew"`                           // Новая цена
	IsNewCostConfirmed bool    `json:"customerNewCostConfirmed"`                  // Новая цена утверждена
	IsPaid             bool    `json:"paid"`                                      // Оплачен
	IsStarted          bool    `json:"moderatorBegin"`                            // Запущен
	IsNewCharge        bool    `json:"supplierFactualCostNew"`                    // Предложена фактическая цена
	Charged_Fee        float64 `json:"supplierFactualCost" validate:"min=0"`      // Фактическая цена
	IsExecuted         bool    `json:"supplierClose"`                             // Выполнен
	IsDocumented       bool    `json:"moderatorDocumentsGotten"`                  // Документы имеются
	IsClosed           bool    `json:"moderatorClose"`                            // Закрыт
	IsArchived         bool    `json:"archive"`                                   // Архивирован
	IsDeleted          bool    `json:"del"`                                       // Удален
}

type ViewFullOrder struct {
	Creator_ID int64 `json:"userId"` // Идентификатор пользователя
	Unit_ID    int64 `json:"unitId"` // Идентификатор объединения
	ViewLongOrder
}

type ApiMetaOrder struct {
	Total         int64 `json:"count"`   // Общее число заказов
	NumOfNew      int64 `json:"new"`     // Число новых заказов
	NumOfOpen     int64 `json:"open"`    // Число заказов в работе
	NumOfClosed   int64 `json:"close"`   // Число закрытых заказов
	NumOfArchived int64 `json:"archive"` // Число архивных заказов
	NumOfAlert    int64 `json:"alert"`   // Число заказов с уведомлениями
}

type ApiFullMetaOrder struct {
	Total            int64 `json:"count"`       // Общее число заказов
	NumOfCompleted   int64 `json:"completed"`   // Число оформленных заказов
	NumOfNew         int64 `json:"new"`         // Число новых заказов
	NumOfOpen        int64 `json:"open"`        // Число заказов в работе
	NumOfClosed      int64 `json:"close"`       // Число закрытых заказов
	NumOfNotPaid     int64 `json:"notPaid"`     // Число неоплаченных заказов
	NumOfOnTheGo     int64 `json:"onTheGo"`     // Число активных заказов
	NumOfNoDocuments int64 `json:"noDocuments"` // Число заказов без документов
	NumOfArchived    int64 `json:"archive"`     // Число архивированных заказов
	NumOfAlert       int64 `json:"alert"`       // Число заказов с предупреждениями
	NumOfDeleted     int64 `json:"deleted"`     // Число удаленных заказов
}

type ApiMetaOrderByProject struct {
	Total      int64 `json:"count"` // Общее число заказов в проекте
	NumOfAlert int64 `json:"alert"` // Число заказов с уведомлениями
}

type ApiTinyOrder struct {
	ID          int64 `json:"id" db:"id"`           // Уникальный идентификатор
	Facility_ID int64 `json:"type" db:"service_id"` // Идентификатор услуги
}

type ApiShortOrder struct {
	ID          int64  `json:"id" db:"id"`                 // Уникальный идентификатор
	IsAssembled bool   `json:"completed" db:"completed"`   // Собран
	Facility_ID int64  `json:"type" db:"type"`             // Идентификатор услуги
	Supplier_ID int64  `json:"supplierId" db:"supplierId"` // Идентификатор поставщика
	IsNew       bool   `json:"new" db:"new"`               // Новый
	IsOpen      bool   `json:"open" db:"open"`             // Открыт
	Name        string `json:"name" db:"name"`             // Название
}

type ApiMiddleOrder struct {
	ID                 int64   `json:"id" db:"id"`                                             // Уникальный идентификатор
	Name               string  `json:"name" db:"name"`                                         // Название
	Step               byte    `json:"step" db:"step"`                                         // Шаг
	IsAssembled        bool    `json:"completed" db:"completed"`                               // Собран
	IsConfirmed        bool    `json:"moderatorConfirmed" db:"moderatorConfirmed"`             // Утвержден
	Facility_ID        int64   `json:"type" db:"type"`                                         // Идентификатор услуги
	Supplier_ID        int64   `json:"supplierId" db:"supplierId"`                             // Идентификатор поставщика
	IsNew              bool    `json:"new" db:"new"`                                           // Новый
	IsOpen             bool    `json:"open" db:"open"`                                         // Открыт
	IsCancelled        bool    `json:"cancel" db:"cancel"`                                     // Отказ
	Reason             string  `json:"cancelDescription" db:"cancelDescription"`               // Причина отказа
	Execution_Forecast int     `json:"supplierForecastWorkDays" db:"supplierForecastWorkDays"` // Прогноз исполнения
	Proposed_Price     float64 `json:"supplierCost" db:"supplierCost"`                         // Предложенная цена
	IsNewCost          bool    `json:"supplierCostNew" db:"supplierCostNew"`                   // Новая цена
	IsNewCostConfirmed bool    `json:"customerNewCostConfirmed" db:"customerNewCostConfirmed"` // Новая цена утверждена
	IsPaid             bool    `json:"paid" db:"paid"`                                         // Оплачен
	IsStarted          bool    `json:"moderatorBegin" db:"moderatorBegin"`                     // Запущен
	IsNewCharge        bool    `json:"supplierFactualCostNew" db:"supplierFactualCostNew"`     // Предложена фактическая цена
	Charged_Fee        float64 `json:"supplierFactualCost" db:"supplierFactualCost"`           // Фактическая цена
	IsExecuted         bool    `json:"supplierClose" db:"supplierClose"`                       // Выполнен
	IsDocumented       bool    `json:"moderatorDocumentsGotten" db:"moderatorDocumentsGotten"` // Документы имеются
	IsClosed           bool    `json:"moderatorClose" db:"moderatorClose"`                     // Закрыт
	IsArchived         bool    `json:"archive" db:"archive"`                                   // Архивирован
}

type ApiBriefOrder struct {
	ID         int64  `json:"id" db:"id"`           // Уникальный идентификатор
	Name       string `json:"name" db:"name"`       // Название
	IsPaid     bool   `json:"paid" db:"paid"`       // Оплачен
	IsArchived bool   `json:"archive" db:"archive"` // Архивирован
	IsDeleted  bool   `json:"del" db:"del"`         // Удален
}

type ApiLongOrder struct {
	ID                 int64   `json:"id"`                       // Уникальный идентификатор
	Name               string  `json:"name"`                     // Название
	Step               byte    `json:"step"`                     // Шаг
	IsAssembled        bool    `json:"completed"`                // Собран
	IsConfirmed        bool    `json:"moderatorConfirmed"`       // Утвержден
	Facility_ID        int64   `json:"type"`                     // Идентификатор услуги
	Supplier_ID        int64   `json:"supplierId"`               // Идентификатор поставщика
	IsNew              bool    `json:"new"`                      // Новый
	IsOpen             bool    `json:"open"`                     // Открыт
	IsCancelled        bool    `json:"cancel"`                   // Отказ
	Reason             string  `json:"cancelDescription"`        // Причина отказа
	Execution_Forecast int     `json:"supplierForecastWorkDays"` // Прогноз исполнения
	Proposed_Price     float64 `json:"supplierCost"`             // Предложенная цена
	IsNewCost          bool    `json:"supplierCostNew"`          // Новая цена
	IsNewCostConfirmed bool    `json:"customerNewCostConfirmed"` // Новая цена утверждена
	IsPaid             bool    `json:"paid"`                     // Оплачен
	IsStarted          bool    `json:"moderatorBegin"`           // Запущен
	IsNewCharge        bool    `json:"supplierFactualCostNew"`   // Предложена фактическая цена
	Charged_Fee        float64 `json:"supplierFactualCost"`      // Фактическая цена
	IsExecuted         bool    `json:"supplierClose"`            // Выполнен
	IsDocumented       bool    `json:"moderatorDocumentsGotten"` // Документы имеются
	IsClosed           bool    `json:"moderatorClose"`           // Закрыт
	IsArchived         bool    `json:"archive"`                  // Архивирован
	IsDeleted          bool    `json:"del"`                      // Удален
}

type ApiListOrder struct {
	ID          int64   `json:"id" db:"id"`                 // Уникальный идентификатор
	Step        byte    `json:"step" db:"step"`             // Шаг
	IsAssembled bool    `json:"completed" db:"completed"`   // Собран
	Facility_ID int64   `json:"type" db:"type"`             // Идентификатор услуги
	Unit_ID     int64   `json:"unitId" db:"unitId"`         // Идентификатор объединения
	User_ID     int64   `json:"customerId" db:"customerId"` // Идентификатор пользователя
	Supplier_ID int64   `json:"supplierId" db:"supplierId"` // Идентификатор поставщика
	IsNew       bool    `json:"new" db:"new"`               // Новый
	IsOpen      bool    `json:"open" db:"open"`             // Открыт
	IsCancelled bool    `json:"cancel" db:"cancel"`         // Отказ
	Charged_Fee float64 `json:"cost" db:"cost"`             // Фактическая цена
	IsPaid      bool    `json:"paid" db:"paid"`             // Оплачен
	Name        string  `json:"name" db:"name"`             // Название
	IsArchived  bool    `json:"archive" db:"archive"`       // Архивирован
	IsDeleted   bool    `json:"del" db:"del"`               // Удален
}

type ApiFullOrder struct {
	User_ID int64     `json:"userId"`  // Идентификатор пользователя
	Unit_ID int64     `json:"unitId"`  // Идентификатор объединения
	Created time.Time `json:"created"` // Время создания
	ApiLongOrder
}

type ApiFinanceOrder struct {
	ID               int64     `json:"orderId" db:"orderId"`             // Уникальный идентификатор
	Project_ID       int64     `json:"projectId" db:"projectId"`         // Идентификатор проекта
	Begin_Date       time.Time `json:"beginDate" db:"beginDate"`         // Дата начала работ
	Charged_Fee      float64   `json:"cost" db:"cost"`                   // Фактическая цена
	Facility_ID      int64     `json:"type" db:"type"`                   // Идентификатор услуги
	Supplier_ID      int64     `json:"supplierId" db:"supplierId"`       // Идентификатор поставщика
	Status_ID        int       `json:"statusId" db:"statusId"`           // Статус заказа
	Act_File_ID      int64     `json:"actFileId" db:"actFileId"`         // Идентификатор файла акта
	EInvoice_File_ID int64     `json:"invoiceFileId" db:"invoiceFileId"` // Идентификатор файла акта
}

type ApiResultOrder struct {
	Total       int64   `json:"count" db:"count"`             // Всего заказов
	Processing  int64   `json:"countInWork" db:"countInWork"` // Заказов в работе
	Charged_Fee float64 `json:"cost" db:"cost"`               // Фактическая цена
}

type OrderFinanceSearch struct {
	ID               int64     `query:"orderId" search:"o.id"`                                                    // Уникальный идентификатор
	Project_ID       int64     `query:"projectId" search:"o.project_id"`                                          // Идентификатор проекта
	Begin_Date       time.Time `query:"beginDate" search:"o.begin_date" group:"convert(o.begin_date using utf8)"` // Дата начала работ
	Charged_Fee      float64   `query:"cost" search:"o.charged_fee" group:"o.charged_fee"`                        // Фактическая цена
	Facility_ID      int64     `query:"type" search:"o.service_id"`                                               // Идентификатор услуги
	Supplier_ID      int64     `query:"supplierId" search:"o.supplier_id"`                                        // Идентификатор поставщика
	Status_ID        int       `query:"statusId" search:"coalesce(s.complex_status_id, 0)"`                       // Статус заказа
	Act_File_ID      int64     `query:"actFileId" search:"coalesce(a.file_id, 0)"`                                // Идентификатор файла акта
	EInvoice_File_ID int64     `query:"invoiceFileId" search:"coalesce(i.file_id, 0)"`                            // Идентификатор файла акта
}

type OrderSearch struct {
	ID          int64  `query:"id" search:"o.id"`                        // Уникальный идентификатор
	IsAssembled bool   `query:"completed" search:"coalesce(c.value, 0)"` // Собран
	Facility_ID int64  `query:"type" search:"o.service_id"`              // Идентификатор услуги
	Supplier_ID int64  `query:"supplierId" search:"o.supplier_id"`       // Идентификатор поставщика
	IsNew       bool   `query:"new" search:"coalesce(n.value, 0)"`       // Новый
	IsOpen      bool   `query:"open" search:"coalesce(p.value, 0)"`      // Открыт
	Name        string `query:"name" search:"o.name" group:"o.name"`     // Название
}

type OrderProjectSearch struct {
	ID                 int64   `query:"id" search:"o.id"`                                                                     // Уникальный идентификатор
	Name               string  `query:"name" search:"o.name" group:"o.name"`                                                  // Название
	Step               byte    `query:"step" search:"o.step"`                                                                 // Шаг
	IsAssembled        bool    `query:"completed" search:"coalesce(c.value, 0)"`                                              // Собран
	IsConfirmed        bool    `query:"moderatorConfirmed" search:"coalesce(m.value, 0)"`                                     // Утвержден
	Facility_ID        int64   `query:"type" search:"o.service_id"`                                                           // Идентификатор услуги
	Supplier_ID        int64   `query:"supplierId" search:"o.supplier_id"`                                                    // Идентификатор поставщика
	IsNew              bool    `query:"new" search:"coalesce(n.value, 0)"`                                                    // Новый
	IsOpen             bool    `query:"open" search:"coalesce(v.value, 0)"`                                                   // Открыт
	IsCancelled        bool    `query:"cancel" search:"coalesce(a.value, 0)"`                                                 // Отказ
	Reason             string  `query:"cancelDescription" search:"coalesce(a.comments, '')" group:"coalesce(a.comments, '')"` // Причина отказа
	Execution_Forecast int     `query:"supplierForecastWorkDays" search:"o.execution_forecast" group:"o.execution_forecast"`  // Прогноз исполнения
	Proposed_Price     float64 `query:"supplierCost" search:"o.proposed_price" group:"o.proposed_price"`                      // Предложенная цена
	IsNewCost          bool    `query:"supplierCostNew" search:"coalesce(s.value, 0)"`                                        // Новая цена
	IsNewCostConfirmed bool    `query:"customerNewCostConfirmed" search:"coalesce(u.value, 0)"`                               // Новая цена утверждена
	IsPaid             bool    `query:"paid" search:"coalesce(p.value, 0)"`                                                   // Оплачен
	IsStarted          bool    `query:"moderatorBegin" search:"coalesce(d.value, 0)"`                                         // Запущен
	IsNewCharge        bool    `query:"supplierFactualCostNew" search:"coalesce(l.value, 0)"`                                 // Предложена фактическая цена
	Charged_Fee        float64 `query:"supplierFactualCost" search:"o.charged_fee" group:"o.charged_fee"`                     // Фактическая цена
	IsExecuted         bool    `query:"supplierClose" search:"coalesce(i.value, 0)"`                                          // Выполнен
	IsDocumented       bool    `query:"moderatorDocumentsGotten" search:"coalesce(e.value, 0)"`                               // Документы имеются
	IsClosed           bool    `query:"moderatorClose" search:"coalesce(r.value, 0)"`                                         // Закрыт
	IsArchived         bool    `query:"archive" search:"coalesce(h.value, 0)"`                                                // Архивирован
}

type OrderAdminSearch struct {
	ID          int64   `query:"id" search:"o.id"`                                  // Уникальный идентификатор
	Step        byte    `query:"step" search:"o.step"`                              // Шаг
	IsAssembled bool    `query:"completed" search:"coalesce(c.value, 0)"`           // Собран
	Facility_ID int64   `query:"type" search:"o.service_id"`                        // Идентификатор услуги
	Unit_ID     int64   `query:"unitId" search:"o.unit_id"`                         // Идентификатор объединения
	Creator_ID  int64   `query:"customerId" search:"o.user_id"`                     // Идентификатор создателя
	Supplier_ID int64   `query:"supplierId" search:"o.supplier_id"`                 // Идентификатор поставщика
	IsNew       bool    `query:"new" search:"coalesce(n.value, 0)"`                 // Новый
	IsOpen      bool    `query:"open" search:"coalesce(p.value, 0)"`                // Открыт
	IsCancelled bool    `query:"cancel" search:"coalesce(a.value, 0)"`              // Отказ
	Charged_Fee float64 `query:"cost" search:"o.charged_fee" group:"o.charged_fee"` // Фактическая цена
	IsPaid      bool    `query:"paid" search:"coalesce(i.value, 0)"`                // Оплачен
	Name        string  `query:"name" search:"o.name" group:"o.name"`               // Название
	IsArchived  bool    `query:"archive" search:"coalesce(r.value, 0)"`             // Архивирован
	IsDeleted   bool    `query:"del" search:"coalesce(e.value, 0)"`                 // Удален
}

type DtoOrder struct {
	ID                 int64     `db:"id"`                 // Уникальный идентификатор
	Project_ID         int64     `db:"project_id"`         // Идентификатор проекта
	Creator_ID         int64     `db:"user_id"`            // Идентификатор создателя
	Unit_ID            int64     `db:"unit_id"`            // Идентификатор объединения
	Supplier_ID        int64     `db:"supplier_id"`        // Идентификатор поставщика
	Facility_ID        int64     `db:"service_id"`         // Идентификатор услуги
	Name               string    `db:"name"`               // Название
	Step               byte      `db:"step"`               // Шаг
	Created            time.Time `db:"created"`            // Время создания
	Proposed_Price     float64   `db:"proposed_price"`     // Предложенная цена
	Charged_Fee        float64   `db:"charged_fee"`        // Фактическая цена
	Execution_Forecast int       `db:"execution_forecast"` // Прогноз исполнения
	Begin_Date         time.Time `db:"begin_date"`         // Дата начала работ
	End_Date           time.Time `db:"end_date"`           // Дата окончания работ
	Act_ID             int64     `db:"act_id"`             // Идентификатор акта
	EInvoice_ID        int64     `db:"eInvoice_id"`        // Идентификатор счета-фактуры
}

// Конструктор создания объекта заказа в api
func NewApiMetaOrder(total int64, numofnew int64, numofopen int64, numofclosed int64, numofarchived int64,
	numofalert int64) *ApiMetaOrder {
	return &ApiMetaOrder{
		Total:         total,
		NumOfNew:      numofnew,
		NumOfOpen:     numofopen,
		NumOfClosed:   numofclosed,
		NumOfArchived: numofarchived,
		NumOfAlert:    numofalert,
	}
}

func NewApiFullMetaOrder(total int64, numofcompleted int64, numofnew int64, numofopen int64, numofclosed int64, numofnotpaid int64,
	numofonthego int64, numofnodocuments int64, numofarchived int64, numofalert int64, numofdeleted int64) *ApiFullMetaOrder {
	return &ApiFullMetaOrder{
		Total:            total,
		NumOfCompleted:   numofcompleted,
		NumOfNew:         numofnew,
		NumOfOpen:        numofopen,
		NumOfClosed:      numofclosed,
		NumOfNotPaid:     numofnotpaid,
		NumOfOnTheGo:     numofonthego,
		NumOfNoDocuments: numofnodocuments,
		NumOfArchived:    numofarchived,
		NumOfAlert:       numofalert,
		NumOfDeleted:     numofdeleted,
	}
}

func NewApiMetaOrderByProject(total int64, numofalert int64) *ApiMetaOrderByProject {
	return &ApiMetaOrderByProject{
		Total:      total,
		NumOfAlert: numofalert,
	}
}

func NewApiTinyOrder(id int64, facility_id int64) *ApiTinyOrder {
	return &ApiTinyOrder{
		ID:          id,
		Facility_ID: facility_id,
	}
}

func NewApiShortOrder(id int64, isassembled bool, facility_id int64, supplier_id int64, isnew bool, isopen bool,
	name string) *ApiShortOrder {
	return &ApiShortOrder{
		ID:          id,
		IsAssembled: isassembled,
		Facility_ID: facility_id,
		Supplier_ID: supplier_id,
		IsNew:       isnew,
		IsOpen:      isopen,
		Name:        name,
	}
}

func NewApiMiddleOrder(id int64, isassembled bool, facility_id int64, supplier_id int64, isnew bool, isopen bool,
	iscancelled bool, ispaid bool, name string) *ApiMiddleOrder {
	return &ApiMiddleOrder{
		ID:          id,
		IsAssembled: isassembled,
		Facility_ID: facility_id,
		Supplier_ID: supplier_id,
		IsNew:       isnew,
		IsOpen:      isopen,
		IsCancelled: iscancelled,
		IsPaid:      ispaid,
		Name:        name,
	}
}

func NewApiBriefOrder(id int64, name string, ispaid bool, isarchived bool, isdeleted bool) *ApiBriefOrder {
	return &ApiBriefOrder{
		ID:         id,
		Name:       name,
		IsPaid:     ispaid,
		IsArchived: isarchived,
		IsDeleted:  isdeleted,
	}
}

func NewApiLongOrder(id int64, name string, step byte, isassembled bool, isconfirmed bool, facility_id int64,
	supplier_id int64, isnew bool, isopen bool, iscancelled bool, reason string, execution_forecast int, proposed_price float64,
	isnewcost bool, isnewcostconfirmed bool, ispaid bool, isstarted bool, isnewcharge bool, charged_fee float64, isexecuted bool,
	isdocumented bool, isclosed bool, isarchived bool, isdeleted bool) *ApiLongOrder {
	return &ApiLongOrder{
		ID:                 id,
		Name:               name,
		Step:               step,
		IsAssembled:        isassembled,
		IsConfirmed:        isconfirmed,
		Facility_ID:        facility_id,
		Supplier_ID:        supplier_id,
		IsNew:              isnew,
		IsOpen:             isopen,
		IsCancelled:        iscancelled,
		Reason:             reason,
		Execution_Forecast: execution_forecast,
		Proposed_Price:     proposed_price,
		IsNewCost:          isnewcost,
		IsNewCostConfirmed: isnewcostconfirmed,
		IsPaid:             ispaid,
		IsStarted:          isstarted,
		IsNewCharge:        isnewcharge,
		Charged_Fee:        charged_fee,
		IsExecuted:         isexecuted,
		IsDocumented:       isdocumented,
		IsClosed:           isclosed,
		IsArchived:         isarchived,
		IsDeleted:          isdeleted,
	}
}

func NewApiListOrder(id int64, step byte, isassembled bool, facility_id int64, unit_id int64, user_id int64, supplier_id int64,
	isnew bool, isopen bool, iscancelled bool, charged_fee float64, ispaid bool, name string, isarchived bool, isdeleted bool) *ApiListOrder {
	return &ApiListOrder{
		ID:          id,
		Step:        step,
		IsAssembled: isassembled,
		Facility_ID: facility_id,
		Unit_ID:     unit_id,
		User_ID:     user_id,
		Supplier_ID: supplier_id,
		IsNew:       isnew,
		IsOpen:      isopen,
		IsCancelled: iscancelled,
		Charged_Fee: charged_fee,
		IsPaid:      ispaid,
		Name:        name,
		IsArchived:  isarchived,
		IsDeleted:   isdeleted,
	}
}

func NewApiFullOrder(user_id int64, unit_id int64, created time.Time, apilongorder ApiLongOrder) *ApiFullOrder {
	return &ApiFullOrder{
		User_ID:      user_id,
		Unit_ID:      unit_id,
		Created:      created,
		ApiLongOrder: apilongorder,
	}
}

func NewApiFinanceOrder(id int64, project_id int64, begin_date time.Time, charged_fee float64, facility_id int64, supplier_id int64,
	status_id int, act_file_id int64, einvoice_file_id int64) *ApiFinanceOrder {
	return &ApiFinanceOrder{
		ID:               id,
		Project_ID:       project_id,
		Begin_Date:       begin_date,
		Charged_Fee:      charged_fee,
		Facility_ID:      facility_id,
		Supplier_ID:      supplier_id,
		Status_ID:        status_id,
		Act_File_ID:      act_file_id,
		EInvoice_File_ID: einvoice_file_id,
	}
}

func NewApiResultOrder(total int64, processing int64, charged_fee float64) *ApiResultOrder {
	return &ApiResultOrder{
		Total:       total,
		Processing:  processing,
		Charged_Fee: charged_fee,
	}
}

// Конструктор создания объекта заказа в бд
func NewDtoOrder(id int64, project_id int64, creator_id int64, unit_id int64, supplier_id int64, facility_id int64,
	name string, step byte, created time.Time, proposed_price float64, charged_fee float64, execution_forecast int,
	begin_date time.Time, end_date time.Time, act_id int64, einvoice_id int64) *DtoOrder {
	return &DtoOrder{
		ID:                 id,
		Project_ID:         project_id,
		Creator_ID:         creator_id,
		Unit_ID:            unit_id,
		Supplier_ID:        supplier_id,
		Facility_ID:        facility_id,
		Name:               name,
		Step:               step,
		Created:            created,
		Proposed_Price:     proposed_price,
		Charged_Fee:        charged_fee,
		Execution_Forecast: execution_forecast,
		Begin_Date:         begin_date,
		End_Date:           end_date,
		Act_ID:             act_id,
		EInvoice_ID:        einvoice_id,
	}
}

func (order *OrderSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, order), nil
}

func (order *OrderSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, order)
	errField = nil
	errValue = nil
	switch infield {
	case "id":
		fallthrough
	case "type":
		fallthrough
	case "supplierId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
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
	case "completed":
		fallthrough
	case "new":
		fallthrough
	case "open":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (order *OrderSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(order)
}

func (order *OrderProjectSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, order), nil
}

func (order *OrderProjectSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, order)
	errField = nil
	errValue = nil
	switch infield {
	case "id":
		fallthrough
	case "type":
		fallthrough
	case "supplierId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "supplierForecastWorkDays":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "step":
		_, errConv := strconv.ParseInt(invalue, 0, 8)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "supplierCost":
		fallthrough
	case "supplierFactualCost":
		_, errConv := strconv.ParseFloat(invalue, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "name":
		fallthrough
	case "cancelDescription":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	case "completed":
		fallthrough
	case "moderatorConfirmed":
		fallthrough
	case "new":
		fallthrough
	case "open":
		fallthrough
	case "cancel":
		fallthrough
	case "supplierCostNew":
		fallthrough
	case "customerNewCostConfirmed":
		fallthrough
	case "paid":
		fallthrough
	case "moderatorBegin":
		fallthrough
	case "supplierFactualCostNew":
		fallthrough
	case "supplierClose":
		fallthrough
	case "moderatorDocumentsGotten":
		fallthrough
	case "moderatorClose":
		fallthrough
	case "archive":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (order *OrderProjectSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(order)
}

func (order *ViewMiddleOrder) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(order, errors, req)
}

func (order *ViewLongOrder) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(order, errors, req)
}

func (order *ViewShortOrder) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(order, errors, req)
}

func (order *ViewLongOrder) ToOrderStatuses(id int64) (orderstatuses *[]DtoOrderStatus) {
	return &[]DtoOrderStatus{
		{Order_ID: id, Status_ID: ORDER_STATUS_COMPLETED, Value: order.IsAssembled, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_MODERATOR_CONFIRMED, Value: order.IsConfirmed, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_NEW, Value: order.IsNew, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_OPEN, Value: order.IsOpen, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_CANCEL, Value: order.IsCancelled, Comments: order.Reason, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_SUPPLIER_COST_NEW, Value: order.IsNewCost, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED, Value: order.IsNewCostConfirmed, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_PAID, Value: order.IsPaid, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_MODERATOR_BEGIN, Value: order.IsStarted, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_SUPPLIER_FACTUAL_COST_NEW, Value: order.IsNewCharge, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_SUPPLIER_CLOSE, Value: order.IsExecuted, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN, Value: order.IsDocumented, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_MODERATOR_CLOSE, Value: order.IsClosed, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_ARCHIVE, Value: order.IsArchived, Created: time.Now()},
		{Order_ID: id, Status_ID: ORDER_STATUS_DEL, Value: order.IsDeleted, Created: time.Now()},
	}
}

func NewApiLongOrderFromDto(dtoorder *DtoOrder, dtoorderstatuses *[]DtoOrderStatus) (apiorder *ApiLongOrder) {
	apiorder = new(ApiLongOrder)
	apiorder.ID = dtoorder.ID
	apiorder.Name = dtoorder.Name
	apiorder.Step = dtoorder.Step
	apiorder.Facility_ID = dtoorder.Facility_ID
	apiorder.Supplier_ID = dtoorder.Supplier_ID
	apiorder.Proposed_Price = dtoorder.Proposed_Price
	apiorder.Charged_Fee = dtoorder.Charged_Fee
	apiorder.Execution_Forecast = dtoorder.Execution_Forecast
	for _, dtoorderstatus := range *dtoorderstatuses {
		switch dtoorderstatus.Status_ID {
		case ORDER_STATUS_COMPLETED:
			apiorder.IsAssembled = dtoorderstatus.Value
		case ORDER_STATUS_MODERATOR_CONFIRMED:
			apiorder.IsConfirmed = dtoorderstatus.Value
		case ORDER_STATUS_NEW:
			apiorder.IsNew = dtoorderstatus.Value
		case ORDER_STATUS_OPEN:
			apiorder.IsOpen = dtoorderstatus.Value
		case ORDER_STATUS_CANCEL:
			apiorder.IsCancelled = dtoorderstatus.Value
			apiorder.Reason = dtoorderstatus.Comments
		case ORDER_STATUS_SUPPLIER_COST_NEW:
			apiorder.IsNewCost = dtoorderstatus.Value
		case ORDER_STATUS_CUSTOMER_NEW_COST_CONFIRMED:
			apiorder.IsNewCostConfirmed = dtoorderstatus.Value
		case ORDER_STATUS_PAID:
			apiorder.IsPaid = dtoorderstatus.Value
		case ORDER_STATUS_MODERATOR_BEGIN:
			apiorder.IsStarted = dtoorderstatus.Value
		case ORDER_STATUS_SUPPLIER_FACTUAL_COST_NEW:
			apiorder.IsNewCharge = dtoorderstatus.Value
		case ORDER_STATUS_SUPPLIER_CLOSE:
			apiorder.IsExecuted = dtoorderstatus.Value
		case ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN:
			apiorder.IsDocumented = dtoorderstatus.Value
		case ORDER_STATUS_MODERATOR_CLOSE:
			apiorder.IsClosed = dtoorderstatus.Value
		case ORDER_STATUS_ARCHIVE:
			apiorder.IsArchived = dtoorderstatus.Value
		case ORDER_STATUS_DEL:
			apiorder.IsDeleted = dtoorderstatus.Value
		}
	}

	return apiorder
}

func NewApiFullOrderFromDto(dtoorder *DtoOrder, dtoorderstatuses *[]DtoOrderStatus) (apiorder *ApiFullOrder) {
	apiorder = new(ApiFullOrder)
	apiorder.Unit_ID = dtoorder.Unit_ID
	apiorder.User_ID = dtoorder.Creator_ID
	apiorder.Created = dtoorder.Created
	apiorder.ApiLongOrder = *NewApiLongOrderFromDto(dtoorder, dtoorderstatuses)
	return apiorder
}

func (order *OrderAdminSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, order), nil
}

func (order *OrderAdminSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, order)
	errField = nil
	errValue = nil

	switch infield {
	case "id":
		fallthrough
	case "step":
		fallthrough
	case "type":
		fallthrough
	case "unitId":
		fallthrough
	case "customerId":
		fallthrough
	case "supplierId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "cost":
		_, errConv := strconv.ParseFloat(invalue, 64)
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
	case "completed":
		fallthrough
	case "new":
		fallthrough
	case "open":
		fallthrough
	case "cancel":
		fallthrough
	case "paid":
		fallthrough
	case "archive":
		fallthrough
	case "del":
		val, errConv := strconv.ParseBool(invalue)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = fmt.Sprintf("%v", val)
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (order *OrderAdminSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(order)
}

func (order *OrderFinanceSearch) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, order), nil
}

func (order *OrderFinanceSearch) Extract(infield string, invalue string) (outfield string, outvalue string, errField error, errValue error) {
	outvalue = ""
	outfield = GetSearchTag(infield, order)
	errField = nil
	errValue = nil

	switch infield {
	case "orderId":
		fallthrough
	case "projectId":
		fallthrough
	case "type":
		fallthrough
	case "supplierId":
		fallthrough
	case "actFileId":
		fallthrough
	case "invoiceFileId":
		_, errConv := strconv.ParseInt(invalue, 0, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "statusId":
		_, errConv := strconv.ParseInt(invalue, 0, 32)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "cost":
		_, errConv := strconv.ParseFloat(invalue, 64)
		if errConv != nil {
			errValue = errConv
			break
		}
		outvalue = invalue
	case "beginDate":
		if strings.Contains(invalue, "'") {
			invalue = strings.Replace(invalue, "'", "''", -1)
		}
		outvalue = "'" + invalue + "'"
	default:
		errField = errors.New("Unknown field")
	}

	return outfield, outvalue, errField, errValue
}

func (order *OrderFinanceSearch) GetAllFields(parameter interface{}) (fields *[]string) {
	return GetAllGroupTags(order)
}
