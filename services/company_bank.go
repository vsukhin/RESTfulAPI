package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type CompanyBankRepository interface {
	Get(id int64) (companybank *models.DtoCompanyBank, err error)
	GetByCompany(company_id int64) (companybanks *[]models.ViewApiCompanyBank, err error)
	Create(dtocompanybank *models.DtoCompanyBank, trans *gorp.Transaction) (err error)
	DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error)
}

type CompanyBankService struct {
	*Repository
}

func NewCompanyBankService(repository *Repository) *CompanyBankService {
	repository.DbContext.AddTableWithName(models.DtoCompanyBank{}, repository.Table).SetKeys(true, "id")
	return &CompanyBankService{Repository: repository}
}

func (companybankservice *CompanyBankService) Get(id int64) (companybank *models.DtoCompanyBank, err error) {
	companybank = new(models.DtoCompanyBank)
	err = companybankservice.DbContext.SelectOne(companybank, "select * from "+companybankservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company bank object from database %v with value %v", err, id)
		return nil, err
	}

	return companybank, nil
}

func (companybankservice *CompanyBankService) GetByCompany(company_id int64) (companybanks *[]models.ViewApiCompanyBank, err error) {
	companybanks = new([]models.ViewApiCompanyBank)
	_, err = companybankservice.DbContext.Select(companybanks,
		"select `primary`, bik, name, checking_account, corresponding_account, not active as del from "+
			companybankservice.Table+" where company_id = ?", company_id)
	if err != nil {
		log.Error("Error during getting all company bank object from database %v with value %v", err, company_id)
		return nil, err
	}

	return companybanks, nil
}

func (companybankservice *CompanyBankService) Create(dtocompanybank *models.DtoCompanyBank, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtocompanybank)
	} else {
		err = companybankservice.DbContext.Insert(dtocompanybank)
	}
	if err != nil {
		log.Error("Error during creating company bank object in database %v", err)
		return err
	}

	return nil
}

func (companybankservice *CompanyBankService) DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+companybankservice.Table+" where company_id = ?", company_id)
	} else {
		_, err = companybankservice.DbContext.Exec("delete from "+companybankservice.Table+" where company_id = ?", company_id)
	}
	if err != nil {
		log.Error("Error during deleting company bank objects for company object in database %v with value %v", err, company_id)
		return err
	}

	return nil
}
