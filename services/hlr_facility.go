package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type HLRFacilityRepository interface {
	Exists(orderid int64) (found bool, err error)
	Get(order_id int64) (hlrfacility *models.DtoHLRFacility, err error)
	SetArrays(hlrfacility *models.DtoHLRFacility, trans *gorp.Transaction) (err error)
	Create(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error)
	Update(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error)
	Save(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error)
}

type HLRFacilityService struct {
	MobileOperatorOperationRepository MobileOperatorOperationRepository
	ResultTableRepository             ResultTableRepository
	WorkTableRepository               WorkTableRepository
	*Repository
}

func NewHLRFacilityService(repository *Repository) *HLRFacilityService {
	repository.DbContext.AddTableWithName(models.DtoHLRFacility{}, repository.Table).SetKeys(false, "order_id")
	return &HLRFacilityService{Repository: repository}
}

func (hlrfacilityservice *HLRFacilityService) Exists(orderid int64) (found bool, err error) {
	var count int64
	count, err = hlrfacilityservice.DbContext.SelectInt("select count(*) from "+hlrfacilityservice.Table+
		" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during checking hlr facility object in database %v with value %v", err, orderid)
		return false, err
	}

	return count != 0, nil
}

func (hlrfacilityservice *HLRFacilityService) Get(orderid int64) (hlrfacility *models.DtoHLRFacility, err error) {
	hlrfacility = new(models.DtoHLRFacility)
	err = hlrfacilityservice.DbContext.SelectOne(hlrfacility, "select * from "+hlrfacilityservice.Table+" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during getting hlr facility object from database %v with value %v", err, orderid)
		return nil, err
	}

	return hlrfacility, nil
}

func (hlrfacilityservice *HLRFacilityService) SetArrays(hlrfacility *models.DtoHLRFacility, trans *gorp.Transaction) (err error) {
	err = hlrfacilityservice.MobileOperatorOperationRepository.DeleteByOrder(hlrfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
		return err
	}
	for _, dtooperator := range hlrfacility.EstimatedOperators {
		err = hlrfacilityservice.MobileOperatorOperationRepository.Create(&dtooperator, trans)
		if err != nil {
			log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
			return err
		}
	}
	err = hlrfacilityservice.ResultTableRepository.DeleteByOrder(hlrfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
		return err
	}
	for _, dtoresulttable := range hlrfacility.ResultTables {
		err = hlrfacilityservice.ResultTableRepository.Create(&dtoresulttable, trans)
		if err != nil {
			log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
			return err
		}
	}
	err = hlrfacilityservice.WorkTableRepository.DeleteByOrder(hlrfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
		return err
	}
	for _, dtoworktable := range hlrfacility.WorkTables {
		err = hlrfacilityservice.WorkTableRepository.Create(&dtoworktable, trans)
		if err != nil {
			log.Error("Error during setting hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
			return err
		}
	}

	return nil
}

func (hlrfacilityservice *HLRFacilityService) Create(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = hlrfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating hlr facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(hlrfacility)
	} else {
		err = hlrfacilityservice.DbContext.Insert(hlrfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating hlr facility object in database %v", err)
		return err
	}

	if !briefly {
		err = hlrfacilityservice.SetArrays(hlrfacility, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating hlr facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (hlrfacilityservice *HLRFacilityService) Update(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = hlrfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating hlr facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(hlrfacility)
	} else {
		_, err = hlrfacilityservice.DbContext.Update(hlrfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
		return err
	}

	if !briefly {
		err = hlrfacilityservice.SetArrays(hlrfacility, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating hlr facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (hlrfacilityservice *HLRFacilityService) Save(hlrfacility *models.DtoHLRFacility, briefly bool, inTrans bool) (err error) {
	count, err := hlrfacilityservice.DbContext.SelectInt("select count(*) from "+hlrfacilityservice.Table+
		" where order_id = ?", hlrfacility.Order_ID)
	if err != nil {
		log.Error("Error during saving hlr facility object in database %v with value %v", err, hlrfacility.Order_ID)
		return err
	}
	if count == 0 {
		err = hlrfacilityservice.Create(hlrfacility, briefly, inTrans)
	} else {
		err = hlrfacilityservice.Update(hlrfacility, briefly, inTrans)
	}

	return err
}
