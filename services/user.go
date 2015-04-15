package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
	"time"
)

type UserRepository interface {
	GetUserArrays(user *models.DtoUser) (*models.DtoUser, error)
	FindByLogin(login string) (user *models.DtoUser, err error)
	FindByCode(code string) (user *models.DtoUser, err error)
	Get(userid int64) (user *models.DtoUser, err error)
	GetAll(filter string) (users *[]models.ApiUserShort, err error)
	GetByUnit(unitid int64) (users *[]models.ApiUserTiny, err error)
	GetMeta() (usermeta *models.ApiUserMeta, err error)
	InitUnit(trans *gorp.Transaction, inTrans bool) (unitid int64, err error)
	Create(user *models.DtoUser, inTrans bool) (err error)
	Update(user *models.DtoUser, briefly bool, inTrans bool) (err error)
	Delete(userid int64, inTrans bool) (err error)
}

type UserService struct {
	SessionRepository     SessionRepository
	EmailRepository       EmailRepository
	UnitRepository        UnitRepository
	GroupRepository       GroupRepository
	MessageRepository     MessageRepository
	MobilePhoneRepository MobilePhoneRepository
	*Repository
}

func NewUserService(repository *Repository) *UserService {
	repository.DbContext.AddTableWithName(models.DtoUser{}, repository.Table).SetKeys(true, "id")
	return &UserService{Repository: repository}
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
		"select * from "+userservice.Table+" where id in (select user_id from emails where `primary` = 1 and email = ?) or"+
			" id in (select user_id from mobile_phones where `primary` = 1 and phone = ?)", login, login)
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
	_, err = userservice.DbContext.Select(users, "select u.id, "+
		"case when e.primary = 1 then e.email else case when m.primary = 1 then m.phone else '' end end as login,"+
		" not u.active as blocked, u.confirmed, u.lastLogin as lastLoginAt, u.surname, u.name, u.middleName from "+userservice.Table+
		" u left join emails e on u.id = e.user_id left join mobile_phones m on u.id = m.user_id"+
		" where (e.primary = 1 or e.primary is null) or (m.primary = 1 or m.primary is null)"+filter)
	if err != nil {
		log.Error("Error during getting user objects from database %v", err)
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

func (userservice *UserService) InitUnit(trans *gorp.Transaction, inTrans bool) (unitid int64, err error) {
	unit := new(models.DtoUnit)
	unit.Created = time.Now()
	unit.Active = true
	err = userservice.UnitRepository.Create(unit)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
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
		user.UnitID, err = userservice.InitUnit(trans, inTrans)
		if err != nil {
			return err
		}
	}

	err = userservice.DbContext.Insert(user)
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
			err = userservice.EmailRepository.Create(&email)
		} else {
			err = userservice.EmailRepository.Update(&email)
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
			err = userservice.MobilePhoneRepository.Create(&phone)
		} else {
			err = userservice.MobilePhoneRepository.Update(&phone)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during creating user object in database %v with value %v", err, phone.Phone)
			return err
		}
	}

	err = userservice.GroupRepository.SetByUser(user.ID, &user.Roles, false)
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
		user.UnitID, err = userservice.InitUnit(trans, inTrans)
		if err != nil {
			return err
		}
	}

	_, err = userservice.DbContext.Update(user)
	if err != nil {
		log.Error("Error during updating user object in database %v with value %v", err, user.ID)
		return err
	}

	if !briefly {
		err = userservice.GroupRepository.SetByUser(user.ID, &user.Roles, false)
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
					err = userservice.EmailRepository.Create(&updEmail)
				} else {
					err = userservice.EmailRepository.Update(&updEmail)
				}
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updEmail.Email)
					return err
				}
			} else {
				err = userservice.EmailRepository.Update(&updEmail)
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
				err = userservice.EmailRepository.Delete(curEmail.Email)
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
					err = userservice.MobilePhoneRepository.Create(&updPhone)
				} else {
					err = userservice.MobilePhoneRepository.Update(&updPhone)
				}
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, updPhone.Phone)
					return err
				}
			} else {
				err = userservice.MobilePhoneRepository.Update(&updPhone)
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
				err = userservice.MobilePhoneRepository.Delete(curPhone.Phone)
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

	_, err = userservice.DbContext.Exec("update orders set user_id = 0 where user_id = ?", userid)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.MobilePhoneRepository.DeleteByUser(userid)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	err = userservice.MessageRepository.DeleteByUser(userid, false)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.SessionRepository.DeleteByUser(userid, false)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.EmailRepository.DeleteByUser(userid)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	_, err = userservice.DbContext.Exec("update "+userservice.Table+" set user_id = 0 where user_id = ?", userid)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.GroupRepository.SetByUser(userid, &[]models.UserRole{}, false)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object in database %v with value %v", err, userid)
		return err
	}

	_, err = userservice.DbContext.Exec("delete from "+userservice.Table+" where id = ?", userid)
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
