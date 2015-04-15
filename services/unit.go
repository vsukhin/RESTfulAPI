package services

import (
	"application/models"
)

type UnitRepository interface {
	FindByUser(userid int64) (unit *models.DtoUnit, err error)
	Get(unitid int64) (unit *models.DtoUnit, err error)
	GetMeta() (unit *models.ApiShortMetaUnit, err error)
	GetAll(filter string) (units *[]models.ApiShortUnit, err error)
	Create(unit *models.DtoUnit) (err error)
	Update(unit *models.DtoUnit) (err error)
	Deactivate(*models.DtoUnit) (err error)
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

func (unitservice *UnitService) GetMeta() (unit *models.ApiShortMetaUnit, err error) {
	unit = new(models.ApiShortMetaUnit)
	unit.Total, err = unitservice.DbContext.SelectInt("select count(*) from " + unitservice.Table)
	if err != nil {
		log.Error("Error during getting unit meta object from database %v", err)
		return nil, err
	}

	return unit, nil
}

func (unitservice *UnitService) GetAll(filter string) (units *[]models.ApiShortUnit, err error) {
	units = new([]models.ApiShortUnit)
	_, err = unitservice.DbContext.Select(units, "select id, name from "+unitservice.Table+filter)
	if err != nil {
		log.Error("Error during getting all unit object from database %v", err)
		return nil, err
	}

	return units, nil
}

func (unitservice *UnitService) Create(unit *models.DtoUnit) (err error) {
	if unit.Name == "" {
		unit.Name = models.UNIT_NAME_DEFAULT
	}

	err = unitservice.DbContext.Insert(unit)
	if err != nil {
		log.Error("Error during creating unit object in database %v", err)
		return err
	}

	return nil
}

func (unitservice *UnitService) Update(unit *models.DtoUnit) (err error) {
	if unit.Name == "" {
		unit.Name = models.UNIT_NAME_DEFAULT
	}

	_, err = unitservice.DbContext.Update(unit)
	if err != nil {
		log.Error("Error during updating unit object in database %v with value %v", err, unit.ID)
		return err
	}

	return nil
}

func (unitservice *UnitService) Deactivate(unit *models.DtoUnit) (err error) {
	_, err = unitservice.DbContext.Exec("update "+unitservice.Table+" set active = 0 where id = ?", unit.ID)
	if err != nil {
		log.Error("Error during deactivating unit object in database %v with value %v", err, unit.ID)
		return err
	}

	return nil
}
