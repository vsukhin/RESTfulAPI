package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type SMSFacilityRepository interface {
	Exists(orderid int64) (found bool, err error)
	Get(order_id int64) (smsfacility *models.DtoSMSFacility, err error)
	SetArrays(smsfacility *models.DtoSMSFacility, trans *gorp.Transaction) (err error)
	Create(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error)
	Update(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error)
	Save(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error)
}

type SMSFacilityService struct {
	MobileOperatorOperationRepository MobileOperatorOperationRepository
	SMSPeriodRepository               SMSPeriodRepository
	SMSEventRepository                SMSEventRepository
	ResultTableRepository             ResultTableRepository
	WorkTableRepository               WorkTableRepository
	*Repository
}

func NewSMSFacilityService(repository *Repository) *SMSFacilityService {
	repository.DbContext.AddTableWithName(models.DtoSMSFacility{}, repository.Table).SetKeys(false, "order_id")
	return &SMSFacilityService{Repository: repository}
}

func (smsfacilityservice *SMSFacilityService) Exists(orderid int64) (found bool, err error) {
	var count int64
	count, err = smsfacilityservice.DbContext.SelectInt("select count(*) from "+smsfacilityservice.Table+
		" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during checking sms facility object in database %v with value %v", err, orderid)
		return false, err
	}

	return count != 0, nil
}

func (smsfacilityservice *SMSFacilityService) Get(orderid int64) (smsfacility *models.DtoSMSFacility, err error) {
	smsfacility = new(models.DtoSMSFacility)
	err = smsfacilityservice.DbContext.SelectOne(smsfacility, "select * from "+smsfacilityservice.Table+" where order_id = ?", orderid)
	if err != nil {
		log.Error("Error during getting sms facility object from database %v with value %v", err, orderid)
		return nil, err
	}

	return smsfacility, nil
}

func (smsfacilityservice *SMSFacilityService) SetArrays(smsfacility *models.DtoSMSFacility, trans *gorp.Transaction) (err error) {
	err = smsfacilityservice.MobileOperatorOperationRepository.DeleteByOrder(smsfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	for _, dtooperator := range smsfacility.EstimatedOperators {
		err = smsfacilityservice.MobileOperatorOperationRepository.Create(&dtooperator, trans)
		if err != nil {
			log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
			return err
		}
	}
	err = smsfacilityservice.SMSPeriodRepository.DeleteByOrder(smsfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	for _, smsevent := range smsfacility.Periods {
		err = smsfacilityservice.SMSPeriodRepository.Create(&smsevent, trans)
		if err != nil {
			log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
			return err
		}
	}
	err = smsfacilityservice.SMSEventRepository.DeleteByOrder(smsfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	for _, smsevent := range smsfacility.Events {
		err = smsfacilityservice.SMSEventRepository.Create(&smsevent, trans)
		if err != nil {
			log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
			return err
		}
	}
	err = smsfacilityservice.ResultTableRepository.DeleteByOrder(smsfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	for _, dtoresulttable := range smsfacility.ResultTables {
		err = smsfacilityservice.ResultTableRepository.Create(&dtoresulttable, trans)
		if err != nil {
			log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
			return err
		}
	}
	err = smsfacilityservice.WorkTableRepository.DeleteByOrder(smsfacility.Order_ID, trans)
	if err != nil {
		log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	for _, dtoworktable := range smsfacility.WorkTables {
		err = smsfacilityservice.WorkTableRepository.Create(&dtoworktable, trans)
		if err != nil {
			log.Error("Error during setting sms facility object in database %v with value %v", err, smsfacility.Order_ID)
			return err
		}
	}

	return nil
}

func (smsfacilityservice *SMSFacilityService) Create(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = smsfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating sms facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		err = trans.Insert(smsfacility)
	} else {
		err = smsfacilityservice.DbContext.Insert(smsfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating sms facility object in database %v", err)
		return err
	}

	if !briefly {
		err = smsfacilityservice.SetArrays(smsfacility, trans)
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
			log.Error("Error during creating sms facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (smsfacilityservice *SMSFacilityService) Update(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = smsfacilityservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating sms facility object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(smsfacility)
	} else {
		_, err = smsfacilityservice.DbContext.Update(smsfacility)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}

	if !briefly {
		err = smsfacilityservice.SetArrays(smsfacility, trans)
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
			log.Error("Error during updating sms facility object in database %v", err)
			return err
		}
	}

	return nil
}

func (smsfacilityservice *SMSFacilityService) Save(smsfacility *models.DtoSMSFacility, briefly bool, inTrans bool) (err error) {
	count, err := smsfacilityservice.DbContext.SelectInt("select count(*) from "+smsfacilityservice.Table+
		" where order_id = ?", smsfacility.Order_ID)
	if err != nil {
		log.Error("Error during saving sms facility object in database %v with value %v", err, smsfacility.Order_ID)
		return err
	}
	if count == 0 {
		err = smsfacilityservice.Create(smsfacility, briefly, inTrans)
	} else {
		err = smsfacilityservice.Update(smsfacility, briefly, inTrans)
	}

	return err
}
