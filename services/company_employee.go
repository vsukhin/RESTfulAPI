package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type CompanyEmployeeRepository interface {
	Get(id int64) (companyemployee *models.DtoCompanyEmployee, err error)
	GetByCompany(company_id int64) (companyemployees *[]models.ViewApiCompanyEmployee, err error)
	Create(dtocompanyemployee *models.DtoCompanyEmployee, trans *gorp.Transaction) (err error)
	DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error)
}

type CompanyEmployeeService struct {
	*Repository
}

func NewCompanyEmployeeService(repository *Repository) *CompanyEmployeeService {
	repository.DbContext.AddTableWithName(models.DtoCompanyEmployee{}, repository.Table).SetKeys(true, "id")
	return &CompanyEmployeeService{Repository: repository}
}

func (companyemployeeservice *CompanyEmployeeService) Get(id int64) (companyemployee *models.DtoCompanyEmployee, err error) {
	companyemployee = new(models.DtoCompanyEmployee)
	err = companyemployeeservice.DbContext.SelectOne(companyemployee, "select * from "+companyemployeeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company employee object from database %v with value %v", err, id)
		return nil, err
	}

	return companyemployee, nil
}

func (companyemployeeservice *CompanyEmployeeService) GetByCompany(company_id int64) (companyemployees *[]models.ViewApiCompanyEmployee, err error) {
	companyemployees = new([]models.ViewApiCompanyEmployee)
	_, err = companyemployeeservice.DbContext.Select(companyemployees,
		"select employee_type, ditto, surname, name, middlename, base, not active as del from "+
			companyemployeeservice.Table+" where company_id = ?", company_id)
	if err != nil {
		log.Error("Error during getting all company employee object from database %v with value %v", err, company_id)
		return nil, err
	}

	return companyemployees, nil
}

func (companyemployeeservice *CompanyEmployeeService) Create(dtocompanyemployee *models.DtoCompanyEmployee, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtocompanyemployee)
	} else {
		err = companyemployeeservice.DbContext.Insert(dtocompanyemployee)
	}
	if err != nil {
		log.Error("Error during creating company employee object in database %v", err)
		return err
	}

	return nil
}

func (companyemployeeservice *CompanyEmployeeService) DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+companyemployeeservice.Table+" where company_id = ?", company_id)
	} else {
		_, err = companyemployeeservice.DbContext.Exec("delete from "+companyemployeeservice.Table+" where company_id = ?", company_id)
	}
	if err != nil {
		log.Error("Error during deleting company employee objects for company object in database %v with value %v", err, company_id)
		return err
	}

	return nil
}
