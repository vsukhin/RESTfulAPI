package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type AppendixRepository interface {
	Get(id int64) (appendix *models.DtoAppendix, err error)
	GetByContract(contract_id int64) (appendices *[]models.ApiAppendix, err error)
	Create(dtoappendix *models.DtoAppendix, trans *gorp.Transaction) (err error)
	DeleteByContract(contract_id int64, trans *gorp.Transaction) (err error)
}

type AppendixService struct {
	*Repository
}

func NewAppendixService(repository *Repository) *AppendixService {
	repository.DbContext.AddTableWithName(models.DtoAppendix{}, repository.Table).SetKeys(true, "id")
	return &AppendixService{Repository: repository}
}

func (appendixservice *AppendixService) Get(id int64) (appendix *models.DtoAppendix, err error) {
	appendix = new(models.DtoAppendix)
	err = appendixservice.DbContext.SelectOne(appendix, "select * from "+appendixservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting appendix object from database %v with value %v", err, id)
		return nil, err
	}

	return appendix, nil
}

func (appendixservice *AppendixService) GetByContract(contract_id int64) (appendices *[]models.ApiAppendix, err error) {
	appendices = new([]models.ApiAppendix)
	_, err = appendixservice.DbContext.Select(appendices,
		"select signed_date, name, file_id from "+appendixservice.Table+" where contract_id = ? and active = 1 order by signed_date asc", contract_id)
	if err != nil {
		log.Error("Error during getting all appendix object from database %v with value %v", err, contract_id)
		return nil, err
	}

	return appendices, nil
}

func (appendixservice *AppendixService) Create(dtoappendix *models.DtoAppendix, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtoappendix)
	} else {
		err = appendixservice.DbContext.Insert(dtoappendix)
	}
	if err != nil {
		log.Error("Error during creating appendix object in database %v", err)
		return err
	}

	return nil
}

func (appendixservice *AppendixService) DeleteByContract(contract_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+appendixservice.Table+" where contract_id = ?", contract_id)
	} else {
		_, err = appendixservice.DbContext.Exec("delete from "+appendixservice.Table+" where contract_id = ? and active = 1", contract_id)
	}
	if err != nil {
		log.Error("Error during deleting appendix objects for contract object in database %v with value %v", err, contract_id)
		return err
	}

	return nil
}
