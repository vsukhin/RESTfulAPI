package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

type FacilityRepository interface {
	Get(id int64) (facility *models.DtoFacility, err error)
	GetAll() (facilities *[]models.DtoFacility, err error)
	GetByUser(user_id int64) (facilities *[]models.ApiShortFacility, err error)
	SetByUser(user_id int64, facilities *[]int64, inTrans bool) (err error)
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

func (facilityservice *FacilityService) Get(id int64) (facility *models.DtoFacility, err error) {
	facility = new(models.DtoFacility)
	err = facilityservice.DbContext.SelectOne(facility, "select * from "+facilityservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting facility object from database %v with value %v", err, id)
		return nil, err
	}

	return facility, nil
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

func (facilityservice *FacilityService) GetByUser(user_id int64) (facilities *[]models.ApiShortFacility, err error) {
	facilities = new([]models.ApiShortFacility)
	_, err = facilityservice.DbContext.Select(facilities, "select id, name, description from "+facilityservice.Table+
		" where id in (select service_id from supplier_services where supplier_id = "+
		"(select unit_id from users where id = ? and active = 1 and confirmed = 1)) and active = 1", user_id)
	if err != nil {
		log.Error("Error during getting all facility object from database %v with value %v", err, user_id)
		return nil, err
	}

	return facilities, nil
}

func (facilityservice *FacilityService) SetByUser(user_id int64, facilities *[]int64, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = facilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during setting facility objects for user object in database %v", err)
			return err
		}
	}

	_, err = facilityservice.DbContext.Exec("delete from supplier_services where supplier_id = "+
		"(select unit_id from users where id = ? and active = 1 and confirmed = 1)", user_id)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during setting facility objects for user object in database %v with value %v", err, user_id)
		return err
	}

	if len(*facilities) > 0 {
		statement := ""
		for _, value := range *facilities {
			if statement != "" {
				statement += " union"
			}
			statement += fmt.Sprintf(" select (select unit_id from users where id = %v and active = 1 and confirmed = 1), %v", user_id, value)
		}
		_, err = facilityservice.DbContext.Exec("insert into supplier_services (supplier_id, service_id)" + statement)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during setting facility objects for user object in database %v with value %v", err, user_id)
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during setting facility objects for user object in database %v", err)
			return err
		}
	}

	return nil
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
