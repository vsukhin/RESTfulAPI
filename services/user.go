package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
	"time"
)

type UserRepository interface {
	CheckReportAccess(user_id int64) (allowed bool, err error)
	CheckUnitAccess(user_id int64) (allowed bool, err error)
	GetUserArrays(user *models.DtoUser) (*models.DtoUser, error)
	FindByLogin(login string) (user *models.DtoUser, err error)
	FindByCode(code string) (user *models.DtoUser, err error)
	Get(userid int64) (user *models.DtoUser, err error)
	GetAll(filter string) (users *[]models.ApiUserShort, err error)
	GetAllByUser(user_id int64, filter string) (users *[]models.ApiSearchUnitUser, err error)
	GetByUnit(unitid int64) (users *[]models.ApiUserTiny, err error)
	GetMeta() (usermeta *models.ApiUserMeta, err error)
	GetMetaByUser(user_id int64) (metaunituser *models.ApiMetaUnitUser, err error)
	GetMetaDashboard(user_id int64) (metadashboard *models.ApiMetaDashboard, err error)
	InitUnit(trans *gorp.Transaction) (unitid int64, err error)
	Create(user *models.DtoUser, inTrans bool) (err error)
	Update(user *models.DtoUser, briefly bool, inTrans bool) (err error)
	UpdateProfile(user *models.DtoUser) (err error)
	UpdatePassword(user *models.DtoUser) (err error)
	UpdateEmails(user *models.DtoUser, inTrans bool) (err error)
	UpdateMobilePhones(user *models.DtoUser, inTrans bool) (err error)
	Delete(userid int64, inTrans bool) (err error)
}

type UserService struct {
	SessionRepository     SessionRepository
	EmailRepository       EmailRepository
	UnitRepository        UnitRepository
	GroupRepository       GroupRepository
	MessageRepository     MessageRepository
	MobilePhoneRepository MobilePhoneRepository
	DeviceRepository      DeviceRepository
	*Repository
}

func NewUserService(repository *Repository) *UserService {
	repository.DbContext.AddTableWithName(models.DtoUser{}, repository.Table).SetKeys(true, "id")
	return &UserService{Repository: repository}
}

func (userservice *UserService) CheckReportAccess(user_id int64) (allowed bool, err error) {
	count, err := userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+
		" where id = ? and reportAccess = 1", user_id)
	if err != nil {
		log.Error("Error during checking report user object from database %v with value %v", err, user_id)
		return false, err
	}

	return count != 0, nil
}

func (userservice *UserService) CheckUnitAccess(user_id int64) (allowed bool, err error) {
	count, err := userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+
		" where id = ? and unitAdmin = 1", user_id)
	if err != nil {
		log.Error("Error during checking unit user object from database %v with value %v", err, user_id)
		return false, err
	}

	return count != 0, nil
}

func (userservice *UserService) GetUserArrays(user *models.DtoUser) (*models.DtoUser, error) {
	roles, err := userservice.GroupRepository.GetByUser(user.ID)
	if err != nil {
		log.Error("Error during getting user roles object from database %v with value %v", err, user.ID)
		return nil, err
	}
	user.Roles = *roles

	emails, err := userservice.EmailRepository.GetByUser(user.ID)
	if err != nil {
		log.Error("Error during getting user email object from database %v with value %v", err, user.ID)
		return nil, err
	}
	user.Emails = emails

	phones, err := userservice.MobilePhoneRepository.GetByUser(user.ID)
	if err != nil {
		log.Error("Error during getting user mobile phone object from database %v with value %v", err, user.ID)
		return nil, err
	}
	user.MobilePhones = phones

	return user, nil
}

func (userservice *UserService) FindByLogin(login string) (user *models.DtoUser, err error) {
	user = new(models.DtoUser)
	err = userservice.DbContext.SelectOne(user,
		"select * from "+userservice.Table+" where id in (select user_id from emails where `primary` = 1 and confirmed = 1 and email = ?) or"+
			" id in (select user_id from mobile_phones where `primary` = 1 and confirmed = 1 and phone = ?)", login, login)
	if err != nil {
		log.Error("Error during finding user object in database %v with value %v", err, login)
		return nil, err
	}

	return userservice.GetUserArrays(user)
}

func (userservice *UserService) FindByCode(code string) (user *models.DtoUser, err error) {
	user = new(models.DtoUser)
	err = userservice.DbContext.SelectOne(user, "select * from "+userservice.Table+" where code = ?", code)
	if err != nil {
		log.Error("Error during finding user object in database %v with value %v", err, code)
		return nil, err
	}

	return userservice.GetUserArrays(user)
}

func (userservice *UserService) Get(userid int64) (user *models.DtoUser, err error) {
	user = new(models.DtoUser)
	err = userservice.DbContext.SelectOne(user, "select * from "+userservice.Table+" where id = ?", userid)
	if err != nil {
		log.Error("Error during getting user object from database %v with value %v", err, userid)
		return nil, err
	}

	return userservice.GetUserArrays(user)
}

func (userservice *UserService) GetAll(filter string) (users *[]models.ApiUserShort, err error) {
	users = new([]models.ApiUserShort)
	_, err = userservice.DbContext.Select(users, "select id, not active as blocked, confirmed, lastLogin as lastLoginAt,"+
		" surname, name, middleName from "+userservice.Table+filter)
	if err != nil {
		log.Error("Error during getting user objects from database %v", err)
		return nil, err
	}

	return users, nil
}

func (userservice *UserService) GetAllByUser(user_id int64, filter string) (users *[]models.ApiSearchUnitUser, err error) {
	users = new([]models.ApiSearchUnitUser)
	_, err = userservice.DbContext.Select(users, "select id, not active as blocked, confirmed, lastLogin as lastLoginAt,"+
		" unitAdmin, surname, name, middleName from "+userservice.Table+" where unit_id = (select unit_id from users where id = ?)"+filter,
		user_id)
	if err != nil {
		log.Error("Error during getting unit user objects from database %v", err)
		return nil, err
	}

	return users, nil
}

func (userservice *UserService) GetByUnit(unitid int64) (users *[]models.ApiUserTiny, err error) {
	users = new([]models.ApiUserTiny)
	_, err = userservice.DbContext.Select(users, "select id from "+userservice.Table+" where unit_id = ?", unitid)
	if err != nil {
		log.Error("Error during getting all user objects from database %v with value%v", err, unitid)
		return nil, err
	}

	return users, nil
}

func (userservice *UserService) GetMeta() (usermeta *models.ApiUserMeta, err error) {
	usermeta = new(models.ApiUserMeta)
	usermeta.NumOfRows, err = userservice.DbContext.SelectInt("select count(*) from " + userservice.Table)
	if err != nil {
		log.Error("Error during getting meta user object from database %v", err)
		return nil, err
	}

	return usermeta, nil
}

func (userservice *UserService) GetMetaByUser(user_id int64) (metaunituser *models.ApiMetaUnitUser, err error) {
	metaunituser = new(models.ApiMetaUnitUser)
	metaunituser.Total, err = userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+
		" where unit_id = (select unit_id from users where id = ?)", user_id)
	if err != nil {
		log.Error("Error during getting meta unit user object from database %v", err)
		return nil, err
	}
	metaunituser.NumOfNotConfirmed, err = userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and confirmed = 0", user_id)
	if err != nil {
		log.Error("Error during getting meta unit user object from database %v", err)
		return nil, err
	}
	metaunituser.NumOfAdmins, err = userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+
		" where unit_id = (select unit_id from users where id = ?) and unitAdmin = 1", user_id)
	if err != nil {
		log.Error("Error during getting meta unit user object from database %v", err)
		return nil, err
	}

	return metaunituser, nil
}

func (userservice *UserService) GetMetaDashboard(user_id int64) (metadashboard *models.ApiMetaDashboard, err error) {
	metadashboard = new(models.ApiMetaDashboard)
	count, err := userservice.DbContext.SelectInt("select count(*) from "+userservice.Table+" where id = ? and newsBlocked = 1", user_id)
	if err != nil {
		log.Error("Error during getting meta unit user object from database %v", err)
		return nil, err
	}
	metadashboard.NewsBlocked = count != 0
	metadashboard.InvoiceTotal, err = userservice.DbContext.SelectInt("select count(*) from invoices "+
		"where company_id in (select id from companies where unit_id = (select unit_id from users where id = ?))", user_id)
	if err != nil {
		log.Error("Error during getting meta unit user object from database %v", err)
		return nil, err
	}

	return metadashboard, nil
}

func (userservice *UserService) InitUnit(trans *gorp.Transaction) (unitid int64, err error) {
	unit := new(models.DtoUnit)
	unit.Created = time.Now()
	unit.Active = true
	unit.Subscribed = false
	unit.Paid = false
	err = userservice.UnitRepository.Create(unit, trans)
	if err != nil {
		log.Error("Error during creating user object in database %v", err)
		return 0, err
	}

	return unit.ID, nil
}

func (userservice *UserService) Create(user *models.DtoUser, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = userservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during creating user object in database %v", err)
			return err
		}
	}

	if user.UnitID == 0 {
		user.UnitID, err = userservice.InitUnit(trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		err = trans.Insert(user)
	} else {
		err = userservice.DbContext.Insert(user)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating user object in database %v", err)
		return err
	}

	for _, email := range *user.Emails {
		email.UserID = user.ID
		if !email.Exists {
			err = userservice.EmailRepository.Create(&email, trans)
		} else {
			err = userservice.EmailRepository.Update(&email, trans)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during creating user object in database %v with value %v", err, email.Email)
			return err
		}
	}

	for _, phone := range *user.MobilePhones {
		phone.UserID = user.ID
		if !phone.Exists {
			err = userservice.MobilePhoneRepository.Create(&phone, trans)
		} else {
			err = userservice.MobilePhoneRepository.Update(&phone, trans)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during creating user object in database %v with value %v", err, phone.Phone)
			return err
		}
	}

	err = userservice.GroupRepository.SetByUser(user.ID, &user.Roles, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during creating user object in database %v with value %v", err, user.ID)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during creating user object in database %v", err)
			return err
		}
	}

	return nil
}

func (userservice *UserService) UpdateProfile(user *models.DtoUser) (err error) {
	_, err = userservice.DbContext.Exec("update "+userservice.Table+
		" set name = ?, language = ?, surname = ?, middleName = ?, workPhone = ?, jobTitle = ? where id = ?",
		user.Name, user.Language, user.Surname, user.MiddleName, user.WorkPhone, user.JobTitle, user.ID)
	if err != nil {
		log.Error("Error during updating user profile object in database %v with value %v", err, user.ID)
		return err
	}

	return nil
}

func (userservice *UserService) UpdatePassword(user *models.DtoUser) (err error) {
	_, err = userservice.DbContext.Exec("update "+userservice.Table+" set password = ? where id = ?", user.Password, user.ID)
	if err != nil {
		log.Error("Error during updating user password object in database %v with value %v", err, user.ID)
		return err
	}

	return nil
}

func (userservice *UserService) UpdateEmails(user *models.DtoUser, inTrans bool) (err error) {
	var trans *gorp.Transaction
	current := new(models.DtoUser)

	current.ID = user.ID
	current, err = userservice.GetUserArrays(current)
	if err != nil {
		log.Error("Error during updating user emails object in database %v with value %v", err, current.ID)
		return err
	}

	if inTrans {
		trans, err = userservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating user emails object in database %v", err)
			return err
		}
	}

	for _, updEmail := range *user.Emails {
		found := false
		for _, curEmail := range *current.Emails {
			if curEmail.Email == updEmail.Email {
				found = true
				break
			}
		}

		if !found {
			if !updEmail.Exists {
				err = userservice.EmailRepository.Create(&updEmail, trans)
			} else {
				err = userservice.EmailRepository.Update(&updEmail, trans)
			}
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user emails object in database %v with value %v", err, updEmail.Email)
				return err
			}
		} else {
			err = userservice.EmailRepository.Update(&updEmail, trans)
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user emails object in database %v with value %v", err, updEmail.Email)
				return err
			}
		}
	}

	for _, curEmail := range *current.Emails {
		found := false
		for _, updEmail := range *user.Emails {
			if curEmail.Email == updEmail.Email {
				found = true
				break
			}
		}
		if !found {
			err = userservice.EmailRepository.Delete(curEmail.Email, trans)
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user emails object in database %v with value %v", err, curEmail.Email)
				return err
			}
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating user emails object in database %v", err)
			return err
		}
	}

	return nil
}

func (userservice *UserService) UpdateMobilePhones(user *models.DtoUser, inTrans bool) (err error) {
	var trans *gorp.Transaction
	current := new(models.DtoUser)

	current.ID = user.ID
	current, err = userservice.GetUserArrays(current)
	if err != nil {
		log.Error("Error during updating user mobile phones object in database %v with value %v", err, current.ID)
		return err
	}

	if inTrans {
		trans, err = userservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating user mobile phones object in database %v", err)
			return err
		}
	}

	for _, updPhone := range *user.MobilePhones {
		found := false
		for _, curPhone := range *current.MobilePhones {
			if curPhone.Phone == updPhone.Phone {
				found = true
				break
			}
		}

		if !found {
			if !updPhone.Exists {
				err = userservice.MobilePhoneRepository.Create(&updPhone, trans)
			} else {
				err = userservice.MobilePhoneRepository.Update(&updPhone, trans)
			}
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user mobile phones object in database %v with value %v", err, updPhone.Phone)
				return err
			}
		} else {
			err = userservice.MobilePhoneRepository.Update(&updPhone, trans)
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user mobile phones object in database %v with value %v", err, updPhone.Phone)
				return err
			}
		}
	}

	for _, curPhone := range *current.MobilePhones {
		found := false
		for _, updPhone := range *user.MobilePhones {
			if curPhone.Phone == updPhone.Phone {
				found = true
				break
			}
		}
		if !found {
			err = userservice.MobilePhoneRepository.Delete(curPhone.Phone, trans)
			if err != nil {
				if inTrans {
					_ = trans.Rollback()
				}
				log.Error("Error during updating user mobile phones object in database %v with value %v", err, curPhone.Phone)
				return err
			}
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating user mobile phones object in database %v", err)
			return err
		}
	}

	return nil
}

func (userservice *UserService) Update(user *models.DtoUser, briefly bool, inTrans bool) (err error) {
	var trans *gorp.Transaction
	current := new(models.DtoUser)

	if !briefly {
		current.ID = user.ID
		current, err = userservice.GetUserArrays(current)
		if err != nil {
			log.Error("Error during updating user object in database %v with value %v", err, current.ID)
			return err
		}
	}

	if inTrans {
		trans, err = userservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during updating user object in database %v", err)
			return err
		}
	}

	if user.UnitID == 0 {
		user.UnitID, err = userservice.InitUnit(trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			return err
		}
	}

	if inTrans {
		_, err = trans.Update(user)
	} else {
		_, err = userservice.DbContext.Update(user)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during updating user object in database %v with value %v", err, user.ID)
		return err
	}

	if !briefly {
		err = userservice.GroupRepository.SetByUser(user.ID, &user.Roles, trans)
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during updating user object in database %v with value %v", err, user.ID)
			return err
		}

		for _, updEmail := range *user.Emails {
			found := false
			for _, curEmail := range *current.Emails {
				if curEmail.Email == updEmail.Email {
					found = true
					break
				}
			}

			if !found {
				if !updEmail.Exists {
					err = userservice.EmailRepository.Create(&updEmail, trans)
				} else {
					err = userservice.EmailRepository.Update(&updEmail, trans)
				}
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updEmail.Email)
					return err
				}
			} else {
				err = userservice.EmailRepository.Update(&updEmail, trans)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updEmail.Email)
					return err
				}
			}
		}

		for _, curEmail := range *current.Emails {
			found := false
			for _, updEmail := range *user.Emails {
				if curEmail.Email == updEmail.Email {
					found = true
					break
				}
			}
			if !found {
				err = userservice.EmailRepository.Delete(curEmail.Email, trans)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, curEmail.Email)
					return err
				}
			}
		}

		for _, updPhone := range *user.MobilePhones {
			found := false
			for _, curPhone := range *current.MobilePhones {
				if curPhone.Phone == updPhone.Phone {
					found = true
					break
				}
			}

			if !found {
				if !updPhone.Exists {
					err = userservice.MobilePhoneRepository.Create(&updPhone, trans)
				} else {
					err = userservice.MobilePhoneRepository.Update(&updPhone, trans)
				}
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updPhone.Phone)
					return err
				}
			} else {
				err = userservice.MobilePhoneRepository.Update(&updPhone, trans)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updPhone.Phone)
					return err
				}
			}
		}

		for _, curPhone := range *current.MobilePhones {
			found := false
			for _, updPhone := range *user.MobilePhones {
				if curPhone.Phone == updPhone.Phone {
					found = true
					break
				}
			}
			if !found {
				err = userservice.MobilePhoneRepository.Delete(curPhone.Phone, trans)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, curPhone.Phone)
					return err
				}
			}
		}
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during updating user object in database %v", err)
			return err
		}
	}

	return nil
}

func (userservice *UserService) Delete(userid int64, inTrans bool) (err error) {
	var trans *gorp.Transaction

	if inTrans {
		trans, err = userservice.DbContext.Begin()
		if err != nil {
			log.Error("Error during deleting user object in database %v", err)
			return err
		}
	}

	if inTrans {
		_, err = trans.Exec("update reports set user_id = 0 where user_id = ?", userid)
	} else {
		_, err = userservice.DbContext.Exec("update reports set user_id = 0 where user_id = ?", userid)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	if inTrans {
		_, err = trans.Exec("update feedback set user_id = 0 where user_id = ?", userid)
	} else {
		_, err = userservice.DbContext.Exec("update feedback set user_id = 0 where user_id = ?", userid)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	if inTrans {
		_, err = trans.Exec("update orders set user_id = 0 where user_id = ?", userid)
	} else {
		_, err = userservice.DbContext.Exec("update orders set user_id = 0 where user_id = ?", userid)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.DeviceRepository.DeleteByUser(userid, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	err = userservice.MobilePhoneRepository.DeleteByUser(userid, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	err = userservice.MessageRepository.DeleteByUser(userid, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.SessionRepository.DeleteByUser(userid, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.EmailRepository.DeleteByUser(userid, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	if inTrans {
		_, err = trans.Exec("update "+userservice.Table+" set user_id = 0 where user_id = ?", userid)
	} else {
		_, err = userservice.DbContext.Exec("update "+userservice.Table+" set user_id = 0 where user_id = ?", userid)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.GroupRepository.SetByUser(userid, &[]models.UserRole{}, trans)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	if inTrans {
		_, err = trans.Exec("delete from "+userservice.Table+" where id = ?", userid)
	} else {
		_, err = userservice.DbContext.Exec("delete from "+userservice.Table+" where id = ?", userid)
	}
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	if inTrans {
		err = trans.Commit()
		if err != nil {
			log.Error("Error during deleting user object in database %v", err)
			return err
		}
	}

	return nil
}
