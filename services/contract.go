package services

import (
	"application/models"
	"fmt"
	"github.com/coopernurse/gorp"
)

const (
	CONTRACT_NAME_TEMPLATE = "Номер контракта %v"
)

type ContractRepository interface {
	Get(id int64) (contract *models.DtoContract, err error)
	GetByUser(userid int64) (contracts *[]models.DtoContract, err error)
	GetByUnit(unitid int64) (contracts *[]models.DtoContract, err error)
	SetArrays(contract *models.DtoContract, trans *gorp.Transaction) (err error)
	Create(contract *models.DtoContract, briefly bool, trans *gorp.Transaction) (err error)
	Update(contract *models.DtoContract, briefly bool, trans *gorp.Transaction) (err error)
	SaveAll(contracts *[]models.DtoContract, briefly bool) (err error)
	Deactivate(contract *models.DtoContract) (err error)
}

type ContractService struct {
	AppendixRepository AppendixRepository
	*Repository
}

func NewContractService(repository *Repository) *ContractService {
	repository.DbContext.AddTableWithName(models.DtoContract{}, repository.Table).SetKeys(true, "id")
	return &ContractService{Repository: repository}
}

func (contractservice *ContractService) Get(id int64) (contract *models.DtoContract, err error) {
	contract = new(models.DtoContract)
	err = contractservice.DbContext.SelectOne(contract, "select * from "+contractservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting contract object from database %v with value %v", err, id)
		return nil, err
	}

	return contract, nil
}

func (contractservice *ContractService) GetByUser(userid int64) (contracts *[]models.DtoContract, err error) {
	contracts = new([]models.DtoContract)
	_, err = contractservice.DbContext.Select(contracts, "select * from "+contractservice.Table+
		" where company_id in (select id from companies where unit_id = (select unit_id from users where id = ?) and active = 1) and active = 1", userid)
	if err != nil {
		log.Error("Error during getting unit contract object from database %v with value %v", err, userid)
		return nil, err
	}

	return contracts, nil
}

func (contractservice *ContractService) GetByUnit(unitid int64) (contracts *[]models.DtoContract, err error) {
	contracts = new([]models.DtoContract)
	_, err = contractservice.DbContext.Select(contracts, "select * from "+contractservice.Table+
		" where company_id in (select id from companies where unit_id = ? and active = 1) and active = 1", unitid)
	if err != nil {
		log.Error("Error during getting unit contract object from database %v with value %v", err, unitid)
		return nil, err
	}

	return contracts, nil
}

func (contractservice *ContractService) SetArrays(contract *models.DtoContract, trans *gorp.Transaction) (err error) {
	err = contractservice.AppendixRepository.DeleteByContract(contract.ID, trans)
	if err != nil {
		log.Error("Error during setting contract object in database %v with value %v", err, contract.ID)
		return err
	}
	for _, dtoappendix := range contract.Appendices {
		dtoappendix.Contract_ID = contract.ID
		err = contractservice.AppendixRepository.Create(&dtoappendix, trans)
		if err != nil {
			log.Error("Error during setting contract object in database %v with value %v", err, contract.ID)
			return err
		}
	}

	return nil
}

func (contractservice *ContractService) Create(contract *models.DtoContract, briefly bool, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(contract)
	} else {
		err = contractservice.DbContext.Insert(contract)
	}
	if err != nil {
		log.Error("Error during creating contract object in database %v", err)
		return err
	}

	if !briefly {
		err = contractservice.SetArrays(contract, trans)
		if err != nil {
			return err
		}
	}

	return nil
}

func (contractservice *ContractService) Update(contract *models.DtoContract, briefly bool, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Update(contract)
	} else {
		_, err = contractservice.DbContext.Update(contract)
	}
	if err != nil {
		log.Error("Error during updating contract object in database %v with value %v", err, contract.ID)
		return err
	}

	if !briefly {
		err = contractservice.SetArrays(contract, trans)
		if err != nil {
			return err
		}
	}

	return nil
}

func (contractservice *ContractService) SaveAll(contracts *[]models.DtoContract, briefly bool) (err error) {
	var trans *gorp.Transaction

	trans, err = contractservice.DbContext.Begin()
	if err != nil {
		log.Error("Error during saving contract object in database %v", err)
		return err
	}

	for i := range *contracts {
		if (*contracts)[i].ID == 0 {
			err = contractservice.Create(&(*contracts)[i], briefly, trans)
			if err != nil {
				_ = trans.Rollback()
				return err
			}
			(*contracts)[i].Name = fmt.Sprintf(CONTRACT_NAME_TEMPLATE, (*contracts)[i].ID)
		}
		err = contractservice.Update(&(*contracts)[i], briefly, trans)
		if err != nil {
			_ = trans.Rollback()
			return err
		}
	}

	err = trans.Commit()
	if err != nil {
		log.Error("Error during saving contract object in database %v", err)
		return err
	}

	return nil
}

func (contractservice *ContractService) Deactivate(contract *models.DtoContract) (err error) {
	_, err = contractservice.DbContext.Exec("update "+contractservice.Table+" set active = 0 where id = ?", contract.ID)
	if err != nil {
		log.Error("Error during deactivating contract object in database %v with value %v", err, contract.ID)
		return err
	}

	return nil
}
