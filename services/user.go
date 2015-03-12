package services

import (
	"application/models"
	"github.com/coopernurse/gorp"
	"time"
)

type UserService struct {
	SessionService *SessionService
	EmailService   *EmailService
	UnitService    *UnitService
	GroupService   *GroupService
	*Repository
}

func NewUserService(repository *Repository) *UserService {
	repository.DbContext.AddTableWithName(models.DtoUser{}, repository.Table).SetKeys(true, "id")
	return &UserService{Repository: repository}
}

func (userservice *UserService) GetUserArrays(user *models.DtoUser) (*models.DtoUser, error) {
	roles, err := userservice.GroupService.GetByUser(user.ID)
	if err != nil {
		log.Error("Error during getting user roles object from database %v with value %v", err, user.ID)
		return nil, err
	}
	user.Roles = *roles

	emails, err := userservice.EmailService.GetByUser(user.ID)
	if err != nil {
		log.Error("Error during getting user roles object from database %v with value %v", err, user.ID)
		return nil, err
	}
	user.Emails = emails

	return user, nil
}

func (userservice *UserService) FindByLogin(login string) (user *models.DtoUser, err error) {
	user = new(models.DtoUser)
	err = userservice.DbContext.SelectOne(user,
		"select u.* from "+userservice.Table+" u inner join emails e on u.id = e.user_id where e.primary = 1 and e.email = ?",
		login)
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
	_, err = userservice.DbContext.Select(users,
		"select u.id, coalesce(e.email, '') as login, not u.active as blocked, u.confirmed, u.lastLogin as lastLoginAt, u.name from "+
			userservice.Table+" u left join emails e on u.id = e.user_id where (e.primary = true or e.primary is null)"+filter)
	if err != nil {
		log.Error("Error during getting user objects from database %v", err)
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
	unit.Name = "Default name for unit"
	unit.Created = time.Now()
	err = userservice.UnitService.Create(unit)
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
			err = userservice.EmailService.Create(&email)
		} else {
			err = userservice.EmailService.Update(&email)
		}
		if err != nil {
			if inTrans {
				_ = trans.Rollback()
			}
			log.Error("Error during creating user object in database %v", err)
			return err
		}
	}

	err = userservice.GroupService.SetByUser(user.ID, &user.Roles, false)
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
		log.Error("Error during updating user object in database %v", err)
		return err
	}

	if !briefly {
		err = userservice.GroupService.SetByUser(user.ID, &user.Roles, false)
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
					err = userservice.EmailService.Create(&updEmail)
				} else {
					err = userservice.EmailService.Update(&updEmail)
				}
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v", err)
					return err
				}
			} else {
				err = userservice.EmailService.Update(&updEmail)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v", err)
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
				err = userservice.EmailService.Delete(curEmail.Email)
				if err != nil {
					if inTrans {
						_ = trans.Rollback()
					}
					log.Error("Error during updating user object in database %v with value %v", err, curEmail.Email)
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

	err = userservice.SessionService.DeleteByUser(userid, false)
	if err != nil {
		if inTrans {
			_ = trans.Rollback()
		}
		log.Error("Error during deleting user object from database %v with value %v", err, userid)
		return err
	}

	err = userservice.EmailService.DeleteByUser(userid)
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

	err = userservice.GroupService.SetByUser(userid, &[]models.UserRole{}, false)
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
