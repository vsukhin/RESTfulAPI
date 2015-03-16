package services

import (
	"application/models"
)

type FacilityRepository interface {
	GetAll() (facilities *[]models.DtoFacility, err error)
	Get(id int64) (facility *models.DtoFacility, err error)
	Create(facility *models.DtoFacility) (err error)
	Update(facility *models.DtoFacility) (err error)
	Delete(id int64) (err error)
}

type FacilityService struct {
	*Repository
}

func NewFacilityService(repository *Repository) *FacilityService {
	repository.DbContext.AddTableWithName(models.DtoFacility{}, repository.Table).SetKeys(true, "id")
	return &FacilityService{
		repository,
	}
}

func (facilityservice *FacilityService) GetAll() (facilities *[]models.DtoFacility, err error) {
	facilities = new([]models.DtoFacility)
	_, err = facilityservice.DbContext.Select(facilities, "select * from "+facilityservice.Table)
	if err != nil {
		log.Error("Error during getting all facility object from database %v", err)
		return nil, err
	}

	return facilities, nil
}

func (facilityservice *FacilityService) Get(id int64) (facility *models.DtoFacility, err error) {
	facility = new(models.DtoFacility)
	err = facilityservice.DbContext.SelectOne(facility, "select * from "+facilityservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting facility object from database %v with value %v", err, id)
		return nil, err
	}

	return facility, nil
}

func (facilityservice *FacilityService) Create(facility *models.DtoFacility) (err error) {
	err = facilityservice.DbContext.Insert(facility)
	if err != nil {
		log.Error("Error during creating facility object in database %v", err)
		return err
	}

	return nil
}

func (facilityservice *FacilityService) Update(facility *models.DtoFacility) (err error) {
	_, err = facilityservice.DbContext.Update(facility)
	if err != nil {
		log.Error("Error during updating facility object in database %v with value %v", err, facility.ID)
		return err
	}

	return nil
}

func (facilityservice *FacilityService) Delete(id int64) (err error) {
	_, err = facilityservice.DbContext.Exec("delete from "+facilityservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during deleting facility object in database %v with value %v", err, id)
		return err
	}

	return nil
}
