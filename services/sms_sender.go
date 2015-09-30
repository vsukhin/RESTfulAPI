package services

import (
	"application/models"
	"fmt"
)

type SMSSenderRepository interface {
	CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error)
	Exists(name string, supplier_id int64) (found bool, err error)
	Belongs(dtotablecolumn *models.DtoTableColumn, unit_id int64, supplier_id int64) (found bool, err error)
	Get(id int64) (smssender *models.DtoSMSSender, err error)
	GetMeta(user_id int64) (smssender *models.ApiMetaSMSSender, err error)
	GetByUser(userid int64, filter string) (smssenders *[]models.ApiLongSMSSender, err error)
	GetByUnit(unitid int64) (smssenders *[]models.ApiLongSMSSender, err error)
	Create(smssender *models.DtoSMSSender) (err error)
	Update(smssender *models.DtoSMSSender) (err error)
	Deactivate(smssender *models.DtoSMSSender) (err error)
}

type SMSSenderService struct {
	*Repository
}

func NewSMSSenderService(repository *Repository) *SMSSenderService {
	repository.DbContext.AddTableWithName(models.DtoSMSSender{}, repository.Table).SetKeys(true, "id")
	return &SMSSenderService{Repository: repository}
}

func (smssenderservice *SMSSenderService) CheckCustomerAccess(user_id int64, id int64) (allowed bool, err error) {
	count, err := smssenderservice.DbContext.SelectInt("select count(*) from "+smssenderservice.Table+
		" where id = ? and unit_id = (select unit_id from users where id = ?)", id, user_id)
	if err != nil {
		log.Error("Error during checking sms sender object from database %v with value %v, %v", err, user_id, id)
		return false, err
	}

	return count != 0, nil
}

func (smssenderservice *SMSSenderService) Exists(name string, supplier_id int64) (found bool, err error) {
	var count int64
	count, err = smssenderservice.DbContext.SelectInt(
		"select count(*) from "+smssenderservice.Table+" where name = ? and supplier_id  = ? and active = 1", name, supplier_id)
	if err != nil {
		log.Error("Error during getting sms sender object from database %v with value %v, %v", err, name, supplier_id)
		return false, err
	}

	return count != 0, nil
}

func (smssenderservice *SMSSenderService) Belongs(dtotablecolumn *models.DtoTableColumn, unit_id int64, supplier_id int64) (found bool, err error) {
	var count int64
	count, err = smssenderservice.DbContext.SelectInt("select count(*) from table_data where active = 1 and customer_table_id = ? and field"+
		fmt.Sprintf("%v", dtotablecolumn.FieldNum)+
		" not in (select name from sms_senders where active = 1 and registered = 1 and unit_id = ? and supplier_id = ?)",
		dtotablecolumn.Customer_Table_ID, unit_id, supplier_id)
	if err != nil {
		log.Error("Error during getting sms sender object from database %v with value %v, %v, %v", err, dtotablecolumn.ID, unit_id, supplier_id)
		return false, err
	}

	return count == 0, nil
}

func (smssenderservice *SMSSenderService) Get(id int64) (smssender *models.DtoSMSSender, err error) {
	smssender = new(models.DtoSMSSender)
	err = smssenderservice.DbContext.SelectOne(smssender, "select * from "+smssenderservice.Table+" where id = ?", id)
	if err != nil {
		log.Error("Error during getting sms sender object from database %v with value %v", err, id)
		return nil, err
	}

	return smssender, nil
}

func (smssenderservice *SMSSenderService) GetMeta(user_id int64) (smssender *models.ApiMetaSMSSender, err error) {
	smssender = new(models.ApiMetaSMSSender)
	smssender.Total, err = smssenderservice.DbContext.SelectInt("select count(*) from "+smssenderservice.Table+
		" where active = 1 and unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta sms sender object from database %v with value %v", err, user_id)
		return nil, err
	}
	smssender.NumOfDeleted, err = smssenderservice.DbContext.SelectInt("select count(*) from "+smssenderservice.Table+
		" where active = 0 and unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta sms sender object from database %v with value %v", err, user_id)
		return nil, err
	}
	smssender.NumOfNew, err = smssenderservice.DbContext.SelectInt("select count(*) from "+smssenderservice.Table+
		" where active =1 and registered = 0 and unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta sms sender object from database %v with value %v", err, user_id)
		return nil, err
	}
	_, err = smssenderservice.DbContext.Select(&smssender.Suppliers, "select supplier_id, count(*) as total from "+smssenderservice.Table+
		" where active = 1 and unit_id = (select unit_id from users where id = ?) group by supplier_id", user_id)
	if err != nil {
		log.Error("Error during getting meta sms sender object from database %v with value %v", err, user_id)
		return nil, err
	}

	return smssender, nil
}

func (smssenderservice *SMSSenderService) GetByUser(userid int64, filter string) (smssenders *[]models.ApiLongSMSSender, err error) {
	smssenders = new([]models.ApiLongSMSSender)
	_, err = smssenderservice.DbContext.Select(smssenders, "select id, name, registered, supplier_id as supplierId,"+
		" date_format(planned_begin, '%Y-%m-%d') as createBegin, date_format(planned_end, '%Y-%m-%d') as createEnd,"+
		" date_format(actual_begin, '%Y-%m-%d') as begin, date_format(actual_end, '%Y-%m-%d') as end, renew as autoRenew,"+
		" not active as del from "+smssenderservice.Table+
		" where unit_id = (select unit_id from users where id = ?)"+filter, userid)
	if err != nil {
		log.Error("Error during getting unit sms sender object from database %v with value %v", err, userid)
		return nil, err
	}

	return smssenders, nil
}

func (smssenderservice *SMSSenderService) GetByUnit(unitid int64) (smssenders *[]models.ApiLongSMSSender, err error) {
	smssenders = new([]models.ApiLongSMSSender)
	_, err = smssenderservice.DbContext.Select(smssenders, "select id, name, registered, supplier_id as supplierId,"+
		" date_format(planned_begin, '%Y-%m-%d') as createBegin, date_format(planned_end, '%Y-%m-%d') as createEnd,"+
		" date_format(actual_begin, '%Y-%m-%d') as begin, date_format(actual_end, '%Y-%m-%d') as end, renew as autoRenew,"+
		" not active as del from "+smssenderservice.Table+" where unit_id = ?", unitid)
	if err != nil {
		log.Error("Error during getting unit sms sender object from database %v with value %v", err, unitid)
		return nil, err
	}

	return smssenders, nil
}

func (smssenderservice *SMSSenderService) Create(smssender *models.DtoSMSSender) (err error) {
	err = smssenderservice.DbContext.Insert(smssender)
	if err != nil {
		log.Error("Error during creating sms sender object in database %v", err)
		return err
	}

	return nil
}

func (smssenderservice *SMSSenderService) Update(smssender *models.DtoSMSSender) (err error) {
	_, err = smssenderservice.DbContext.Update(smssender)
	if err != nil {
		log.Error("Error during updating sms sender object in database %v with value %v", err, smssender.ID)
		return err
	}

	return nil
}

func (smssenderservice *SMSSenderService) Deactivate(smssender *models.DtoSMSSender) (err error) {
	_, err = smssenderservice.DbContext.Exec("update "+smssenderservice.Table+" set active = 0 where id = ?", smssender.ID)
	if err != nil {
		log.Error("Error during deactivating sms sender object in database %v with value %v", err, smssender.ID)
		return err
	}

	return nil
}
