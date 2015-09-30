package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type InputProductRepository interface {
	Get(order_id int64, product_id int) (inputproduct *models.DtoInputProduct, err error)
	GetByOrder(order_id int64) (inputproducts *[]models.ViewApiInputProduct, err error)
	GetAll(order_id int64) (inputproducts *[]models.DtoInputProduct, err error)
	Create(dtoinputproduct *models.DtoInputProduct, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type InputProductService struct {
	*Repository
}

func NewInputProductService(repository *Repository) *InputProductService {
	repository.DbContext.AddTableWithName(models.DtoInputProduct{}, repository.Table).SetKeys(false, "order_id", "product_id")
	return &InputProductService{Repository: repository}
}

func (inputproductservice *InputProductService) Get(order_id int64, product_id int) (inputproduct *models.DtoInputProduct, err error) {
	inputproduct = new(models.DtoInputProduct)
	err = inputproductservice.DbContext.SelectOne(inputproduct, "select * from "+inputproductservice.Table+
		" where order_id = ? and product_id = ?", order_id, product_id)
	if err != nil {
		log.Error("Error during getting input product object from database %v with value %v, %v", err, order_id, product_id)
		return nil, err
	}

	return inputproduct, nil
}

func (inputproductservice *InputProductService) GetByOrder(order_id int64) (inputproducts *[]models.ViewApiInputProduct, err error) {
	inputproducts = new([]models.ViewApiInputProduct)
	_, err = inputproductservice.DbContext.Select(inputproducts,
		"select product_id from "+inputproductservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all input product object from database %v with value %v", err, order_id)
		return nil, err
	}

	return inputproducts, nil
}

func (inputproductservice *InputProductService) GetAll(order_id int64) (inputproducts *[]models.DtoInputProduct, err error) {
	inputproducts = new([]models.DtoInputProduct)
	_, err = inputproductservice.DbContext.Select(inputproducts,
		"select * from "+inputproductservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all input product object from database %v with value %v", err, order_id)
		return nil, err
	}

	return inputproducts, nil
}

func (inputproductservice *InputProductService) Create(dtoinputproduct *models.DtoInputProduct, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoinputproduct)
	} else {
		err = inputproductservice.DbContext.Insert(dtoinputproduct)
	}
	if err != nil {
		log.Error("Error during creating input product object in database %v", err)
		return err
	}

	return nil
}

func (inputproductservice *InputProductService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+inputproductservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = inputproductservice.DbContext.Exec("delete from "+inputproductservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting input product objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
