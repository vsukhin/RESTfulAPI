package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type InvoiceItemRepository interface {
	Get(id int64) (invoiceitem *models.DtoInvoiceItem, err error)
	GetByInvoice(invoice_id int64) (invoiceitems *[]models.ApiInvoiceItem, err error)
	Create(dtoinvoiceitem *models.DtoInvoiceItem, trans *gorp.Transaction) (err error)
	DeleteByInvoice(invoice_id int64, trans *gorp.Transaction) (err error)
}

type InvoiceItemService struct {
	*Repository
}

func NewInvoiceItemService(repository *Repository) *InvoiceItemService {
	repository.DbContext.AddTableWithName(models.DtoInvoiceItem{}, repository.Table).SetKeys(true, "id")
	return &InvoiceItemService{Repository: repository}
}

func (invoiceitemservice *InvoiceItemService) Get(id int64) (invoiceitem *models.DtoInvoiceItem, err error) {
	invoiceitem = new(models.DtoInvoiceItem)
	err = invoiceitemservice.DbContext.SelectOne(invoiceitem, "select * from "+invoiceitemservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting invoice item object from database %v with value %v", err, id)
		return nil, err
	}

	return invoiceitem, nil
}

func (invoiceitemservice *InvoiceItemService) GetByInvoice(invoice_id int64) (invoiceitems *[]models.ApiInvoiceItem, err error) {
	invoiceitems = new([]models.ApiInvoiceItem)
	_, err = invoiceitemservice.DbContext.Select(invoiceitems,
		"select id, name, measure, amount, price, total from "+invoiceitemservice.Table+" where invoice_id = ?", invoice_id)
	if err != nil {
		log.Error("Error during getting all invoice item object from database %v with value %v", err, invoice_id)
		return nil, err
	}

	return invoiceitems, nil
}

func (invoiceitemservice *InvoiceItemService) Create(dtoinvoiceitem *models.DtoInvoiceItem, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoinvoiceitem)
	} else {
		err = invoiceitemservice.DbContext.Insert(dtoinvoiceitem)
	}
	if err != nil {
		log.Error("Error during creating invoice item object in database %v", err)
		return err
	}

	return nil
}

func (invoiceitemservice *InvoiceItemService) DeleteByInvoice(invoice_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+invoiceitemservice.Table+" where invoice_id = ?", invoice_id)
	} else {
		_, err = invoiceitemservice.DbContext.Exec("delete from "+invoiceitemservice.Table+" where invoice_id = ?", invoice_id)
	}
	if err != nil {
		log.Error("Error during deleting invoice item objects for invoice object in database %v with value %v", err, invoice_id)
		return err
	}

	return nil
}
