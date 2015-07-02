package services

import (
	"application/models"
)

type TransactionTypeRepository interface {
	Get(id int) (transactiontype *models.DtoTransactionType, err error)
}

type TransactionTypeService struct {
	*Repository
}

func NewTransactionTypeService(repository *Repository) *TransactionTypeService {
	repository.DbContext.AddTableWithName(models.DtoTransactionType{}, repository.Table).SetKeys(false, "id")
	return &TransactionTypeService{
		repository,
	}
}

func (transactiontypeservice *TransactionTypeService) Get(id int) (transactiontype *models.DtoTransactionType, err error) {
	transactiontype = new(models.DtoTransactionType)
	err = transactiontypeservice.DbContext.SelectOne(transactiontype, "select * from "+transactiontypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting transaction type object from database %v with value %v", err, id)
		return nil, err
	}

	return transactiontype, nil
}
