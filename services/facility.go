package services

import (
	"application/models"
)

type FacilityRepository interface {
	Get(id int64) (facility *models.DtoFacility, err error)
	GetAll() (facilities *[]models.ApiFullFacility, err error)
	GetAllAvailable() (facilities *[]models.ApiShortFacility, err error)
	GetByUnit(unitid int64) (facilities *[]models.ApiLongFacility, err error)
	GetByUser(user_id int64) (facilities *[]models.ApiShortFacility, err error)
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

func (facilityservice *FacilityService) GetAll() (facilities *[]models.ApiFullFacility, err error) {
	facilities = new([]models.ApiFullFacility)
	_, err = facilityservice.DbContext.Select(facilities,
		"select id, category_id, alias, name, description, description_soon, active,"+
			" picNormal_id, picOver_id, picSoon_id, picDisable_id from "+facilityservice.Table)
	if err != nil {
		log.Error("Error during getting all facility object from database %v", err)
		return nil, err
	}

	return facilities, nil
}

func (facilityservice *FacilityService) GetAllAvailable() (facilities *[]models.ApiShortFacility, err error) {
	facilities = new([]models.ApiShortFacility)
	_, err = facilityservice.DbContext.Select(facilities, "select distinct f.id, f.name, f.description from "+facilityservice.Table+
		" f inner join supplier_services s on f.id = s.service_id where f.active = 1 and s.supplier_id in (select id from units where active = 1) and"+
		" f.id in (select service_id from price_properties p where published = 1 and customer_table_id in"+
		" (select id from customer_tables where active = 1 and permanent = 1 and type_id = ? and unit_id = s.supplier_id)"+
		" and ((date(end) != '0001-01-01' and now() <= end) or (date(end) = '0001-01-01'))"+
		" and ((date(begin) != '0001-01-01' and now() >= begin) or (date(begin) = '0001-01-01' and after_id = 0)"+
		" or (date(begin) = '0001-01-01' and after_id != 0 and"+
		" (select date(end) from price_properties where customer_table_id = p.after_id) != '0001-01-01' and"+
		" now() > (select end from price_properties where customer_table_id = p.after_id))))", models.TABLE_TYPE_PRICE)
	if err != nil {
		log.Error("Error during getting all facility object from database %v", err)
		return nil, err
	}

	return facilities, nil
}

func (facilityservice *FacilityService) GetByUnit(unitid int64) (facilities *[]models.ApiLongFacility, err error) {
	facilities = new([]models.ApiLongFacility)
	_, err = facilityservice.DbContext.Select(facilities, "select id, name, description, active from "+facilityservice.Table+
		" where id in (select service_id from supplier_services where supplier_id = ?)", unitid)
	if err != nil {
		log.Error("Error during getting all facility object from database %v with value %v", err, unitid)
		return nil, err
	}

	return facilities, nil
}

func (facilityservice *FacilityService) GetByUser(user_id int64) (facilities *[]models.ApiShortFacility, err error) {
	facilities = new([]models.ApiShortFacility)
	_, err = facilityservice.DbContext.Select(facilities, "select id, name, description from "+facilityservice.Table+
		" where active = 1 and id in (select service_id from supplier_services where supplier_id = "+
		"(select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting all facility object from database %v with value %v", err, user_id)
		return nil, err
	}

	return facilities, nil
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
