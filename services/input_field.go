package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type InputFieldRepository interface {
	Get(order_id int64, product_id int) (inputfield *models.DtoInputField, err error)
	GetByOrder(order_id int64) (inputfields *[]models.ViewApiInputField, err error)
	Create(dtoinputfield *models.DtoInputField, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type InputFieldService struct {
	*Repository
}

func NewInputFieldService(repository *Repository) *InputFieldService {
	repository.DbContext.AddTableWithName(models.DtoInputField{}, repository.Table).SetKeys(false, "order_id", "product_id")
	return &InputFieldService{Repository: repository}
}

func (inputfieldservice *InputFieldService) Get(order_id int64, product_id int) (inputfield *models.DtoInputField, err error) {
	inputfield = new(models.DtoInputField)
	err = inputfieldservice.DbContext.SelectOne(inputfield, "select * from "+inputfieldservice.Table+
		" where order_id = ? and product_id = ?", order_id, product_id)
	if err != nil {
		log.Error("Error during getting input field object from database %v with value %v, %v", err, order_id, product_id)
		return nil, err
	}

	return inputfield, nil
}

func (inputfieldservice *InputFieldService) GetByOrder(order_id int64) (inputfields *[]models.ViewApiInputField, err error) {
	inputfields = new([]models.ViewApiInputField)
	_, err = inputfieldservice.DbContext.Select(inputfields,
		"select product_id, count from "+inputfieldservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all input field object from database %v with value %v", err, order_id)
		return nil, err
	}

	return inputfields, nil
}

func (inputfieldservice *InputFieldService) Create(dtoinputfield *models.DtoInputField, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoinputfield)
	} else {
		err = inputfieldservice.DbContext.Insert(dtoinputfield)
	}
	if err != nil {
		log.Error("Error during creating input field object in database %v", err)
		return err
	}

	return nil
}

func (inputfieldservice *InputFieldService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+inputfieldservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = inputfieldservice.DbContext.Exec("delete from "+inputfieldservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting input field objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
