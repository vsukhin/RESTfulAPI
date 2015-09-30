package services

import (
	"application/models"
)

type FinanceRepository interface {
	Get(unit_id int64) (finance *models.ApiFinance, err error)
}

type FinanceService struct {
	OperationRepository OperationRepository
	*Repository
}

func NewFinanceService(repository *Repository) *FinanceService {
	return &FinanceService{Repository: repository}
}

func (financeservice *FinanceService) Get(unit_id int64) (finance *models.ApiFinance, err error) {
	finance = new(models.ApiFinance)
	finance.Balance, err = financeservice.OperationRepository.CalculateBalance(unit_id)
	if err != nil {
		return nil, err
	}
	finance.Balance = models.Round(finance.Balance, 0.5, 2)
	finance.TotalInvoiceAll, err = financeservice.DbContext.SelectFloat(
		"select coalesce(sum(total), 0) from invoices where company_id in (select id from companies where unit_id = ?) and active = 1", unit_id)
	if err != nil {
		log.Error("Error during getting finance object from database %v with value %v", err, unit_id)
		return nil, err
	}
	finance.TotalInvoiceAll = models.Round(finance.TotalInvoiceAll, 0.5, 2)
	finance.TotalInvoicePaid, err = financeservice.DbContext.SelectFloat(
		"select coalesce(sum(total), 0) from invoices where company_id in (select id from companies where unit_id = ?) and active = 1 and paid = 1", unit_id)
	if err != nil {
		log.Error("Error during getting finance object from database %v with value %v", err, unit_id)
		return nil, err
	}
	finance.TotalInvoicePaid = models.Round(finance.TotalInvoicePaid, 0.5, 2)
	finance.TotalOrderExecuted, err = financeservice.DbContext.SelectFloat("select coalesce(sum(charged_fee), 0) from orders o"+
		" inner join order_statuses s on o.id = s.order_id inner join order_statuses t on o.id = t.order_id"+
		" where (s.status_id = ? and s.value = 1) and (t.status_id = ? and t.value = 1)",
		models.ORDER_STATUS_SUPPLIER_CLOSE, models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN)
	if err != nil {
		log.Error("Error during getting finance object from database %v with value %v", err, unit_id)
		return nil, err
	}
	finance.TotalOrderExecuted = models.Round(finance.TotalOrderExecuted, 0.5, 2)
	finance.TotalOrderProcessing, err = financeservice.DbContext.SelectFloat("select coalesce(sum(charged_fee), 0) from orders o"+
		" inner join order_statuses s on o.id = s.order_id where s.status_id = ? and s.value = 1"+
		" and o.id not in (select order_id from order_statuses where status_id = ? and value = 1)",
		models.ORDER_STATUS_MODERATOR_BEGIN, models.ORDER_STATUS_MODERATOR_DOCUMENTS_GOTTEN)
	if err != nil {
		log.Error("Error during getting finance object from database %v with value %v", err, unit_id)
		return nil, err
	}
	finance.TotalOrderProcessing = models.Round(finance.TotalOrderProcessing, 0.5, 2)

	return finance, nil
}
