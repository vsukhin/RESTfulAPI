package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type DataProductRepository interface {
	Get(order_id int64, product_id int) (dataproduct *models.DtoDataProduct, err error)
	GetByOrder(order_id int64) (dataproducts *[]models.ViewApiDataProduct, err error)
	Create(dtodataproduct *models.DtoDataProduct, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type DataProductService struct {
	*Repository
}

func NewDataProductService(repository *Repository) *DataProductService {
	repository.DbContext.AddTableWithName(models.DtoDataProduct{}, repository.Table).SetKeys(false, "order_id", "product_id")
	return &DataProductService{Repository: repository}
}

func (dataproductservice *DataProductService) Get(order_id int64, product_id int) (dataproduct *models.DtoDataProduct, err error) {
	dataproduct = new(models.DtoDataProduct)
	err = dataproductservice.DbContext.SelectOne(dataproduct, "select * from "+dataproductservice.Table+
		" where order_id = ? and product_id = ?", order_id, product_id)
	if err != nil {
		log.Error("Error during getting data product object from database %v with value %v, %v", err, order_id, product_id)
		return nil, err
	}

	return dataproduct, nil
}

func (dataproductservice *DataProductService) GetByOrder(order_id int64) (dataproducts *[]models.ViewApiDataProduct, err error) {
	dataproducts = new([]models.ViewApiDataProduct)
	_, err = dataproductservice.DbContext.Select(dataproducts,
		"select product_id from "+dataproductservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all data product object from database %v with value %v", err, order_id)
		return nil, err
	}

	return dataproducts, nil
}

func (dataproductservice *DataProductService) Create(dtodataproduct *models.DtoDataProduct, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtodataproduct)
	} else {
		err = dataproductservice.DbContext.Insert(dtodataproduct)
	}
	if err != nil {
		log.Error("Error during creating data product object in database %v", err)
		return err
	}

	return nil
}

func (dataproductservice *DataProductService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+dataproductservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = dataproductservice.DbContext.Exec("delete from "+dataproductservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting data product objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
