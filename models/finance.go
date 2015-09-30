package models

// Структура для организации хранения финансов
type ApiFinance struct {
	Balance              float64 `json:"amountTotal"`           // Сумма неизрасходованных средств
	TotalInvoiceAll      float64 `json:"amountInvoicesCreated"` // Сумма счетов
	TotalInvoicePaid     float64 `json:"amountInvoicesPaid"`    // Сумма оплаченных счетов
	TotalOrderExecuted   float64 `json:"amountOrdersPerformed"` // Сумма выполненных счетов
	TotalOrderProcessing float64 `json:"amountOrdersInWork"`    // Сумма заказов в работе
}

// Конструктор создания объекта финансов в api
func NewApiFinance(balance float64, totalinvoiceall float64, totalinvoicepaid float64, totalorderexecuted float64, totalorderprocessing float64) *ApiFinance {
	return &ApiFinance{
		Balance:              balance,
		TotalInvoiceAll:      totalinvoiceall,
		TotalInvoicePaid:     totalinvoicepaid,
		TotalOrderExecuted:   totalorderexecuted,
		TotalOrderProcessing: totalorderprocessing,
	}
}
