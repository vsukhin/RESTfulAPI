package services

import (
	"application/models"
)

type CompanyClassRepository interface {
	Get(id int) (companyclass *models.DtoCompanyClass, err error)
	GetAll(filter string) (companyclasses *[]models.ApiCompanyClass, err error)
}

type CompanyClassService struct {
	*Repository
}

func NewCompanyClassService(repository *Repository) *CompanyClassService {
	repository.DbContext.AddTableWithName(models.DtoCompanyClass{}, repository.Table).SetKeys(false, "id")
	return &CompanyClassService{Repository: repository}
}

func (companyclassservice *CompanyClassService) Get(id int) (companyclass *models.DtoCompanyClass, err error) {
	companyclass = new(models.DtoCompanyClass)
	err = companyclassservice.DbContext.SelectOne(companyclass, "select * from "+companyclassservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company class object from database %v with value %v", err, id)
		return nil, err
	}

	return companyclass, nil
}

func (companyclassservice *CompanyClassService) GetAll(filter string) (companyclasses *[]models.ApiCompanyClass, err error) {
	companyclasses = new([]models.ApiCompanyClass)
	_, err = companyclassservice.DbContext.Select(companyclasses,
		"select id, fullname as nameFull, shortname as nameShort, format, required, visible as outward,"+
			" multiple as multiplicity, position, not active as del from "+companyclassservice.Table+filter)
	if err != nil {
		log.Error("Error during getting all company class object from database %v", err)
		return nil, err
	}

	return companyclasses, nil
}
