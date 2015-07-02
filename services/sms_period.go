package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type SMSPeriodRepository interface {
	Get(order_id int64, period_id int) (smsperiod *models.DtoSMSPeriod, err error)
	GetByOrder(order_id int64) (smsperiods *[]models.ViewApiSMSPeriod, err error)
	Create(dtosmsperiod *models.DtoSMSPeriod, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type SMSPeriodService struct {
	*Repository
}

func NewSMSPeriodService(repository *Repository) *SMSPeriodService {
	repository.DbContext.AddTableWithName(models.DtoSMSPeriod{}, repository.Table).SetKeys(false, "order_id", "period_id")
	return &SMSPeriodService{Repository: repository}
}

func (smsperiodservice *SMSPeriodService) Get(order_id int64, period_id int) (smsperiod *models.DtoSMSPeriod, err error) {
	smsperiod = new(models.DtoSMSPeriod)
	err = smsperiodservice.DbContext.SelectOne(smsperiod, "select * from "+smsperiodservice.Table+
		" where order_id = ? and period_id = ?", order_id, period_id)
	if err != nil {
		log.Error("Error during getting sms period object from database %v with value %v, %v", err, order_id, period_id)
		return nil, err
	}

	return smsperiod, nil
}

func (smsperiodservice *SMSPeriodService) GetByOrder(order_id int64) (smsperiods *[]models.ViewApiSMSPeriod, err error) {
	smsperiods = new([]models.ViewApiSMSPeriod)
	_, err = smsperiodservice.DbContext.Select(smsperiods,
		"select period_id from "+smsperiodservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all sms period object from database %v with value %v", err, order_id)
		return nil, err
	}

	return smsperiods, nil
}

func (smsperiodservice *SMSPeriodService) Create(dtosmsperiod *models.DtoSMSPeriod, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtosmsperiod)
	} else {
		err = smsperiodservice.DbContext.Insert(dtosmsperiod)
	}
	if err != nil {
		log.Error("Error during creating sms period object in database %v", err)
		return err
	}

	return nil
}

func (smsperiodservice *SMSPeriodService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+smsperiodservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = smsperiodservice.DbContext.Exec("delete from "+smsperiodservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting sms period objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
