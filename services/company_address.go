package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type CompanyAddressRepository interface {
	Get(id int64) (companyaddress *models.DtoCompanyAddress, err error)
	GetByCompany(company_id int64) (companyaddresss *[]models.ViewApiCompanyAddress, err error)
	Create(dtocompanyaddress *models.DtoCompanyAddress, trans *gorp.Transaction) (err error)
	DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error)
}

type CompanyAddressService struct {
	*Repository
}

func NewCompanyAddressService(repository *Repository) *CompanyAddressService {
	repository.DbContext.AddTableWithName(models.DtoCompanyAddress{}, repository.Table).SetKeys(true, "id")
	return &CompanyAddressService{Repository: repository}
}

func (companyaddressservice *CompanyAddressService) Get(id int64) (companyaddress *models.DtoCompanyAddress, err error) {
	companyaddress = new(models.DtoCompanyAddress)
	err = companyaddressservice.DbContext.SelectOne(companyaddress, "select * from "+companyaddressservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company address object from database %v with value %v", err, id)
		return nil, err
	}

	return companyaddress, nil
}

func (companyaddressservice *CompanyAddressService) GetByCompany(company_id int64) (companyaddresss *[]models.ViewApiCompanyAddress, err error) {
	companyaddresss = new([]models.ViewApiCompanyAddress)
	_, err = companyaddressservice.DbContext.Select(companyaddresss,
		"select c.`primary`, c.ditto, c.address_type_id, c.zip, c.country, c.region, c.city, c.street, c.building, c.postbox, c.company, c.comments from "+
			companyaddressservice.Table+" c inner join address_types a on c.address_type_id = a.id where company_id = ? and a.active = 1 order by a.position asc",
		company_id)
	if err != nil {
		log.Error("Error during getting all company address object from database %v with value %v", err, company_id)
		return nil, err
	}

	return companyaddresss, nil
}

func (companyaddressservice *CompanyAddressService) Create(dtocompanyaddress *models.DtoCompanyAddress, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtocompanyaddress)
	} else {
		err = companyaddressservice.DbContext.Insert(dtocompanyaddress)
	}
	if err != nil {
		log.Error("Error during creating company address object in database %v", err)
		return err
	}

	return nil
}

func (companyaddressservice *CompanyAddressService) DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+companyaddressservice.Table+" where company_id = ?", company_id)
	} else {
		_, err = companyaddressservice.DbContext.Exec("delete from "+companyaddressservice.Table+" where company_id = ?", company_id)
	}
	if err != nil {
		log.Error("Error during deleting company address objects for company object in database %v with value %v", err, company_id)
		return err
	}

	return nil
}
