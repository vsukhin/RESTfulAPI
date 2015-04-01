package services

import (
	"application/models"
)

type StatusRepository interface {
	Get(id int) (status *models.DtoStatus, err error)
	GetAll() (statuses *[]models.DtoStatus, err error)
	FindByName(name string) (id int64, err error)
}

type StatusService struct {
	*Repository
}

func NewStatusService(repository *Repository) *StatusService {
	repository.DbContext.AddTableWithName(models.DtoStatus{}, repository.Table).SetKeys(false, "id")
	return &StatusService{
		repository,
	}
}

func (statusservice *StatusService) Get(id int) (status *models.DtoStatus, err error) {
	status = new(models.DtoStatus)
	err = statusservice.DbContext.SelectOne(status, "select * from "+statusservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting status object from database %v with value %v", err, id)
		return nil, err
	}

	return status, nil
}

func (statusservice *StatusService) GetAll() (statuses *[]models.DtoStatus, err error) {
	statuses = new([]models.DtoStatus)
	_, err = statusservice.DbContext.Select(statuses, "select * from "+statusservice.Table)
	if err != nil {
		log.Error("Error during getting all status object from database %v", err)
		return nil, err
	}

	return statuses, nil
}

func (statusservice *StatusService) FindByName(name string) (id int64, err error) {
	err = statusservice.DbContext.SelectOne(&id, "select id from "+statusservice.Table+" where name = ?", name)
	if err != nil {
		log.Error("Error during finding data format object from database %v with value %v", err, name)
		return 0, err
	}

	return id, nil
}
