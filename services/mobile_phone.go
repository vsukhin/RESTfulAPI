package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
)

type MobilePhoneRepository interface {
	Exists(phone string) (found bool, err error)
	Get(phone string) (dtomobilephone *models.DtoMobilePhone, err error)
	GetByUser(userid int64) (mobilephones *[]models.DtoMobilePhone, err error)
	Create(mobilephone *models.DtoMobilePhone, trans *gorp.Transaction) (err error)
	Update(mobilephone *models.DtoMobilePhone, trans *gorp.Transaction) (err error)
	Delete(phone string, trans *gorp.Transaction) (err error)
	DeleteByUser(userid int64, trans *gorp.Transaction) (err error)
}

type MobilePhoneService struct {
	*Repository
}

func NewMobilePhoneService(repository *Repository) *MobilePhoneService {
	repository.DbContext.AddTableWithName(models.DtoMobilePhone{}, repository.Table).SetKeys(false, "phone")
	return &MobilePhoneService{
		repository,
	}
}

func (mobilephoneservice *MobilePhoneService) Exists(phone string) (found bool, err error) {
	var count int64
	count, err = mobilephoneservice.DbContext.SelectInt("select count(*) from "+mobilephoneservice.Table+" where phone = ?", phone)
	if err != nil {
		log.Error("Error during getting mobile phone object from database %v with value %v", err, phone)
		return false, err
	}

	return count != 0, nil
}

func (mobilephoneservice *MobilePhoneService) Get(phone string) (dtomobilephone *models.DtoMobilePhone, err error) {
	dtomobilephone = new(models.DtoMobilePhone)
	err = mobilephoneservice.DbContext.SelectOne(dtomobilephone, "select * from "+mobilephoneservice.Table+" where phone = ?", phone)
	if err != nil {
		log.Error("Error during getting mobile phone object from database %v with value %v", err, phone)
		return nil, err
	}

	return dtomobilephone, nil
}

func (mobilephoneservice *MobilePhoneService) GetByUser(userid int64) (mobilephones *[]models.DtoMobilePhone, err error) {
	mobilephones = new([]models.DtoMobilePhone)
	_, err = mobilephoneservice.DbContext.Select(mobilephones, "select * from "+mobilephoneservice.Table+" where user_id = ?", userid)
	if err != nil {
		log.Error("Error during getting mobile phone object from database %v with value %v", err, userid)
		return nil, err
	}

	return mobilephones, nil
}

func (mobilephoneservice *MobilePhoneService) Create(mobilephone *models.DtoMobilePhone, trans *gorp.Transaction) (err error) {
	if trans != nil {
		err = trans.Insert(mobilephone)
	} else {
		err = mobilephoneservice.DbContext.Insert(mobilephone)
	}
	if err != nil {
		log.Error("Error during creating mobile phone object in database %v", err)
		return err
	}

	return nil
}

func (mobilephoneservice *MobilePhoneService) Update(mobilephone *models.DtoMobilePhone, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Update(mobilephone)
	} else {
		_, err = mobilephoneservice.DbContext.Update(mobilephone)
	}
	if err != nil {
		log.Error("Error during updating mobile phone object in database %v with value %v", err, mobilephone.Phone)
		return err
	}

	return nil
}

func (mobilephoneservice *MobilePhoneService) Delete(phone string, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+mobilephoneservice.Table+" where phone = ?", phone)
	} else {
		_, err = mobilephoneservice.DbContext.Exec("delete from "+mobilephoneservice.Table+" where phone = ?", phone)
	}
	if err != nil {
		log.Error("Error during deleting mobile phone object in database %v with value %v", err, phone)
		return err
	}

	return nil
}

func (mobilephoneservice *MobilePhoneService) DeleteByUser(userid int64, trans *gorp.Transaction) (err error) {
	if trans != nil {
		_, err = trans.Exec("delete from "+mobilephoneservice.Table+" where user_id = ?", userid)
	} else {
		_, err = mobilephoneservice.DbContext.Exec("delete from "+mobilephoneservice.Table+" where user_id = ?", userid)
	}
	if err != nil {
		log.Error("Error during deleting mobile phone object in database %v with value %v", err, userid)
		return err
	}

	return nil
}
