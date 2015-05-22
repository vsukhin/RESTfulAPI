package services

import (
	"application/models"
)

type CompanyTypeRepository interface {
	Get(id int) (companytype *models.DtoCompanyType, err error)
	GetAll() (companytypes *[]models.ApiCompanyType, err error)
}

type CompanyTypeService struct {
	*Repository
}

func NewCompanyTypeService(repository *Repository) *CompanyTypeService {
	repository.DbContext.AddTableWithName(models.DtoCompanyType{}, repository.Table).SetKeys(false, "id")
	return &CompanyTypeService{Repository: repository}
}

func (companytypeservice *CompanyTypeService) Get(id int) (companytype *models.DtoCompanyType, err error) {
	companytype = new(models.DtoCompanyType)
	err = companytypeservice.DbContext.SelectOne(companytype, "select * from "+companytypeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company type object from database %v with value %v", err, id)
		return nil, err
	}

	return companytype, nil
}

func (companytypeservice *CompanyTypeService) GetAll() (companytypes *[]models.ApiCompanyType, err error) {
	companytypes = new([]models.ApiCompanyType)
	_, err = companytypeservice.DbContext.Select(companytypes,
		"select id, fullname_rus, fullname_eng, shortname_rus, shortname_eng, position from "+companytypeservice.Table+
			" where active = 1 order by position asc")
	if err != nil {
		log.Error("Error during getting all company type object from database %v", err)
		return nil, err
	}

	return companytypes, nil
}
