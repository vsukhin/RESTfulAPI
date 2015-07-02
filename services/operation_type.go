package services

import (
	"application/models"
)

type OperationTypeRepository interface {
	Get(id int) (operationtype *models.DtoOperationType, err error)
}

type OperationTypeService struct {
	*Repository
}

func NewOperationTypeService(repository *Repository) *OperationTypeService {
	repository.DbContext.AddTableWithName(models.DtoOperationType{}, repository.Table).SetKeys(false, "id")
	return &OperationTypeService{
		repository,
	}
}

func (operationtypeservice *OperationTypeService) Get(id int) (operationtype *models.DtoOperationType, err error) {
	operationtype = new(models.DtoOperationType)
	err = operationtypeservice.DbContext.SelectOne(operationtype, "select * from "+operationtypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting operation type object from database %v with value %v", err, id)
		return nil, err
	}

	return operationtype, nil
}
