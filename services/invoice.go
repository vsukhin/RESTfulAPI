package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
	"time"
)

type InvoiceRepository interface {
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (invoice *models.DtoInvoice, err error)
	GetMeta(user_id int64) (invoice *models.ApiMetaInvoice, err error)
	GetByUser(userid int64, filter string) (invoices *[]models.ApiShortInvoice, err error)
	GetByUnit(unitid int64, filter string) (invoices *[]models.ApiShortInvoice, err error)
	SetArrays(invoice *models.DtoInvoice, trans *gorp.Transaction) (err error)
	PaidForOrder(order_id int64, dtoinvoice *models.DtoInvoice, dtotransaction *models.DtoTransaction, inTrans bool) (err error)
	Create(invoice *models.DtoInvoice, trans *gorp.Transaction, inTrans bool) (err error)
	Update(invoice *models.DtoInvoice, trans *gorp.Transaction, inTrans bool) (err error)
	Deactivate(invoice *models.DtoInvoice) (err error)
}

type InvoiceService struct {
	InvoiceItemRepository  InvoiceItemRepository
	TransactionRepository  TransactionRepository
	OperationRepository    OperationRepository
	OrderInvoiceRepository OrderInvoiceRepository
	OrderStatusRepository  OrderStatusRepository

	*Repository
}

func NewInvoiceService(repository *Repository) *InvoiceService {
	repository.DbContext.AddTableWithName(models.DtoInvoice{}, repository.Table).SetKeys(true, "id")
	return &InvoiceService{Repository: repository}
}

func (invoiceservice *InvoiceService) CheckUserAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := invoiceservice.DbContext.SelectInt("select count(*) from "+invoiceservice.Table+
		" where id = ? and company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", id, user_id)
	if err != nil {
		log.Error("Error during checking invoice object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (invoiceservice *InvoiceService) Get(id int64) (invoice *models.DtoInvoice, err error) {
	invoice = new(models.DtoInvoice)
	err = invoiceservice.DbContext.SelectOne(invoice, "select * from "+invoiceservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting invoice object from database %v with value %v", err, id)
		return nil, err
	}

	return invoice, nil
}

func (invoiceservice *InvoiceService) GetMeta(user_id int64) (invoice *models.ApiMetaInvoice, err error) {
	invoice = new(models.ApiMetaInvoice)
	invoice.Total, err = invoiceservice.DbContext.SelectInt("select count(*) from "+invoiceservice.Table+
		" where company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting meta invoice object from database %v with value %v", err, user_id)
		return nil, err
	}
	invoice.Unpaid, err = invoiceservice.DbContext.SelectInt("select count(*) from "+invoiceservice.Table+
		" where paid = 0 and company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting meta invoice object from database %v with value %v", err, user_id)
		return nil, err
	}
	invoice.Companies, err = invoiceservice.DbContext.SelectInt("select count(distinct company_id) from "+invoiceservice.Table+
		" where company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting meta invoice object from database %v with value %v", err, user_id)
		return nil, err
	}
	invoice.Deleted, err = invoiceservice.DbContext.SelectInt("select count(*) from "+invoiceservice.Table+
		" where active = 0 and company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting meta invoice object from database %v with value %v", err, user_id)
		return nil, err
	}

	return invoice, nil
}

func (invoiceservice *InvoiceService) GetByUser(userid int64, filter string) (invoices *[]models.ApiShortInvoice, err error) {
	invoices = new([]models.ApiShortInvoice)
	_, err = invoiceservice.DbContext.Select(invoices,
		"select id, company_id as organisationId, total, paid, not active as del from "+invoiceservice.Table+
			" where company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit invoice object from database %v with value %v", err, userid)
		return nil, err
	}

	return invoices, nil
}

func (invoiceservice *InvoiceService) GetByUnit(unitid int64, filter string) (invoices *[]models.ApiShortInvoice, err error) {
	invoices = new([]models.ApiShortInvoice)
	_, err = invoiceservice.DbContext.Select(invoices,
		"select id, company_id as organisationId, total, paid, not active as del from "+invoiceservice.Table+
			" where company_id in (select id from companies where unit_id = ?)"+filter, unitid)
	if err != nil {
		log.Error("Error during getting unit invoice object from database %v with value %v", err, unitid)
		return nil, err
	}

	return invoices, nil
}

func (invoiceservice *InvoiceService) SetArrays(invoice *models.DtoInvoice, trans *gorp.Transaction) (err error) {
	err = invoiceservice.InvoiceItemRepository.DeleteByInvoice(invoice.ID, trans)
	if err != nil {
		log.Error("Error during setting invoice object in database %v with value %v", err, invoice.ID)
		return err
	}
	for _, dtoinvoiceitem := range invoice.InvoiceItems {
		dtoinvoiceitem.Invoice_ID = invoice.ID
		err = invoiceservice.InvoiceItemRepository.Create(&dtoinvoiceitem, trans)
		if err != nil {
			log.Error("Error during setting invoice object in database %v with value %v", err, invoice.ID)
			return err
		}
	}

	return nil
}

func (invoiceservice *InvoiceService) PaidForOrder(order_id int64, dtoinvoice *models.DtoInvoice, dtotransaction *models.DtoTransaction,
	inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = invoiceservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during paying invoice in database %v", err)
			return err
		}
	}

	err = invoiceservice.Create(dtoinvoice, trans, false)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	dtoorderinvoice := new(models.DtoOrderInvoice)
	dtoorderinvoice.Order_ID = order_id
	dtoorderinvoice.Invoice_ID = dtoinvoice.ID
	err = invoiceservice.OrderInvoiceRepository.Create(dtoorderinvoice, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	err = invoiceservice.TransactionRepository.Create(dtotransaction, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	dtocredit := new(models.DtoOperation)
	dtocredit.Transaction_ID = dtotransaction.ID
	dtocredit.Invoice_ID = dtoinvoice.ID
	dtocredit.Money = dtoinvoice.Total
	dtocredit.Type_ID = models.OPERATION_TYPE_WITHDRAW
	dtocredit.Created = time.Now()
	err = invoiceservice.OperationRepository.Create(dtocredit, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}
	dtodebet := new(models.DtoOperation)
	dtodebet.Transaction_ID = dtotransaction.ID
	dtodebet.Invoice_ID = dtoinvoice.ID
	dtodebet.Money = dtoinvoice.Total
	dtodebet.Type_ID = models.OPERATION_TYPE_RECEIVE
	dtodebet.Created = time.Now()
	err = invoiceservice.OperationRepository.Create(dtodebet, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	dtoorderstatus := models.NewDtoOrderStatus(order_id, models.ORDER_STATUS_PAID, true, "", time.Now())
	err = invoiceservice.OrderStatusRepository.Save(dtoorderstatus, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during paying invoice in database %v", err)
			return err
		}
	}

	return nil
}

func (invoiceservice *InvoiceService) Create(invoice *models.DtoInvoice, trans *gorp.Transaction, inTrans bool) (err error) {
	if inTrans {
		trans, err = invoiceservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating invoice object in database %v", err)
			return err
		}
	}

	if trans != nil {
		err = trans.Insert(invoice)
	} else {
		err = invoiceservice.DbContext.Insert(invoice)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating invoice object in database %v", err)
		return err
	}

	err = invoiceservice.SetArrays(invoice, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating invoice object in database %v", err)
			return err
		}
	}

	return nil
}

func (invoiceservice *InvoiceService) Update(invoice *models.DtoInvoice, trans *gorp.Transaction, inTrans bool) (err error) {
	if inTrans {
		trans, err = invoiceservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating invoice object in database %v", err)
			return err
		}
	}

	if trans != nil {
		_, err = trans.Update(invoice)
	} else {
		_, err = invoiceservice.DbContext.Update(invoice)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating invoice object in database %v with value %v", err, invoice.ID)
		return err
	}

	err = invoiceservice.SetArrays(invoice, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating invoice object in database %v", err)
			return err
		}
	}

	return nil
}

func (invoiceservice *InvoiceService) Deactivate(invoice *models.DtoInvoice) (err error) {
	_, err = invoiceservice.DbContext.Exec("update "+invoiceservice.Table+" set active = 0 where id = ?", invoice.ID)
	if err != nil {
		log.Error("Error during deactivating invoice object in database %v with value %v", err, invoice.ID)
		return err
	}

	return nil
}
