package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type SMSEventRepository interface {
	Get(order_id int64, event_id int) (smsevent *models.DtoSMSEvent, err error)
	GetByOrder(order_id int64) (smsevents *[]models.ViewApiSMSEvent, err error)
	Create(dtosmsevent *models.DtoSMSEvent, trans *gorp.Transaction) (err error)
	DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error)
}

type SMSEventService struct {
	*Repository
}

func NewSMSEventService(repository *Repository) *SMSEventService {
	repository.DbContext.AddTableWithName(models.DtoSMSEvent{}, repository.Table).SetKeys(false, "order_id", "event_id")
	return &SMSEventService{Repository: repository}
}

func (smseventservice *SMSEventService) Get(order_id int64, event_id int) (smsevent *models.DtoSMSEvent, err error) {
	smsevent = new(models.DtoSMSEvent)
	err = smseventservice.DbContext.SelectOne(smsevent, "select * from "+smseventservice.Table+
		" where order_id = ? and event_id = ?", order_id, event_id)
	if err != nil {
		log.Error("Error during getting sms event object from database %v with value %v, %v", err, order_id, event_id)
		return nil, err
	}

	return smsevent, nil
}

func (smseventservice *SMSEventService) GetByOrder(order_id int64) (smsevents *[]models.ViewApiSMSEvent, err error) {
	smsevents = new([]models.ViewApiSMSEvent)
	_, err = smseventservice.DbContext.Select(smsevents,
		"select event_id from "+smseventservice.Table+" where order_id = ?", order_id)
	if err != nil {
		log.Error("Error during getting all sms event object from database %v with value %v", err, order_id)
		return nil, err
	}

	return smsevents, nil
}

func (smseventservice *SMSEventService) Create(dtosmsevent *models.DtoSMSEvent, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(dtosmsevent)
	} else {
		err = smseventservice.DbContext.Insert(dtosmsevent)
	}
	if err != nil {
		log.Error("Error during creating sms event object in database %v", err)
		return err
	}

	return nil
}

func (smseventservice *SMSEventService) DeleteByOrder(order_id int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+smseventservice.Table+" where order_id = ?", order_id)
	} else {
		_, err = smseventservice.DbContext.Exec("delete from "+smseventservice.Table+" where order_id = ?", order_id)
	}
	if err != nil {
		log.Error("Error during deleting sms event objects for order object in database %v with value %v", err, order_id)
		return err
	}

	return nil
}
