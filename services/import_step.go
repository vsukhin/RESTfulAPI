package services

import (
	"application/models"
)

type ImportStepRepository interface {
	Get(tableid int64, step byte) (importstep *models.DtoImportStep, err error)
	GetByTable(tableid int64) (importsteps *[]models.ApiImportStep, err error)
	Create(importstep *models.DtoImportStep) (err error)
	Update(importstep *models.DtoImportStep) (err error)
}

type ImportStepService struct {
	*Repository
}

func NewImportStepService(repository *Repository) *ImportStepService {
	repository.DbContext.AddTableWithName(models.DtoImportStep{}, repository.Table).SetKeys(false, "customer_table_id", "step")
	return &ImportStepService{
		repository,
	}
}

func (importstepservice *ImportStepService) Get(tableid int64, step byte) (importstep *models.DtoImportStep, err error) {
	importstep = new(models.DtoImportStep)
	err = importstepservice.DbContext.SelectOne(importstep, "select * from "+importstepservice.Table+" where customer_table_id = ? and step = ?", tableid, step)
	if err != nil {
		log.Error("Error during getting import step object from database %v with value %v, %v", err, tableid, step)
		return nil, err
	}

	return importstep, nil
}

func (importstepservice *ImportStepService) GetByTable(tableid int64) (importsteps *[]models.ApiImportStep, err error) {
	importsteps = new([]models.ApiImportStep)
	_, err = importstepservice.DbContext.Select(importsteps, "select step, ready, percentage from "+importstepservice.Table+" where customer_table_id = ?", tableid)
	if err != nil {
		log.Error("Error during getting all import step object from database %v with value %v", err, tableid)
		return nil, err
	}

	return importsteps, nil
}

func (importstepservice *ImportStepService) Create(importstep *models.DtoImportStep) (err error) {
	err = importstepservice.DbContext.Insert(importstep)
	if err != nil {
		log.Error("Error during creating import step object in database %v", err)
		return err
	}

	return nil
}

func (importstepservice *ImportStepService) Update(importstep *models.DtoImportStep) (err error) {
	_, err = importstepservice.DbContext.Update(importstep)
	if err != nil {
		log.Error("Error during updating import step object in database %v with value %v, %v", err, importstep.Customer_Table_ID, importstep.Step)
		return err
	}

	return nil
}
