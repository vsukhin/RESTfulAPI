package services

import (
	"application/models"
)

type HeaderFacilityRepository interface {
	Exists(orderid int64) (found bool, err error)
	Get(order_id int64) (headerfacility *models.DtoHeaderFacility, err error)
	Create(headerfacility *models.DtoHeaderFacility) (err error)
	Update(headerfacility *models.DtoHeaderFacility) (err error)
	Save(headerfacility *models.DtoHeaderFacility) (err error)
}

type HeaderFacilityService struct {
	*Repository
}

func NewHeaderFacilityService(repository *Repository) *HeaderFacilityService {
	repository.DbContext.AddTableWithName(models.DtoHeaderFacility{}, repository.Table).SetKeys(false, "order_id")
	return &HeaderFacilityService{Repository: repository}
}

func (headerfacilityservice *HeaderFacilityService) Exists(orderid int64) (found bool, err error) {
	var count int64
	count, err = headerfacilityservice.DbContext.SelectInt("select count(*) from "+headerfacilityservice.Table+
		" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during checking header facility object in database %v with value %v", err, orderid)
		return false, err
	}

	return count != 0, nil
}

func (headerfacilityservice *HeaderFacilityService) Get(orderid int64) (headerfacility *models.DtoHeaderFacility, err error) {
	headerfacility = new(models.DtoHeaderFacility)
	err = headerfacilityservice.DbContext.SelectOne(headerfacility, "select * from "+headerfacilityservice.Table+" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during getting header facility object from database %v with value %v", err, orderid)
		return nil, err
	}

	return headerfacility, nil
}

func (headerfacilityservice *HeaderFacilityService) Create(headerfacility *models.DtoHeaderFacility) (err error) {
	err = headerfacilityservice.DbContext.Insert(headerfacility)
	if err != nil {
		log.Error("Error during creating header facility object in database %v", err)
		return err
	}

	return nil
}

func (headerfacilityservice *HeaderFacilityService) Update(headerfacility *models.DtoHeaderFacility) (err error) {
	_, err = headerfacilityservice.DbContext.Update(headerfacility)
	if err != nil {
		log.Error("Error during updating header facility object in database %v with value %v", err, headerfacility.Order_ID)
		return err
	}

	return nil
}

func (headerfacilityservice *HeaderFacilityService) Save(headerfacility *models.DtoHeaderFacility) (err error) {
	count, err := headerfacilityservice.DbContext.SelectInt("select count(*) from "+headerfacilityservice.Table+
		" where order_id = ?", headerfacility.Order_ID)
	if err != nil {
		log.Error("Error during saving header facility object in database %v with value %v", err, headerfacility.Order_ID)
		return err
	}
	if count == 0 {
		err = headerfacilityservice.Create(headerfacility)
	} else {
		err = headerfacilityservice.Update(headerfacility)
	}

	return err
}
