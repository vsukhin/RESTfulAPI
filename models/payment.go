package models

import (
	"github.com/martini-contrib/binding"
	"net/http"
	"time"
)

// Структура для организации хранения платежа
type ViewPayment struct {
	Tariff_Plan_ID int  `json:"planId"`             // Идентификатор тарифного плана
	Renew          bool `json:"autoNextPaymentDue"` // Продлевать
}

type ApiPayment struct {
	Tariff_Plan_ID   int       `json:"planId" db:"tariff_plan_id"`           // Идентификатор тарифного плана
	Paid             bool      `json:"paid" db:"paid"`                       // Оплачен
	Payment_Date     time.Time `json:"paymentDate" db:"payment_date"`        // Дата оплаты
	Next_Payment_Due time.Time `json:"nextPaymentDue" db:"next_payment_due"` // Следущая оплата
	Renew            bool      `json:"autoNextPaymentDue" db:"renew"`        // Продлевать
}

type DtoPayment struct {
	Unit_ID          int64     `db:"unit_id"`          // Идентификатор объединения
	Tariff_Plan_ID   int       `db:"tariff_plan_id"`   // Идентификатор тарифного плана
	Paid             bool      `db:"paid"`             // Оплачен
	Payment_Date     time.Time `db:"payment_date"`     // Дата оплаты
	Next_Payment_Due time.Time `db:"next_payment_due"` // Следущая оплата
	Renew            bool      `db:"renew"`            // Продлевать
}

// Конструктор создания объекта платежа в api
func NewApiPayment(tariff_plan_id int, paid bool, payment_date, next_payment_due time.Time, renew bool) *ApiPayment {
	return &ApiPayment{
		Tariff_Plan_ID:   tariff_plan_id,
		Paid:             paid,
		Payment_Date:     payment_date,
		Next_Payment_Due: next_payment_due,
		Renew:            renew,
	}
}

// Конструктор создания объекта платежа в бд
func NewDtoPayment(unit_id int64, tariff_plan_id int, paid bool, payment_date, next_payment_due time.Time, renew bool) *DtoPayment {
	return &DtoPayment{
		Unit_ID:          unit_id,
		Tariff_Plan_ID:   tariff_plan_id,
		Paid:             paid,
		Payment_Date:     payment_date,
		Next_Payment_Due: next_payment_due,
		Renew:            renew,
	}
}

func (payment *ViewPayment) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	return Validate(payment, errors, req)
}
