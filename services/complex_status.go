package services

import (
	"application/models"
)

type ComplexStatusRepository interface {
	Get(id int) (complexstatus *models.DtoComplexStatus, err error)
	GetAll() (complexstatuses *[]models.ApiComplexStatus, err error)
}

type ComplexStatusService struct {
	*Repository
}

func NewComplexStatusService(repository *Repository) *ComplexStatusService {
	repository.DbContext.AddTableWithName(models.DtoComplexStatus{}, repository.Table).SetKeys(false, "id")
	return &ComplexStatusService{Repository: repository}
}

func (complexstatusservice *ComplexStatusService) Get(id int) (complexstatus *models.DtoComplexStatus, err error) {
	complexstatus = new(models.DtoComplexStatus)
	err = complexstatusservice.DbContext.SelectOne(complexstatus, "select * from "+complexstatusservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting complex status object from database %v with value %v", err, id)
		return nil, err
	}

	return complexstatus, nil
}

func (complexstatusservice *ComplexStatusService) GetAll() (complexstatuses *[]models.ApiComplexStatus, err error) {
	complexstatuses = new([]models.ApiComplexStatus)
	_, err = complexstatusservice.DbContext.Select(complexstatuses, "select id, final, name, description from "+complexstatusservice.Table+
		" where active = 1")
	if err != nil {
		log.Error("Error during getting all complex status object from database %v", err)
		return nil, err
	}

	return complexstatuses, nil
}
