package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type SupplierRequestRepository interface {
	Get(order_id int64, supplier_id int64) (supplierrequest *models.DtoSupplierRequest, err error)
	GetByOrder(order_id int64) (supplierrequests *[]models.ApiSupplierRequest, err error)
	Create(dtosupplierrequest *models.DtoSupplierRequest, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type SupplierRequestService struct {
	*Repository
}

func NewSupplierRequestService(repository *Repository) *SupplierRequestService {
	repository.DbContext.AddTableWithName(models.DtoSupplierRequest{}, repository.Table).SetKeys(false, "order_id", "supplier_id")
	return &SupplierRequestService{Repository: repository}
}

func (supplierrequestservice *SupplierRequestService) Get(order_id int64, supplier_id int64) (supplierrequest *models.DtoSupplierRequest, err error) {
	supplierrequest = new(models.DtoSupplierRequest)
	err = supplierrequestservice.DbContext.SelectOne(supplierrequest, "select * from "+supplierrequestservice.Table+
		" where order_id = ? and supplier_id = ?", order_id, supplier_id)
	if err != nil {
		log.Error("Error during getting supplier request object from database %v with value %v, %v", err, order_id, supplier_id)
		return nil, err
	}

	return supplierrequest, nil
}

func (supplierrequestservice *SupplierRequestService) GetByOrder(order_id int64) (supplierrequests *[]models.ApiSupplierRequest, err error) {
	supplierrequests = new([]models.ApiSupplierRequest)
	_, err = supplierrequestservice.DbContext.Select(supplierrequests,
		"select supplier_id, requestDate, responded, respondedDate, estimatedCost, myChoice from "+
			supplierrequestservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all supplier request object from database %v with value %v", err, order_id)
		return nil, err
	}

	return supplierrequests, nil
}

func (supplierrequestservice *SupplierRequestService) Create(dtosupplierrequest *models.DtoSupplierRequest, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtosupplierrequest)
	} else {
		err = supplierrequestservice.DbContext.Insert(dtosupplierrequest)
	}
	if err != nil {
		log.Error("Error during creating supplier request object in database %v", err)
		return err
	}

	return nil
}

func (supplierrequestservice *SupplierRequestService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+supplierrequestservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = supplierrequestservice.DbContext.Exec("delete from "+supplierrequestservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting supplier request objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
