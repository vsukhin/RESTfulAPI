package services

import (
	"application/models"
)

type UnitRepository interface {
	FindByUser(userid int64) (unit *models.DtoUnit, err error)
	Get(unitid int64) (unit *models.DtoUnit, err error)
	Create(unit *models.DtoUnit) (err error)
	Update(unit *models.DtoUnit) (err error)
	Delete(unitid int64) (err error)
}

type UnitService struct {
	*Repository
}

func NewUnitService(repository *Repository) *UnitService {
	repository.DbContext.AddTableWithName(models.DtoUnit{}, repository.Table).SetKeys(true, "id")
	return &UnitService{
		repository,
	}
}

func (unitservice *UnitService) FindByUser(userid int64) (unit *models.DtoUnit, err error) {
	unit = new(models.DtoUnit)
	err = unitservice.DbContext.SelectOne(unit,
		"select un.* from "+unitservice.Table+" un inner join users us on un.id = us.unit_id where us.id = ?", userid)

	if err != nil {
		log.Error("Error during getting unit object from database %v with value %v", err, userid)
		return nil, err
	}

	return unit, nil
}

func (unitservice *UnitService) Get(unitid int64) (unit *models.DtoUnit, err error) {
	unit = new(models.DtoUnit)
	err = unitservice.DbContext.SelectOne(unit, "select * from "+unitservice.Table+" where id = ?", unitid)
	if err != nil {
		log.Error("Error during getting unit object from database %v with value %v", err, unitid)
		return nil, err
	}

	return unit, nil
}

func (unitservice *UnitService) Create(unit *models.DtoUnit) (err error) {
	err = unitservice.DbContext.Insert(unit)
	if err != nil {
		log.Error("Error during creating unit object in database %v", err)
		return err
	}

	return nil
}

func (unitservice *UnitService) Update(unit *models.DtoUnit) (err error) {
	_, err = unitservice.DbContext.Update(unit)
	if err != nil {
		log.Error("Error during updating unit object in database %v with value %v", err, unit.ID)
		return err
	}

	return nil
}

func (unitservice *UnitService) Delete(unitid int64) (err error) {
	_, err = unitservice.DbContext.Exec("delete from "+unitservice.Table+" where id = ?", unitid)
	if err != nil {
		log.Error("Error during deleting unit object in database %v with value %v", err, unitid)
		return err
	}

	return nil
}
