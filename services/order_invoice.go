package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type OrderInvoiceRepository interface {
	Get(id int64) (orderinvoice *models.DtoOrderInvoice, err error)
	GetByOrder(order_id int64) (orderinvoices *[]models.DtoOrderInvoice, err error)
	Create(dtoorderinvoice *models.DtoOrderInvoice, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type OrderInvoiceService struct {
	*Repository
}

func NewOrderInvoiceService(repository *Repository) *OrderInvoiceService {
	repository.DbContext.AddTableWithName(models.DtoOrderInvoice{}, repository.Table).SetKeys(true, "id")
	return &OrderInvoiceService{Repository: repository}
}

func (orderinvoiceservice *OrderInvoiceService) Get(id int64) (orderinvoice *models.DtoOrderInvoice, err error) {
	orderinvoice = new(models.DtoOrderInvoice)
	err = orderinvoiceservice.DbContext.SelectOne(orderinvoice, "select * from "+orderinvoiceservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting order invoice object from database %v with value %v", err, id)
		return nil, err
	}

	return orderinvoice, nil
}

func (orderinvoiceservice *OrderInvoiceService) GetByOrder(order_id int64) (orderinvoices *[]models.DtoOrderInvoice, err error) {
	orderinvoices = new([]models.DtoOrderInvoice)
	_, err = orderinvoiceservice.DbContext.Select(orderinvoices,
		"select i* from "+orderinvoiceservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all order invoice object from database %v with value %v", err, order_id)
		return nil, err
	}

	return orderinvoices, nil
}

func (orderinvoiceservice *OrderInvoiceService) Create(dtoorderinvoice *models.DtoOrderInvoice, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoorderinvoice)
	} else {
		err = orderinvoiceservice.DbContext.Insert(dtoorderinvoice)
	}
	if err != nil {
		log.Error("Error during creating order invoice object in database %v", err)
		return err
	}

	return nil
}

func (orderinvoiceservice *OrderInvoiceService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+orderinvoiceservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = orderinvoiceservice.DbContext.Exec("delete from "+orderinvoiceservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting order invoice objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
