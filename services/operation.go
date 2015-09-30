package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type OperationRepository interface {
	Get(id int64) (operation *models.DtoOperation, err error)
	GetByUnit(unit_id int64) (operations *[]models.DtoOperation, err error)
	CalculateBalance(unit_id int64) (money float64, err error)
	Create(dtooperation *models.DtoOperation, trans *gorp.Transaction) (err error)
}

type OperationService struct {
	*Repository
}

func NewOperationService(repository *Repository) *OperationService {
	repository.DbContext.AddTableWithName(models.DtoOperation{}, repository.Table).SetKeys(true, "id")
	return &OperationService{Repository: repository}
}

func (operationservice *OperationService) Get(id int64) (operation *models.DtoOperation, err error) {
	operation = new(models.DtoOperation)
	err = operationservice.DbContext.SelectOne(operation, "select * from "+operationservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting operation object from database %v with value %v", err, id)
		return nil, err
	}

	return operation, nil
}

func (operationservice *OperationService) GetByUnit(unit_id int64) (operations *[]models.DtoOperation, err error) {
	operations = new([]models.DtoOperation)
	_, err = operationservice.DbContext.Select(operations,
		"select * from "+operationservice.Table+" where unit_id = ?", unit_id)
	if err != nil {
		log.Error("Error during getting all operation objects from database %v with value %v", err, unit_id)
		return nil, err
	}

	return operations, nil
}

func (operationservice *OperationService) CalculateBalance(unit_id int64) (money float64, err error) {
	money, err = operationservice.DbContext.SelectFloat(
		"select coalesce(sum(d.money), 0) - coalesce(sum(c.money), 0) from debet d inner join credit c on d.unit_id = c.unit_id where d.unit_id = ?", unit_id)
	if err != nil {
		log.Error("Error during getting operation object from database %v with value %v", err, unit_id)
		return 0, err
	}

	return money, nil
}

func (operationservice *OperationService) Create(dtooperation *models.DtoOperation, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtooperation)
	} else {
		err = operationservice.DbContext.Insert(dtooperation)
	}
	if err != nil {
		log.Error("Error during creating operation object in database %v", err)
		return err
	}

	return nil
}
