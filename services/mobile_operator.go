package services

import (
	"application/models"
)

type MobileOperatorRepository interface {
	Get(id int) (mobileoperator *models.DtoMobileOperator, err error)
	GetDefault() (mobileoperator *models.DtoMobileOperator, err error)
	GetAll() (mobileoperators *[]models.ApiMobileOperator, err error)
	FindAll() (mobileoperators *[]models.DtoMobileOperator, err error)
}

type MobileOperatorService struct {
	*Repository
}

func NewMobileOperatorService(repository *Repository) *MobileOperatorService {
	repository.DbContext.AddTableWithName(models.DtoMobileOperator{}, repository.Table).SetKeys(true, "id")
	return &MobileOperatorService{Repository: repository}
}

func (mobileoperatorservice *MobileOperatorService) Get(id int) (mobileoperator *models.DtoMobileOperator, err error) {
	mobileoperator = new(models.DtoMobileOperator)
	err = mobileoperatorservice.DbContext.SelectOne(mobileoperator, "select * from "+mobileoperatorservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting mobile operator object from database %v with value %v", err, id)
		return nil, err
	}

	return mobileoperator, nil
}

func (mobileoperatorservice *MobileOperatorService) GetDefault() (mobileoperator *models.DtoMobileOperator, err error) {
	mobileoperator = new(models.DtoMobileOperator)
	err = mobileoperatorservice.DbContext.SelectOne(mobileoperator,
		"select * from "+mobileoperatorservice.Table+" where active = 1 and `default` = 1")
	if err != nil {
		log.Error("Error during getting default mobile operator object from database %v", err)
		return nil, err
	}

	return mobileoperator, nil
}

func (mobileoperatorservice *MobileOperatorService) GetAll() (mobileoperators *[]models.ApiMobileOperator, err error) {
	mobileoperators = new([]models.ApiMobileOperator)
	_, err = mobileoperatorservice.DbContext.Select(mobileoperators, "select id, shortname, longname, position from "+mobileoperatorservice.Table+
		" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all mobile operator object from database %v", err)
		return nil, err
	}

	return mobileoperators, nil
}

func (mobileoperatorservice *MobileOperatorService) FindAll() (mobileoperators *[]models.DtoMobileOperator, err error) {
	mobileoperators = new([]models.DtoMobileOperator)
	_, err = mobileoperatorservice.DbContext.Select(mobileoperators, "select * from "+mobileoperatorservice.Table+
		" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during finding all mobile operator object from database %v", err)
		return nil, err
	}

	return mobileoperators, nil
}
