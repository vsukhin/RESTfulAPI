package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type CompanyCodeRepository interface {
	Get(id int64) (companycode *models.DtoCompanyCode, err error)
	GetByCompany(company_id int64) (companycodes *[]models.ViewApiCompanyCode, err error)
	Create(dtocompanycode *models.DtoCompanyCode, trans *gorp.Transaction) (err error)
	DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error)
}

type CompanyCodeService struct {
	*Repository
}

func NewCompanyCodeService(repository *Repository) *CompanyCodeService {
	repository.DbContext.AddTableWithName(models.DtoCompanyCode{}, repository.Table).SetKeys(true, "id")
	return &CompanyCodeService{Repository: repository}
}

func (companycodeservice *CompanyCodeService) Get(id int64) (companycode *models.DtoCompanyCode, err error) {
	companycode = new(models.DtoCompanyCode)
	err = companycodeservice.DbContext.SelectOne(companycode, "select * from "+companycodeservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company code object from database %v with value %v", err, id)
		return nil, err
	}

	return companycode, nil
}

func (companycodeservice *CompanyCodeService) GetByCompany(company_id int64) (companycodes *[]models.ViewApiCompanyCode, err error) {
	companycodes = new([]models.ViewApiCompanyCode)
	_, err = companycodeservice.DbContext.Select(companycodes,
		"select c.company_class_id, group_concat(c.code separator ',') as codes from "+companycodeservice.Table+
			" c inner join company_classes l on c.company_class_id =l.id  where c.company_id = ? and l.active = 1"+
			" group by company_class_id order by l.position asc", company_id)
	if err != nil {
		log.Error("Error during getting all company code object from database %v with value %v", err, company_id)
		return nil, err
	}

	return companycodes, nil
}

func (companycodeservice *CompanyCodeService) Create(dtocompanycode *models.DtoCompanyCode, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtocompanycode)
	} else {
		err = companycodeservice.DbContext.Insert(dtocompanycode)
	}
	if err != nil {
		log.Error("Error during creating company code object in database %v", err)
		return err
	}

	return nil
}

func (companycodeservice *CompanyCodeService) DeleteByCompany(company_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+companycodeservice.Table+" where company_id = ?", company_id)
	} else {
		_, err = companycodeservice.DbContext.Exec("delete from "+companycodeservice.Table+" where company_id = ?", company_id)
	}
	if err != nil {
		log.Error("Error during deleting company code objects for company object in database %v with value %v", err, company_id)
		return err
	}

	return nil
}
