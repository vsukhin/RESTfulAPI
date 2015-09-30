package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type CompanyRepository interface {
	CheckUserAccess(user_id int64, id int64) (allowed bool, err error)
	Get(id int64) (company *models.DtoCompany, err error)
	GetMeta(user_id int64) (company *models.ApiMetaCompany, err error)
	GetByUser(userid int64, filter string) (companies *[]models.ApiShortCompany, err error)
	GetByUnit(unitid int64) (companies *[]models.ApiShortCompany, err error)
	GetPrimaryByUser(userid int64) (company *models.DtoCompany, err error)
	GetPrimaryByUnit(unitid int64) (company *models.DtoCompany, err error)
	ClearPrimary(company *models.DtoCompany, trans *gorp.Transaction) (err error)
	SetArrays(company *models.DtoCompany, trans *gorp.Transaction) (err error)
	Create(company *models.DtoCompany, inTrans bool) (err error)
	Update(company *models.DtoCompany, inTrans bool) (err error)
	Deactivate(company *models.DtoCompany) (err error)
}

type CompanyService struct {
	CompanyCodeRepository     CompanyCodeRepository
	CompanyAddressRepository  CompanyAddressRepository
	CompanyBankRepository     CompanyBankRepository
	CompanyEmployeeRepository CompanyEmployeeRepository
	*Repository
}

func NewCompanyService(repository *Repository) *CompanyService {
	repository.DbContext.AddTableWithName(models.DtoCompany{}, repository.Table).SetKeys(true, "id")
	return &CompanyService{Repository: repository}
}

func (companyservice *CompanyService) CheckUserAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := companyservice.DbContext.SelectInt("select count(*) from "+companyservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking company object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (companyservice *CompanyService) Get(id int64) (company *models.DtoCompany, err error) {
	company = new(models.DtoCompany)
	err = companyservice.DbContext.SelectOne(company, "select * from "+companyservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting company object from database %v with value %v", err, id)
		return nil, err
	}

	return company, nil
}

func (companyservice *CompanyService) GetMeta(user_id int64) (company *models.ApiMetaCompany, err error) {
	company = new(models.ApiMetaCompany)
	company.Total, err = companyservice.DbContext.SelectInt("select count(*) from "+companyservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and active = 1", user_id)
	if err != nil {
		log.Error("Error during getting meta company object from database %v with value %v", err, user_id)
		return nil, err
	}

	return company, nil
}

func (companyservice *CompanyService) GetByUser(userid int64, filter string) (companies *[]models.ApiShortCompany, err error) {
	companies = new([]models.ApiShortCompany)
	_, err = companyservice.DbContext.Select(companies,
		"select id, shortname_rus as nameShortRus, shortname_eng as nameShortEng, unit_id as unitId, locked as `lock`, `primary` from "+
			companyservice.Table+" where unit_id = (select unit_id from users where id = ?) and active = 1"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit company object from database %v with value %v", err, userid)
		return nil, err
	}

	return companies, nil
}

func (companyservice *CompanyService) GetByUnit(unitid int64) (companies *[]models.ApiShortCompany, err error) {
	companies = new([]models.ApiShortCompany)
	_, err = companyservice.DbContext.Select(companies,
		"select id, shortname_rus as nameShortRus, shortname_eng as nameShortEng, unit_id as unitId, locked as `lock`, `primary` from "+companyservice.Table+
			" where unit_id = ? and active = 1", unitid)
	if err != nil {
		log.Error("Error during getting unit company object from database %v with value %v", err, unitid)
		return nil, err
	}

	return companies, nil
}

func (companyservice *CompanyService) GetPrimaryByUser(userid int64) (company *models.DtoCompany, err error) {
	company = new(models.DtoCompany)
	err = companyservice.DbContext.SelectOne(company, "select * from "+companyservice.Table+" where active = 1 and `primary` = 1"+
		" and unit_id = (select unit_id from users where id = ?)", userid)
	if err != nil {
		log.Error("Error during getting unit company object from database %v with value %v", err, userid)
		return nil, err
	}

	return company, nil
}

func (companyservice *CompanyService) GetPrimaryByUnit(unitid int64) (company *models.DtoCompany, err error) {
	company = new(models.DtoCompany)
	err = companyservice.DbContext.SelectOne(company, "select * from "+companyservice.Table+" where active = 1 and `primary` = 1"+
		" and unit_id = ?", unitid)
	if err != nil {
		log.Error("Error during getting unit company object from database %v with value %v", err, unitid)
		return nil, err
	}

	return company, nil
}

func (companyservice *CompanyService) ClearPrimary(company *models.DtoCompany, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("update "+companyservice.Table+" set `primary` = 0 where unit_id = ? and active = 1", company.Unit_ID)
	} else {
		_, err = companyservice.DbContext.Exec("update "+companyservice.Table+" set `primary` = 0 where unit_id = ? and active = 1", company.Unit_ID)
	}
	if err != nil {
		log.Error("Error during preparing company object in database %v", err)
		return err
	}

	return nil
}

func (companyservice *CompanyService) SetArrays(company *models.DtoCompany, trans *gorp.Transaction) (err error) {
	err = companyservice.CompanyCodeRepository.DeleteByCompany(company.ID, trans)
	if err != nil {
		log.Error("Error during setting company object in database %v with value %v", err, company.ID)
		return err
	}
	for _, dtocompanycode := range company.CompanyCodes {
		dtocompanycode.Company_ID = company.ID
		err = companyservice.CompanyCodeRepository.Create(&dtocompanycode, trans)
		if err != nil {
			log.Error("Error during setting company object in database %v with value %v", err, company.ID)
			return err
		}
	}
	err = companyservice.CompanyAddressRepository.DeleteByCompany(company.ID, trans)
	if err != nil {
		log.Error("Error during setting company object in database %v with value %v", err, company.ID)
		return err
	}
	for _, dtocompanyaddress := range company.CompanyAddresses {
		dtocompanyaddress.Company_ID = company.ID
		err = companyservice.CompanyAddressRepository.Create(&dtocompanyaddress, trans)
		if err != nil {
			log.Error("Error during setting company object in database %v with value %v", err, company.ID)
			return err
		}
	}
	err = companyservice.CompanyBankRepository.DeleteByCompany(company.ID, trans)
	if err != nil {
		log.Error("Error during setting company object in database %v with value %v", err, company.ID)
		return err
	}
	for _, dtocompanybank := range company.CompanyBanks {
		dtocompanybank.Company_ID = company.ID
		err = companyservice.CompanyBankRepository.Create(&dtocompanybank, trans)
		if err != nil {
			log.Error("Error during setting company object in database %v with value %v", err, company.ID)
			return err
		}
	}
	err = companyservice.CompanyEmployeeRepository.DeleteByCompany(company.ID, trans)
	if err != nil {
		log.Error("Error during setting company object in database %v with value %v", err, company.ID)
		return err
	}
	for _, dtocompanyemployee := range company.CompanyStaff {
		dtocompanyemployee.Company_ID = company.ID
		err = companyservice.CompanyEmployeeRepository.Create(&dtocompanyemployee, trans)
		if err != nil {
			log.Error("Error during setting company object in database %v with value %v", err, company.ID)
			return err
		}
	}

	return nil
}

func (companyservice *CompanyService) Create(company *models.DtoCompany, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = companyservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating company object in database %v", err)
			return err
		}
	}

	if company.Primary {
		err = companyservice.ClearPrimary(company, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		err = trans.Insert(company)
	} else {
		err = companyservice.DbContext.Insert(company)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating company object in database %v", err)
		return err
	}

	err = companyservice.SetArrays(company, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating company object in database %v", err)
			return err
		}
	}

	return nil
}

func (companyservice *CompanyService) Update(company *models.DtoCompany, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = companyservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating company object in database %v", err)
			return err
		}
	}

	if company.Primary {
		err = companyservice.ClearPrimary(company, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(company)
	} else {
		_, err = companyservice.DbContext.Update(company)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating company object in database %v with value %v", err, company.ID)
		return err
	}

	err = companyservice.SetArrays(company, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating company object in database %v", err)
			return err
		}
	}

	return nil
}

func (companyservice *CompanyService) Deactivate(company *models.DtoCompany) (err error) {
	_, err = companyservice.DbContext.Exec("update "+companyservice.Table+" set active = 0 where id = ?", company.ID)
	if err != nil {
		log.Error("Error during deactivating company object in database %v with value %v", err, company.ID)
		return err
	}

	return nil
}
