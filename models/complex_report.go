package models

type DoughnutValues struct {
	Count int64
	Sum   float64
}

// Структура для организации хранения составного отчета
type ApiDoughnutReport struct {
	ID       int64   `json:"id"`       // Идентификатор
	Position int64   `json:"position"` // Позиция
	Percent  float64 `json:"percent"`  // Процент
	Value    float64 `json:"number"`   // Значение
}

type ApiBudgetedReport struct {
	Budgeted string `json:"budgetBy"` // Вид бюджетирования
	ApiDoughnutReport
}

type ApiOrderReport struct {
	Project_ID         int64   `json:"projectId" db:"projectId" query:"projectId"`          // Идентификатор проекта
	Project_Name       string  `json:"projectName" db:"projectName" query:"projectName"`    // Название проекта
	Order_ID           int64   `json:"orderId" db:"orderId" query:"orderId"`                // Идентификатор заказа
	Order_Name         string  `json:"orderName" db:"orderName" query:"orderName"`          // Название заказа
	Facility_ID        int64   `json:"serviceId" db:"serviceId" query:"serviceId"`          // Идентификатор поставщика
	Facility_Name      string  `json:"serviceName" db:"serviceName" query:"serviceName"`    // Название поставщика
	Supplier_ID        int64   `json:"supplierId" db:"supplierId" query:"supplierId"`       // Идентификатор поставщика
	Supplier_Name      string  `json:"supplierName" db:"supplierName" query:"supplierName"` // Название поставщика
	ComplexStatus_ID   int64   `json:"statusId" db:"statusId" query:"statusId"`             // Идентификатор статуса
	ComplexStatus_Name string  `json:"statusName" db:"statusName" query:"statusName"`       // Название статуса
	Begin              string  `json:"dateBegin" db:"dateBegin" query:"dateBegin"`          // Начало
	End                string  `json:"dateEnd" db:"dateEnd" query:"dateEnd"`                // Окончание
	Budget             float64 `json:"budget" db:"budget" query:"budget"`                   // Бюджет
}

type ApiComplexReport struct {
	Facilities      []ApiDoughnutReport `json:"donutServices,omitempty"`  // Услуги
	ComplexStatuses []ApiDoughnutReport `json:"donutStatuses,omitempty"`  // Статусы
	Suppliers       []ApiDoughnutReport `json:"donutSuppliers,omitempty"` // Поставщики
	Budgets         []ApiBudgetedReport `json:"donutBudget,omitempty"`    // Бюджеты
	Orders          []ApiOrderReport    `json:"orders,omitempty"`         // Заказы
}

// Конструктор создания объекта составного отчета в api
func NewApiDoughnutReport(id int64, position int64, percent float64, value float64) *ApiDoughnutReport {
	return &ApiDoughnutReport{
		ID:       id,
		Position: position,
		Percent:  percent,
		Value:    value,
	}
}

func NewApiBudgetedReport(budgeted string, apidoughnutreport ApiDoughnutReport) *ApiBudgetedReport {
	return &ApiBudgetedReport{
		Budgeted:          budgeted,
		ApiDoughnutReport: apidoughnutreport,
	}
}

func NewApiOrderReport(project_id int64, project_name string, order_id int64, order_name string, facility_id int64, facility_name string,
	supplier_id int64, supplier_name string, complexstatus_id int64, complexstatus_name string, begin string, end string,
	budget float64) *ApiOrderReport {
	return &ApiOrderReport{
		Project_ID:         project_id,
		Project_Name:       project_name,
		Order_ID:           order_id,
		Order_Name:         order_name,
		Facility_ID:        facility_id,
		Facility_Name:      facility_name,
		Supplier_ID:        supplier_id,
		Supplier_Name:      supplier_name,
		ComplexStatus_ID:   complexstatus_id,
		ComplexStatus_Name: complexstatus_name,
		Begin:              begin,
		End:                end,
		Budget:             budget,
	}
}

func NewApiComplexReport(facilities []ApiDoughnutReport, complexstatuses []ApiDoughnutReport, suppliers []ApiDoughnutReport,
	budgets []ApiBudgetedReport, orders []ApiOrderReport) *ApiComplexReport {
	return &ApiComplexReport{
		Facilities:      facilities,
		ComplexStatuses: complexstatuses,
		Suppliers:       suppliers,
		Budgets:         budgets,
		Orders:          orders,
	}
}

func (orderreport *ApiOrderReport) Check(field string) (valid bool, err error) {
	return CheckQueryTag(field, orderreport), nil
}
