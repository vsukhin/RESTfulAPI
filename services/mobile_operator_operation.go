package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type MobileOperatorOperationRepository interface {
	Get(order_id int64, mobileoperator_id int) (mobileoperatoroperation *models.DtoMobileOperatorOperation, err error)
	GetByOrder(order_id int64) (mobileoperatoroperations *[]models.ViewApiMobileOperatorOperation, err error)
	Create(dtomobileoperatoroperation *models.DtoMobileOperatorOperation, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type MobileOperatorOperationService struct {
	*Repository
}

func NewMobileOperatorOperationService(repository *Repository) *MobileOperatorOperationService {
	repository.DbContext.AddTableWithName(models.DtoMobileOperatorOperation{}, repository.Table).SetKeys(false, "order_id", "mobileoperator_id")
	return &MobileOperatorOperationService{Repository: repository}
}

func (mobileoperatoroperationservice *MobileOperatorOperationService) Get(order_id int64,
	mobileoperator_id int) (mobileoperatoroperation *models.DtoMobileOperatorOperation, err error) {
	mobileoperatoroperation = new(models.DtoMobileOperatorOperation)
	err = mobileoperatoroperationservice.DbContext.SelectOne(mobileoperatoroperation, "select * from "+mobileoperatoroperationservice.Table+
		" where order_id = ? and mobileoperator_id = ?", order_id, mobileoperator_id)
	if err != nil {
		log.Error("Error during getting mobile operator operation object from database %v with value %v, %v", err, order_id, mobileoperator_id)
		return nil, err
	}

	return mobileoperatoroperation, nil
}

func (mobileoperatoroperationservice *MobileOperatorOperationService) GetByOrder(
	order_id int64) (mobileoperatoroperations *[]models.ViewApiMobileOperatorOperation, err error) {
	mobileoperatoroperations = new([]models.ViewApiMobileOperatorOperation)
	_, err = mobileoperatoroperationservice.DbContext.Select(mobileoperatoroperations,
		"select mobileoperator_id, percent, count from "+mobileoperatoroperationservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all mobile operator operation object from database %v with value %v", err, order_id)
		return nil, err
	}

	return mobileoperatoroperations, nil
}

func (mobileoperatoroperationservice *MobileOperatorOperationService) Create(dtomobileoperatoroperation *models.DtoMobileOperatorOperation,
	trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtomobileoperatoroperation)
	} else {
		err = mobileoperatoroperationservice.DbContext.Insert(dtomobileoperatoroperation)
	}
	if err != nil {
		log.Error("Error during creating mobile operator operation object in database %v", err)
		return err
	}

	return nil
}

func (mobileoperatoroperationservice *MobileOperatorOperationService) DeleteByOrder(order_id int64,
	trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+mobileoperatoroperationservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = mobileoperatoroperationservice.DbContext.Exec("delete from "+mobileoperatoroperationservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting mobile operator operation objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
