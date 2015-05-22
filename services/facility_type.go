package services

import (
	"application/models"
)

type FacilityTypeRepository interface {
	Get(id int) (facilitytype *models.DtoFacilityType, err error)
	GetAll() (facilitytypes *[]models.ApiFacilityType, err error)
}

type FacilityTypeService struct {
	*Repository
}

func NewFacilityTypeService(repository *Repository) *FacilityTypeService {
	repository.DbContext.AddTableWithName(models.DtoFacilityType{}, repository.Table).SetKeys(false, "id")
	return &FacilityTypeService{Repository: repository}
}

func (facilitytypeservice *FacilityTypeService) Get(id int) (facilitytype *models.DtoFacilityType, err error) {
	facilitytype = new(models.DtoFacilityType)
	err = facilitytypeservice.DbContext.SelectOne(facilitytype, "select * from "+facilitytypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting facility type object from database %v with value %v", err, id)
		return nil, err
	}

	return facilitytype, nil
}

func (facilitytypeservice *FacilityTypeService) GetAll() (facilitytypes *[]models.ApiFacilityType, err error) {
	facilitytypes = new([]models.ApiFacilityType)
	_, err = facilitytypeservice.DbContext.Select(facilitytypes, "select id, name from "+facilitytypeservice.Table+" where active = 1")
	if err != nil {
		log.Error("Error during getting all facility type object from database %v", err)
		return nil, err
	}

	return facilitytypes, nil
}
