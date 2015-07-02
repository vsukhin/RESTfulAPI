package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type TransactionRepository interface {
	Get(id int64) (transaction *models.DtoTransaction, err error)
	GetByUnit(unit_id int64) (transactions *[]models.DtoTransaction, err error)
	Create(dtotransaction *models.DtoTransaction, trans *gorp.Transaction) (err error)
}

type TransactionService struct {
	*Repository
}

func NewTransactionService(repository *Repository) *TransactionService {
	repository.DbContext.AddTableWithName(models.DtoTransaction{}, repository.Table).SetKeys(true, "id")
	return &TransactionService{Repository: repository}
}

func (transactionservice *TransactionService) Get(id int64) (transaction *models.DtoTransaction, err error) {
	transaction = new(models.DtoTransaction)
	err = transactionservice.DbContext.SelectOne(transaction, "select * from "+transactionservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting transaction object from database %v with value %v", err, id)
		return nil, err
	}

	return transaction, nil
}

func (transactionservice *TransactionService) GetByUnit(unit_id int64) (transactions *[]models.DtoTransaction, err error) {
	transactions = new([]models.DtoTransaction)
	_, err = transactionservice.DbContext.Select(transactions,
		"select * from "+transactionservice.Table+" where source_id = ? or destination_id = ?", unit_id, unit_id)
	if err != nil {
		log.Error("Error during getting all transaction objects from database %v with value %v", err, unit_id)
		return nil, err
	}

	return transactions, nil
}

func (transactionservice *TransactionService) Create(dtotransaction *models.DtoTransaction, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtotransaction)
	} else {
		err = transactionservice.DbContext.Insert(dtotransaction)
	}
	if err != nil {
		log.Error("Error during creating transaction object in database %v", err)
		return err
	}

	return nil
}
