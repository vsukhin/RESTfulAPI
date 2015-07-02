package services

import (
	"application/models"
	"errors"
)

type SMSTableRepository interface {
	Get(customertable_id int64) (smstables *models.ApiSMSTable, err error)
	GetAll(user_id int64) (smstables *[]models.ApiSMSTable, err error)
}

type SMSTableService struct {
	FacilityTableRepository FacilityTableRepository
	*Repository
}

func NewSMSTableService(repository *Repository) *SMSTableService {
	return &SMSTableService{Repository: repository}
}

func (smstableservice *SMSTableService) Get(customertable_id int64) (smstable *models.ApiSMSTable, err error) {
	smstable = new(models.ApiSMSTable)

	table := new(models.DtoCustomerTable)
	err = smstableservice.DbContext.SelectOne(table, "select * from "+smstableservice.Table+
		" where id = ? and active = 1 and permanent = 1", customertable_id)
	if err != nil {
		log.Error("Error during getting sms facility table in database %v with value %v", err, customertable_id)
		return nil, err
	}

	mobilephones, goodmobilephones, err := smstableservice.FacilityTableRepository.GetColumnsByCustomerTable(customertable_id,
		models.COLUMN_TYPE_MOBILE_PHONE)
	if err != nil {
		return nil, err
	}

	messages, goodmessages, err := smstableservice.FacilityTableRepository.GetColumnsByCustomerTable(customertable_id,
		models.COLUMN_TYPE_SMS)
	if err != nil {
		return nil, err
	}

	smssenders, goodsmssenders, err := smstableservice.FacilityTableRepository.GetColumnsByCustomerTable(customertable_id,
		models.COLUMN_TYPE_SMS_SENDER)
	if err != nil {
		return nil, err
	}

	birthdays, goodbirthdays, err := smstableservice.FacilityTableRepository.GetColumnsByCustomerTable(customertable_id,
		models.COLUMN_TYPE_BIRTHDAY)
	if err != nil {
		return nil, err
	}

	smstable.ID = table.ID
	smstable.Name = table.Name
	smstable.UnitID = table.UnitID
	smstable.TypeID = table.TypeID
	for _, mobilephone := range *mobilephones {
		if goodmobilephones[mobilephone.ID] != 0 {
			smstable.MobilePhones = append(smstable.MobilePhones, *models.NewApiTableColumn(mobilephone.ID, mobilephone.Name,
				mobilephone.Column_Type_ID, mobilephone.Position))
		}
	}
	if len(smstable.MobilePhones) == 0 {
		log.Error("Mobile phones was not found for customer table %v", customertable_id)
		return nil, errors.New("Mobile phones not found")
	}
	for _, message := range *messages {
		if goodmessages[message.ID] != 0 {
			smstable.Messages = append(smstable.Messages, *models.NewApiTableColumn(message.ID, message.Name,
				message.Column_Type_ID, message.Position))
		}
	}
	for _, smssender := range *smssenders {
		if goodsmssenders[smssender.ID] != 0 {
			smstable.SMSSenders = append(smstable.SMSSenders, *models.NewApiTableColumn(smssender.ID, smssender.Name,
				smssender.Column_Type_ID, smssender.Position))
		}
	}
	for _, birthday := range *birthdays {
		if goodbirthdays[birthday.ID] != 0 {
			smstable.Birthdays = append(smstable.Birthdays, *models.NewApiTableColumn(birthday.ID, birthday.Name,
				birthday.Column_Type_ID, birthday.Position))
		}
	}

	return smstable, nil
}

func (smstableservice *SMSTableService) GetAll(user_id int64) (smstables *[]models.ApiSMSTable, err error) {
	smstables = new([]models.ApiSMSTable)

	tables := new([]models.DtoCustomerTable)
	_, err = smstableservice.DbContext.Select(tables, "select * from "+smstableservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and active = 1 and permanent = 1", user_id)
	if err != nil {
		log.Error("Error during getting all sms facility tables in database %v with value %v", err, user_id)
		return nil, err
	}

	mobilephones, goodmobilephones, err := smstableservice.FacilityTableRepository.GetColumnsByType(user_id, models.COLUMN_TYPE_MOBILE_PHONE)
	if err != nil {
		return nil, err
	}

	messages, goodmessages, err := smstableservice.FacilityTableRepository.GetColumnsByType(user_id, models.COLUMN_TYPE_SMS)
	if err != nil {
		return nil, err
	}

	smssenders, goodsmssenders, err := smstableservice.FacilityTableRepository.GetColumnsByType(user_id, models.COLUMN_TYPE_SMS_SENDER)
	if err != nil {
		return nil, err
	}

	birthdays, goodbirthdays, err := smstableservice.FacilityTableRepository.GetColumnsByType(user_id, models.COLUMN_TYPE_BIRTHDAY)
	if err != nil {
		return nil, err
	}

	for _, table := range *tables {
		smstable := new(models.ApiSMSTable)
		smstable.ID = table.ID
		smstable.Name = table.Name
		smstable.UnitID = table.UnitID
		smstable.TypeID = table.TypeID
		for _, mobilephone := range *mobilephones {
			if mobilephone.Customer_Table_ID == table.ID && goodmobilephones[mobilephone.ID] != 0 {
				smstable.MobilePhones = append(smstable.MobilePhones, *models.NewApiTableColumn(mobilephone.ID, mobilephone.Name,
					mobilephone.Column_Type_ID, mobilephone.Position))
			}
		}
		if len(smstable.MobilePhones) == 0 {
			continue
		}
		for _, message := range *messages {
			if message.Customer_Table_ID == table.ID && goodmessages[message.ID] != 0 {
				smstable.Messages = append(smstable.Messages, *models.NewApiTableColumn(message.ID, message.Name,
					message.Column_Type_ID, message.Position))
			}
		}
		for _, smssender := range *smssenders {
			if smssender.Customer_Table_ID == table.ID && goodsmssenders[smssender.ID] != 0 {
				smstable.SMSSenders = append(smstable.SMSSenders, *models.NewApiTableColumn(smssender.ID, smssender.Name,
					smssender.Column_Type_ID, smssender.Position))
			}
		}
		for _, birthday := range *birthdays {
			if birthday.Customer_Table_ID == table.ID && goodbirthdays[birthday.ID] != 0 {
				smstable.Birthdays = append(smstable.Birthdays, *models.NewApiTableColumn(birthday.ID, birthday.Name,
					birthday.Column_Type_ID, birthday.Position))
			}
		}

		*smstables = append(*smstables, *smstable)
	}

	return smstables, nil
}
